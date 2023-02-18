package xrpc

import (
	"context"
	"dousheng/rpcserver/kitex_gen/favorite"
	"dousheng/rpcserver/kitex_gen/favorite/favoritesrv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"time"
)

var favoriteClient favoritesrv.Client

// Favorite RPC 客户端初始化
func initFavoriteRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	c, err := favoritesrv.NewClient(
		"favoriteService",
		client.WithTransportProtocol(transport.TTHeader),        //to support mux read check protocol
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler), //to support kerrors
		client.WithMuxConnection(1),                             // mux
		client.WithRPCTimeout(30*time.Second),                   // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond),       // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),       // retry
		client.WithResolver(r),                                  // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "favoriteService"}),
	)
	if err != nil {
		panic(err)
	}
	favoriteClient = c
}

// 传递 点赞操作 的上下文, 并获取 RPC Server 端的响应.
func FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	resp, err = favoriteClient.FavoriteAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 传递 获取点赞列表操作 的上下文, 并获取 RPC Server 端的响应.
func FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (resp *favorite.DouyinFavoriteListResponse, err error) {
	resp, err = favoriteClient.FavoriteList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
