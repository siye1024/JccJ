package xrpc

import (
	"context"
	"dousheng/rpcserver/feed/kitex_gen/feed"
	"dousheng/rpcserver/feed/kitex_gen/feed/feedsrv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"time"
)

var feedClient feedsrv.Client

func initFeedRpc() {

	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	client, err := feedsrv.NewClient(
		"videoFeed",
		client.WithTransportProtocol(transport.TTHeader),        //to support mux read check protocol
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler), //to support kerrors
		//client.WithMiddleware(middleware.CommonMiddleware),
		//client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		//client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r), // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "videoFeed"}),
	)
	if err != nil {
		panic(err)
	}
	feedClient = client
}

// 传递 获取视频流操作 的上下文, 并获取 RPC Server 端的响应.
func Feed(ctx context.Context, req *feed.DouyinFeedRequest) (resp *feed.DouyinFeedResponse, err error) {
	resp, err = feedClient.GetUserFeed(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, err
	}
	return resp, nil
}
