package xminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var (
	minioClient     *minio.Client
	endpoint        = "play.min.io"
	accessKeyID     = "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey = "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL          = true
	bucketName      = "dousheng_JccJ"
	location        = "chengdu"
)

func init() {

	ctx := context.Background()

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("minio client init failed: %v", err)
	}

	log.Println("minio client init successfully")
	minioClient = client

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Println("bucket %s already exists", bucketName)
		} else {
			log.Fatalln("minio client init failed: %v", err)
		}
	} else {
		log.Println("bucket %s create successfully", bucketName)
	}
}
