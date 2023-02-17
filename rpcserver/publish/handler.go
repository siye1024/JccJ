package publishsrv

import (
	"bytes"
	"context"
	"dousheng/controller/xhttp"
	"dousheng/pkg/minio"
	publish "dousheng/rpcserver/kitex_gen/publish"
	publishsrv "dousheng/rpcserver/kitex_gen/publish/publishsrv"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/gofrs/uuid"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

// PublishSrvImpl implements the last service interface defined in the IDL.
type PublishSrvImpl struct{}

// PublishAction implements the PublishSrvImpl interface.
func (s *PublishSrvImpl) PublishAction(ctx context.Context, req *publish.DouyinPublishActionRequest) (resp *publish.DouyinPublishActionResponse, err error) {
	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if len(req.Data) == 0 || len(req.Title) == 0 {
		return nil, kerrors.NewBizStatusError(20001, "Empty Video Data or Empty Title")
	}

	MinioVideoBucketName := xminio.bucketName
	videoData := []byte(req.Data)

	reader := bytes.NewReader(videoData)
	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	fileName := u2.String() + ".mp4" // assign video name across uuid
	err = xminio.UploadFile(MinioVideoBucketName, fileName, reader, int64(len(videoData)))
	if err != nil {
		return nil, err
	}

	return
}

// PublishList implements the PublishSrvImpl interface.
func (s *PublishSrvImpl) PublishList(ctx context.Context, req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	// TODO: Your code here...
	return
}

func (s *PublishSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := publishsrv.NewServer(new(PublishSrvImpl),
		server.WithServiceAddr(addr),                            // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler), //support kerrors
		//server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		//server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		//server.WithSuite(tracing.NewServerSuite()),                         // trace
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "videoPublish"}),
	)

	log.Println("Start Publish RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *PublishSrvImpl) Stop() {
	log.Println("Stop Publish RPC service...")
}
