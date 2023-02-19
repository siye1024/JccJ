package relationsrv

import (
	"context"
	"dousheng/controller/xhttp"
	relation "dousheng/rpcserver/kitex_gen/relation"
	relationsrv "dousheng/rpcserver/kitex_gen/relation/relationsrv"
	"dousheng/rpcserver/relation/api"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

// RelationSrvImpl implements the last service interface defined in the IDL.
type RelationSrvImpl struct{}

// RelationAction implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationAction(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	var (
		respStatusMsg = "User Relation Operate Successfully"
	)
	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if claim.Id > 0 {
		req.UserId = claim.Id
	} else {
		return nil, kerrors.NewBizStatusError(40003, "Invalid Token or User ID")
	}

	if req.ActionType < 1 || req.ActionType > 2 {
		err := kerrors.NewBizStatusError(40004, "Invalid Action Type")
		return nil, err
	}
	// TODO Here
	err = api.NewRelationActionOp(ctx).RelationAction(req)
	if err != nil {
		return nil, err
	}
	resp = &relation.DouyinRelationActionResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
	}

	return resp, nil
}

// RelationFollowList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	var (
		respStatusMsg = "Get Follow List Successfully"
		user_id       int64
	)
	if len(req.Token) == 0 { // tourist can read follow list
		user_id = 0
	} else {
		claim, err := xhttp.Jwt.ParseToken(req.Token)
		if err != nil {
			return nil, err
		}
		user_id = claim.Id
	}

	if req.UserId == 0 {
		req.UserId = user_id
	}
	if req.UserId == 0 && user_id == 0 {
		return nil, kerrors.NewBizStatusError(40004, "Invalid Action Type")
	}

	users, err := api.NewFollowingListOp(ctx).FollowingList(req, user_id)
	if err != nil {
		return nil, err
	}

	resp = &relation.DouyinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserList:   users,
	}

	return resp, nil
}

// RelationFollowerList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	var (
		respStatusMsg = "Get Follower List Successfully"
		user_id       int64
	)
	if len(req.Token) == 0 { // tourist can read follow list
		user_id = 0
	} else {
		claim, err := xhttp.Jwt.ParseToken(req.Token)
		if err != nil {
			return nil, err
		}
		user_id = claim.Id
	}

	if req.UserId == 0 {
		req.UserId = user_id
	}
	if req.UserId == 0 && user_id == 0 {
		return nil, kerrors.NewBizStatusError(40004, "Invalid Action Type")
	}

	users, err := api.NewFollowerListOp(ctx).FollowerList(req, user_id)
	if err != nil {
		return nil, err
	}

	resp = &relation.DouyinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		UserList:   users,
	}

	return resp, nil
}
func (s *RelationSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := relationsrv.NewServer(new(RelationSrvImpl),
		server.WithServiceAddr(addr),                            // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler), //support kerrors
		//server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "relationService"}),
	)

	log.Println("Start Relation RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *RelationSrvImpl) Stop() {
	log.Println("Stop Relation RPC service...")
}
