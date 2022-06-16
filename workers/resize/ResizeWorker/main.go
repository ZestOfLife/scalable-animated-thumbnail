package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func doMinioStartup() *minio.Client {
	cntx := context.Background()
	endpoint := "minio-svc"
	accessKeyID := "minio"
	secretAccessKey := "minio_pass"
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = minioClient.MakeBucket(cntx, "resized", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(cntx, "resized")
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", "resized")
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", "resized")
	}
	return minioClient
}

func main() {
	cntx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "resize-projector:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(cntx).Result()
	if err != nil {
		log.Fatal("Unable to connect to queue: ", err)
	}
	log.Println("Connected to queue")

	minioClient := doMinioStartup()

	var wg sync.WaitGroup
	wg.Add(1)
	work(client, minioClient)
	wg.Wait()
}
