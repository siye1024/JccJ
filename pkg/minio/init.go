package xminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var (
	minioClient     *minio.Client
	Endpoint        = "play.min.io"
	AccessKeyID     = "Q3AM3UQ867SPQQA43P2F"
	SecretAccessKey = "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	UseSSL          = true
	BucketName      = "dousheng_JccJ"
	Location        = "chengdu"
)

func init() {

	ctx := context.Background()

	client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: UseSSL,
	})
	if err != nil {
		log.Fatalln("minio client init failed: %v", err)
	}

	log.Println("minio client init successfully")
	minioClient = client

	err = minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: Location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, BucketName)
		if errBucketExists == nil && exists {
			log.Println("bucket %s already exists", BucketName)
		} else {
			log.Fatalln("minio client init failed: %v", err)
		}
	} else {
		log.Println("bucket %s create successfully", BucketName)
	}
}
