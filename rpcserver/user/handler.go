package usersrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/db"
	"dousheng/pkg/jwt"
	user "dousheng/rpcserver/kitex_gen/user"
	svr "dousheng/rpcserver/kitex_gen/user/usersrv"
	api2 "dousheng/rpcserver/user/api"
	"errors"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/gorm"
	"log"
	"net"
	"time"
)

// Arg2Config Directly set the parameters
var (
	Arg2Config = &api2.Argon2Params{
		Memory:      65536,
		Iterations:  3,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Register Success"
	)

	// empty username or password has been processed by user client
	if len(req.Username) == 0 || len(req.Password) == 0 {
		err := kerrors.NewBizStatusError(10001, "Empty Username or Password")
		return nil, err
	}
	// len of username or password exceed the 32 bit
	if len(req.Username) > 32 || len(req.Password) > 32 {
		err := kerrors.NewBizStatusError(10014, "Too Long Username or Password")
		return nil, err
	}

	err = api2.NewCreateUserOp(ctx).CreateUser(req, Arg2Config)
	if err != nil {
		return nil, err
	}

	//Auto Login
	uid, err := api2.NewCheckUserOp(ctx).CheckUser(req)
	if err != nil {
		return nil, err
	}

	token, err := xhttp.Jwt.CreateToken(jwt.CustomClaims{ //Claim is payload
		Id:   int64(uid),
		Time: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	resp = &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserId:     uid,
		Token:      token, // successful resp must have token
	}

	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Login Success"
		respErrorMag  = "Signature is Invalid"
	)
	// empty username or password has been processed by dousheng client
	if len(req.Username) == 0 || len(req.Password) == 0 {
		err := kerrors.NewBizStatusError(10001, "Empty Username or Password")
		return nil, err
	}
	// len of username or password exceed 32 bit
	if len(req.Username) > 32 || len(req.Password) > 32 {
		err := kerrors.NewBizStatusError(10014, "Too Long Username or Password")
		return nil, err
	}

	// check the user's information
	uid, err := api2.NewCheckUserOp(ctx).CheckUser(req)
	if err != nil {
		resp = &user.DouyinUserRegisterResponse{
			StatusCode: 10012,
			StatusMsg:  &respErrorMag,
		}
		return resp, nil
	}

	token, err := xhttp.Jwt.CreateToken(jwt.CustomClaims{ //Claim is payload
		Id:   int64(uid),
		Time: time.Now().Unix(),
	})
	if err != nil {
		resp = &user.DouyinUserRegisterResponse{
			StatusCode: 10012,
			StatusMsg:  &respErrorMag,
		}
		return resp, nil
	}

	resp = &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserId:     uid,
		Token:      token, // successful resp must have token
	}

	return resp, nil
}

// GetUserById implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserById(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	var (
		respStatusMsg = "Get User's Info By ID Successfully"
	)

	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if req.UserId < 0 {
		return nil, kerrors.NewBizStatusError(10007, "Invalid Username")
	}

	u := new(user.User)
	if err := db.DB.WithContext(ctx).First(&u, req.UserId).Error; err != nil {
		return nil, err
	}

	if u == nil || u.Id == 0 {
		err := kerrors.NewBizStatusError(10009, "User Already Withdraw")
		return nil, err
	}

	// true means the claim.id has follow the modelUser.id, false means not follow
	isFollow := false
	if claim.Id == u.Id {
		isFollow = true
	} else {
		relation := new(db.Relation)
		err := db.DB.WithContext(ctx).First(&relation, "user_id = ? and to_user_id = ?", claim.Id, u.Id).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if relation != nil && relation.UserID != 0 { //double check is necessary
			isFollow = true
		}

	}
	works, err := db.PublishList(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	workCount := int64(len(works))

	favs, err := db.FavoriteList(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	favCount := int64(len(favs))

	userInfo := &user.User{
		Id:             int64(u.Id),
		Name:           u.Name,
		FollowCount:    u.FollowCount,
		FollowerCount:  u.FollowerCount,
		IsFollow:       isFollow,
		WorkCount:      &workCount,
		FavoriteCount:  &favCount,
		TotalFavorited: u.TotalFavorited,
	}

	resp = &user.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		User:       userInfo,
	}
	return resp, nil
}

func (s *UserSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := svr.NewServer(new(UserSrvImpl),
		server.WithServiceAddr(addr),                            // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler), //support kerrors
		//server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "userRegisterLoginGetInfo"}),
	)

	log.Println("Start User RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *UserSrvImpl) Stop() {
	log.Println("Stop User RPC service...")
}
