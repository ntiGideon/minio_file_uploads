package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"os"
)

func CreateBucket(minioClient *minio.Client) error {
	// creating a bucket at region 'us-east-1' with object locking enabled
	err := minioClient.MakeBucket(context.Background(), "mybucket", minio.MakeBucketOptions{
		Region:        "us-east-1",
		ObjectLocking: true,
	})
	if err != nil {
		fmt.Println("error creating bucket:", err)
		return err
	}
	fmt.Println("Successfully created bucket")
	return nil
}

func CreateBucketWithChecks(minioClient *minio.Client, bucketName string) {
	location := "us-east-1a"
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: location,
	})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			fmt.Println("Bucket already exists")
		} else {
			fmt.Println("Error creating bucket:", err)
		}
	}
	fmt.Println("Successfully created bucket", bucketName)
}

func PutObject(minionClient *minio.Client, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:", err)
		return err
	}
	uploadInfo, err := minionClient.PutObject(context.Background(), "mybucket-1", "myobject", file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return err
	}
	fmt.Println("Successfully uploaded file", uploadInfo)
	return nil
}

func GetObject(minioClient *minio.Client, bucketName string, objectName string) error {
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object:", err)
		return err
	}
	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			return
		}
	}(object)
	localFile, err := os.Create("/tmp/local-file.text")
	if err != nil {
		fmt.Println("Error creating local file:", err)
		return err
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {
			return
		}
	}(localFile)
	_, err = io.Copy(localFile, object)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return err
	}
	fmt.Println("Successfully downloaded file", objectName)
	return nil
}

func ListObjects(minioClient *minio.Client, bucketName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "myprefix",
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println("error while streaming the response from the object: ", object.Err)
			return object.Err
		}
		fmt.Println(object)
	}
	return nil
}
