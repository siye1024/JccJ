package minio

//func init() {
//	ctx := context.Background()
//	endpoint := "play.min.io"
//	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
//	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
//	useSSL := true
//
//	minioClient, err := minio.New(endpoint, &minio.Options{
//		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
//		Secure: useSSL,
//	})
//	if err != nil {
//		klog.Errorf("minio client init failed: %v", err)
//	}
//
//	klog.Debug("minio client init successfully")
//
//	bucketName := "dousheng_JJCC"
//	location := "us-east-1"
//
//	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
//	if err != nil {
//		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
//		if errBucketExists == nil && exists {
//			klog.Debugf("bucket %s already exists", bucketName)
//			return nil
//		} else {
//			return err
//		}
//	} else {
//		klog.Infof("bucket %s create successfully", bucketName)
//	}
//	return nil
//	if err := CreateBucket(bucketName); err != nil {
//		klog.Errorf("minio client init failed: %v", err)
//	}
//}
