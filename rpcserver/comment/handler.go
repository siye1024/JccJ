package commentsrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/comment/api"
	comment "dousheng/rpcserver/kitex_gen/comment"
	commentsrv "dousheng/rpcserver/kitex_gen/comment/commentsrv"
	"dousheng/rpcserver/kitex_gen/user"
	"fmt"
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

// CommentSrvImpl implements the last service interface defined in the IDL.
type CommentSrvImpl struct{}

// CommentAction implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentAction(ctx context.Context, req *comment.DouyinCommentActionRequest) (resp *comment.DouyinCommentActionResponse, err error) {
	var (
		respStatusMsg = "User Comment Success"
	)
	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		err := kerrors.NewBizStatusError(30007, "Get Token Error")
		return nil, err
	}

	if req.UserId == 0 || claim.Id != 0 {
		req.UserId = claim.Id
	}

	if req.ActionType != 1 && req.ActionType != 2 || req.UserId == 0 || req.VideoId == 0 {
		err := kerrors.NewBizStatusError(30010, "Error occurred while binding the request body to the struct")
		return nil, err
	}

	err = api.NewCommentActionService(ctx).CommentAction(req)
	if err != nil {
		return nil, err
	}

	userinfo, err := xrpc.GetUserById(ctx, &user.DouyinUserRequest{
		UserId: req.UserId,
		Token:  req.Token,
	})
	t := time.Now()
	tFormat := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	commentInfo := comment.Comment{
		Id:         int64(req.VideoId),
		User:       userinfo.User,
		Content:    *req.CommentText,
		CreateDate: tFormat,
	}

	resp = &comment.DouyinCommentActionResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		Comment:    &commentInfo,
	}
	return resp, nil
}

// CommentList implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentList(ctx context.Context, req *comment.DouyinCommentListRequest) (resp *comment.DouyinCommentListResponse, err error) {
	var (
		respStatusMsg = "Get User's Comment List Successfully"
	)

	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if req.VideoId == 0 || claim.Id == 0 {
		err := kerrors.NewBizStatusError(30010, "Error occurred while binding the request body to the struct")
		return nil, err
	}

	comments, err := api.NewCommentListService(ctx).CommentList(req, claim.Id)
	if err != nil {
		return nil, err
	}

	resp = &comment.DouyinCommentListResponse{
		StatusCode:  0,
		StatusMsg:   &respStatusMsg,
		CommentList: comments,
	}
	return resp, nil
}

func (s *CommentSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := commentsrv.NewServer(new(CommentSrvImpl),
		server.WithServiceAddr(addr),                            // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler), //support kerrors
		//server.WithMiddleware(middleware.CommonMiddleware),               // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),               // middleware
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "commentService"}),
	)

	log.Println("Start Comment RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *CommentSrvImpl) Stop() {
	log.Println("Stop Comment RPC service...")
}
