package feedsrv

import (
	"context"
	feed "dousheng/rpcserver/feed/kitex_gen/feed"
	feedsrv "dousheng/rpcserver/feed/kitex_gen/feed/feedsrv"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

// FeedSrvImpl implements the last service interface defined in the IDL.
type FeedSrvImpl struct{}

// GetUserFeed implements the FeedSrvImpl interface.
func (s *FeedSrvImpl) GetUserFeed(ctx context.Context, req *feed.DouyinFeedRequest) (resp *feed.DouyinFeedResponse, err error) {
	// TODO: Your code here...
	log.Println("correct")
	var (
		respStatusMsg = "User Register Success"
	)
	resp = &feed.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
	}
	return resp, nil
}

// GetVideoById implements the FeedSrvImpl interface.
func (s *FeedSrvImpl) GetVideoById(ctx context.Context, req *feed.VideoIdRequest) (resp *feed.Video, err error) {
	// TODO: Your code here...
	return
}
func (s *FeedSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port

	svr := feedsrv.NewServer(new(FeedSrvImpl),
		server.WithServiceAddr(addr), // address
		//server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "videoFeed"}))

	log.Println("Start Feed RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *FeedSrvImpl) Stop() {
	log.Println("Stop Feed RPC service...")
}
