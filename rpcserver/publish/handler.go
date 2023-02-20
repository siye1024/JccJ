package publishsrv

import (
	"bytes"
	"context"
	"dousheng/controller/xhttp"
	"dousheng/db"
	"dousheng/pkg/minio"
	"dousheng/pkg/pack"
	publish "dousheng/rpcserver/kitex_gen/publish"
	publishsrv "dousheng/rpcserver/kitex_gen/publish/publishsrv"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/gofrs/uuid"
	etcd "github.com/kitex-contrib/registry-etcd"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/jpeg"
	"log"
	"net"
	"os"
	"strings"
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
	//process video
	MinioVideoBucketName := xminio.BucketName
	videoData := []byte(req.Data)

	videoReader := bytes.NewReader(videoData)
	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	fileName := u2.String() + ".mp4" // assign video name across uuid
	// upload video
	err = xminio.UploadFile(MinioVideoBucketName, fileName, videoReader, int64(len(videoData)), "video/mp4")
	if err != nil {
		return nil, err
	}
	// get video url
	url, err := xminio.GetFileUrl(MinioVideoBucketName, fileName, 0) // 0 is the default expiry time (1 day)
	if err != nil {
		return nil, err
	}

	playUrl := strings.Split(url.String(), "?")[0]
	if err != nil {
		return nil, err
	}

	//process video cover
	u3, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	coverName := u3.String() + ".jpg"
	coverData, err := ReadFrameAsJpeg(playUrl)
	if err != nil {
		return nil, err
	}
	//upload video cover
	coverReader := bytes.NewReader(coverData)
	err = xminio.UploadFile(MinioVideoBucketName, coverName, coverReader, int64(len(coverData)), "cover/jpg")
	if err != nil {
		return nil, err
	}
	//get video cover url
	coverUrl, err := xminio.GetFileUrl(MinioVideoBucketName, coverName, 0)
	if err != nil {
		return nil, err
	}

	CoverUrl := strings.Split(coverUrl.String(), "?")[0]
	// 封装video
	videoModel := &db.Video{
		AuthorID:      claim.Id,
		PlayUrl:       playUrl,
		CoverUrl:      CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         req.Title,
	}
	err = db.CreateVideo(ctx, videoModel)
	if err != nil {
		return nil, err
	}

	respMsg := "Publish Video Success"
	resp = &publish.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  &respMsg,
	}
	return resp, nil
}

// PublishList implements the PublishSrvImpl interface.
func (s *PublishSrvImpl) PublishList(ctx context.Context, req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	var (
		respMsg = "Get Pulish List Successfully"
	)

	if req.UserId <= 0 {
		return nil, kerrors.NewBizStatusError(21001, "Invalid User or User Token")
	}

	videos, err := db.PublishList(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	video_list, err := pack.Videos(ctx, videos, &req.UserId)
	if err != nil {
		return nil, err
	}

	resp = &publish.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  &respMsg,
		VideoList:  video_list,
	}

	return resp, nil

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
		server.WithRegistry(r), // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(), // Multiplex
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

// ReadFrameAsJpeg
// 从视频流中截取一帧并返回 需要在本地环境中安装ffmpeg并将bin添加到环境变量
func ReadFrameAsJpeg(filePath string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)

	err := ffmpeg.Input(filePath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()

	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil) // jpeg存储空间比png小，更适合做缩略图

	return buf.Bytes(), err
}
