package main

import (
	"ExtractWorker/queue"
	"ExtractWorker/senders"
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
			go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
			continue
		}

		err4 := downloader(minioClient, job.BucketID, job.VideoName)
		if err4 != nil {
			log.Println("Error 4:")
			log.Println(err4)
			go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
			continue
		}

		timeAt := fmt.Sprintf("'%v", job.Timestamp) + "ms'"

		cmd := exec.Command("ffmpeg", "-ss", timeAt, "-i", path+"/"+job.VideoName, "-frames:v", "1", "-q:v", "2", path+"/"+job.FileName)
		err5 := cmd.Run()
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err5 != nil {
			log.Println("Error 5:")
			log.Println(err5)
			log.Println(out)
			log.Println(stderr)
			go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
			continue
		}

		err6 := uploader(minioClient, job.BucketID, job.VideoName, job.FileName)
		if err6 != nil {
			log.Println("Error 6:")
			log.Println(err6)
			go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
			continue
		}
		go senders.ExtractSuccess(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
	}
}
