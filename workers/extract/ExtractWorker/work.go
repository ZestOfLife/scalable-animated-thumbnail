package main

import (
	"ExtractWorker/queue"
	"ExtractWorker/senders"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

		for i := 0; i < len(res); i++ {
			var job JobCertificate
			err2 := json.Unmarshal([]byte(res[i]), &job)
			if err2 != nil {
				log.Fatal(err2)
			}

			path := filepath.Join(".", fmt.Sprintf("%d", job.BucketID), job.VideoName)
			err3 := os.MkdirAll(path, os.ModePerm)
			if err3 != nil {
				go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
				continue
			}

			err4 := downloader(minioClient, job.BucketID, job.VideoName)
			if err4 != nil {
				go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
				continue
			}

			timeAt := fmt.Sprintf("'%v", job.Timestamp) + "ms'"

			cmd := exec.Command("ffmpeg", "-ss", timeAt, "-i", path, "-frames:v", "1", "-q:v", "2", path+"/"+job.FileName)
			err5 := cmd.Run()
			if err5 != nil {
				go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
				continue
			}

			err6 := uploader(minioClient, job.BucketID, job.VideoName, job.FileName)
			if err6 != nil {
				go senders.ExtractFailure(job.BucketID, job.VideoName, job.FileName, job.Timestamp, job.ExpectedFrames)
				continue
			}
			go senders.ExtractSuccess(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
		}
	}
}
