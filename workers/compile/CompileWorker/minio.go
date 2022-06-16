package main

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func uploader(minioClient *minio.Client, BucketID int, FileName string) error {
	cntx := context.Background()
	filePath := fmt.Sprintf("%v/%v", fmt.Sprintf("%d", BucketID), FileName)

	_, err := minioClient.FPutObject(cntx, "gifs", filePath, "./"+filePath, minio.PutObjectOptions{ContentType: "image/gif"})
	return err
}

func downloader(minioClient *minio.Client, BucketID int, VideoName string, FileName string) error {
	cntx := context.Background()
	filePath := fmt.Sprintf("%v/%v/%v", fmt.Sprintf("%d", BucketID), VideoName, FileName)

	err := minioClient.FGetObject(cntx, "resize", filePath, "./"+filePath, minio.GetObjectOptions{})
	return err
}
