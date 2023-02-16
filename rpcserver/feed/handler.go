package feedsrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/db"
	"dousheng/pkg/pack"
	feed "dousheng/rpcserver/feed/kitex_gen/feed"
	feedsrv "dousheng/rpcserver/feed/kitex_gen/feed/feedsrv"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

const (
	LIMIT = 30 // 单次返回最大视频数
)

// FeedSrvImpl implements the last service interface defined in the IDL.
type FeedSrvImpl struct{}

// GetUserFeed implements the FeedSrvImpl interface.
func (s *FeedSrvImpl) GetUserFeed(ctx context.Context, req *feed.DouyinFeedRequest) (resp *feed.DouyinFeedResponse, err error) {
	var (
		uid           int64 = 0
		nextTime      int64
		respStatusMsg = "User Register Success"
	)
	//do not need to check latest time again
	//check Token
	if len(*req.Token) != 0 {
		claim, err := xhttp.Jwt.ParseToken(*req.Token)
		if err != nil {
			return nil, err
		} else {
			uid = claim.Id
		}
	}

	videos, err := db.MGetVideos(ctx, LIMIT, req.LatestTime)
	if err != nil {
		return nil, err
	}

	if len(videos) < LIMIT {
		nextTime = time.Now().UnixMilli() //reset time if rest videos are less than LIMIT
	} else {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}

	vis, err := pack.PackVideos(ctx, videos, &uid)
	if err != nil {
		return nil, kerrors.NewBizStatusError(20002, "Videos Pack Error")
	}

	resp = &feed.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  &respStatusMsg,
		VideoList:  vis,
		NextTime:   &nextTime,
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
