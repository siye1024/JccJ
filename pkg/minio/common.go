package xminio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/url"
	"os"
	"time"
)

// UploadFile 上传文件（提供reader）至 minio
func UploadFile(bucketName string, objectName string, reader io.Reader, objectsize int64) error {
	ctx := context.Background()
	n, err := minioClient.PutObject(ctx, bucketName, objectName, reader, objectsize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		log.Println("upload %s of size %d failed, %s", bucketName, objectsize, err)
		return kerrors.NewBizStatusError(20002, "upload video failed")
	}
	log.Println("upload %s of bytes %d successfully", objectName, n.Size)
	return nil
}

// GetFileUrl 从 minio 获取文件Url
func GetFileUrl(bucketName string, fileName string, expires time.Duration) (*url.URL, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	if expires <= 0 { // set url expiry time , 1sec <= expires <= 7days
		expires = time.Second * 60 * 60 * 24
	}
	presignedUrl, err := minioClient.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		log.Println("get url of file %s from bucket %s failed, %s", fileName, bucketName, err)
		return nil, kerrors.NewBizStatusError(20002, "Fail to Get url of the Video")
	}
	// TODO: url可能要做截取
	return presignedUrl, nil
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
