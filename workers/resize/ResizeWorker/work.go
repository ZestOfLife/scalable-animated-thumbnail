package main

import (
	"ResizeWorker/queue"
	"ResizeWorker/senders"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
)

func work(client *redis.Client, minioClient *minio.Client) {
	for {
		if queue.GetDeadLen() > 0 {
			if queue.GetDeadLen() > 10 {
				log.Printf("More than 10 messages have failed")
			} else if queue.GetDeadLen() > 100 {
				log.Fatal("More than 100 messages have failed")
			}
			go senders.DoDeadMessages()
		}
		cntx := context.Background()
		res, err := client.BRPop(cntx, 0, "queue").Result()
		if err != nil {
			log.Fatal(err)
		}

		var job JobCertificate
		err2 := json.Unmarshal([]byte(res[1]), &job)
		if err2 != nil {
			log.Fatal(err2)
		}

		path := fmt.Sprintf("./%v/%v", fmt.Sprintf("%d", job.BucketID), job.VideoName)
		err3 := os.MkdirAll(path, os.ModePerm)
		if err3 != nil {
			log.Println("Error 3:")
			log.Println(err3)
			go senders.ResizeFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
			continue
		}

		err4 := downloader(minioClient, job.BucketID, job.VideoName, job.FileName)
		if err4 != nil {
			log.Println("Error 4:")
			log.Println(err4)
			go senders.ResizeFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
			continue
		}

		cmd := exec.Command("mogrify", path+"/"+job.FileName, "-resize", "720x540", path+"/"+job.FileName)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err5 := cmd.Run()
		if err5 != nil {
			log.Println("Error 5:")
			log.Println(err5)
			log.Println(out.String())
			log.Println(stderr.String())
			go senders.ResizeFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
			continue
		}

		err6 := uploader(minioClient, job.BucketID, job.VideoName, job.FileName)
		if err6 != nil {
			log.Println("Error 6:")
			log.Println(err6)
			go senders.ResizeFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
			continue
		}
		go senders.ResizeSuccess(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
	}
}
