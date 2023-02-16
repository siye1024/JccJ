package usersrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/db"
	"dousheng/pkg/jwt"
	user "dousheng/rpcserver/user/kitex_gen/user"
	svr "dousheng/rpcserver/user/kitex_gen/user/usersrv"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	"time"
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Register Success"
	)
	// empty username or password has been processed by user client
	users, err := db.QueryUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		err := kerrors.NewBizStatusError(10002, "User Already Exist")
		return nil, err
	}

	err = db.CreateUser(ctx, []*db.User{{
		UserName: req.Username,
		Password: req.Password,
	}})
	if err != nil {
		return nil, err
	}

	//TODO : AUOTO LOGIN
	//TODO please complete login func and replace code here
	users, err = db.QueryUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	loginUser := users[0]
	uid := int64(loginUser.ID)

	// Sign Key refers to xttp.common, is SoundDance here
	token, err := xhttp.Jwt.CreateToken(jwt.CustomClaims{ //Claim is payload
		Id:   int64(uid),
		Time: time.Now().Unix(),
	})

	//Register Success
	resp = &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserId:     uid,
		Token:      token,
	}

	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	var (
		respStatusMsg = "User Login Success"
	)
	// empty username or password has been processed by dousheng client
	users, err := db.QueryUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		err := kerrors.NewBizStatusError(10007, "Invalid Username")
		return nil, err
	}
	userLogin := users[0]
	if req.Password != userLogin.Password {
		err := kerrors.NewBizStatusError(100008, "Invalid Password")
		return nil, err
	}
	uid := int64(userLogin.ID)
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
		return nil, err
	}

	u := new(user.User)
	if err := db.DB.WithContext(ctx).First(&u, req.UserId).Error; err != nil {
		return nil, err
	}

	if u == nil {
		err := kerrors.NewBizStatusError(10009, "User Already Withdraw")
		return nil, err
	}

	follow_count := int64(*u.FollowCount)
	follower_count := int64(*u.FollowerCount)

	// true means the claim.id has follow the modelUser.id, false means not follow

	isFollow := false
	/*
		relation := new(db.Relation)
		if err := db.DB.WithContext(ctx).First(&relation, "user_id = ? and to_user_id = ?", claim.Id, int64(u.Id)).Error; err != nil {
			return nil, err
		}

		if relation != nil {
			isFollow = true
		}
	*/
	userInfo := &user.User{
		Id:            int64(u.Id),
		Name:          u.Name,
		FollowCount:   &follow_count,
		FollowerCount: &follower_count,
		IsFollow:      isFollow,
	}

	if claim.Id == req.UserId {
		userInfo.IsFollow = true
	} else {
		userInfo.IsFollow = false
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
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
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
