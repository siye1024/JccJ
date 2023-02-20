package xminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net"
	"strings"
)

var (
	minioClient     *minio.Client
	Endpoint        = ":9000"
	AccessKeyID     = "doushengMinio"
	SecretAccessKey = "doushengMinio"
	UseSSL          = false
	BucketName      = "doushengjccj"
	//Location        = "chengdu"
)

func InitMInio() {

	ctx := context.Background()
	//only for test
	ip, err := GetIPv4()
	if err != nil {
		log.Fatal("Minio Get IP Failed")
		return
	}

	//ip := "47.93.27.219"
	log.Println(ip)
	Endpoint = ip + Endpoint

	client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: UseSSL,
	})
	if err != nil {
		log.Fatalln("minio client init failed: %v", err)
	}

	log.Println("minio client init successfully")
	minioClient = client

	err = minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
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

func GetIPv4() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
