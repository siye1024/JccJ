package xminio

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
)

// UploadFile 上传文件（提供reader）至 minio
// put a video into the bucketName
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
