package xminio

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/url"
	"time"
)

// UploadFile 上传文件（提供reader）至 minio
func UploadFile(bucketName string, objectName string, reader io.Reader, objectsize int64, contentType string) error {
	ctx := context.Background()
	n, err := minioClient.PutObject(ctx, bucketName, objectName, reader, objectsize, minio.PutObjectOptions{
		ContentType: contentType, //"application/octet-stream"
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
