package xrpc

import (
	"context"
	"dousheng/rpcserver/publish/kitex_gen/publish"
	"dousheng/rpcserver/publish/kitex_gen/publish/publishsrv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"time"
)

var publishClient publishsrv.Client

// Publish RPC 客户端初始化
func initPublishRpc() {

	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	c, err := publishsrv.NewClient(
		"videoPublish",
		//client.WithMiddleware(middleware.CommonMiddleware),
		//client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		//client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r), // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "videoPublish"}),
	)
	if err != nil {
		panic(err)
	}
	publishClient = c
}

// 传递 发布视频操作 的上下文, 并获取 RPC Server 端的响应.
func PublishAction(ctx context.Context, req *publish.DouyinPublishActionRequest) (resp *publish.DouyinPublishActionResponse, err error) {
	resp, err = publishClient.PublishAction(ctx, req)
	if err != nil {
		return nil, err
	}
	//if resp.StatusCode != 0 {
	//	return nil, errno.NewErrNo(int(resp.StatusCode), *resp.StatusMsg)
	//}
	return resp, nil
}

// 传递 获取用户发布视频列表操作 的上下文, 并获取 RPC Server 端的响应.
func PublishList(ctx context.Context, req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	resp, err = publishClient.PublishList(ctx, req)
	if err != nil {
		return nil, err
	}
	//if resp.StatusCode != 0 {
	//	return nil, errno.NewErrNo(int(resp.StatusCode), *resp.StatusMsg)
	//}
	return resp, nil
}
