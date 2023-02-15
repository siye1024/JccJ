package xrpc

import (
	"context"
	"dousheng/rpcserver/user/kitex_gen/user"
	"dousheng/rpcserver/user/kitex_gen/user/usersrv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"time"
)

var userClient usersrv.Client

func initUserRpc() {

	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	client, err := usersrv.NewClient(
		"userRegisterLoginGetInfo",
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
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "userRegisterLoginGetInfo"}),
	)
	if err != nil {
		panic(err)
	}
	userClient = client
}

func Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp, err = userClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
