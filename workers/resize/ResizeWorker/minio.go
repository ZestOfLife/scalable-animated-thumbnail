package main

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func uploader(minioClient *minio.Client, BucketID int, VideoName string, FileName string) error {
	cntx := context.Background()
	filePath := fmt.Sprintf("%v/%v/%v", fmt.Sprintf("%d", BucketID), VideoName, FileName)

	_, err := minioClient.FPutObject(cntx, "resized", filePath, "./"+filePath, minio.PutObjectOptions{ContentType: "image/jpeg"})
	return err
}

func downloader(minioClient *minio.Client, BucketID int, VideoName string, FileName string) error {
	cntx := context.Background()
	filePath := fmt.Sprintf("%v/%v/%v", fmt.Sprintf("%d", BucketID), VideoName, FileName)

	err := minioClient.FGetObject(cntx, "extracted", filePath, "./"+filePath, minio.GetObjectOptions{})
	return err
}
