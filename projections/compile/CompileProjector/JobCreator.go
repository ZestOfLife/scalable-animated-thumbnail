package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v9"
)

func queueJob(wg *sync.WaitGroup, client *redis.Client, BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	cntx := context.Background()
	payload, _ := json.Marshal(JobCertificate{BucketID: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	client.LPush(cntx, "queue", string(payload))
	wg.Done()
}

func redoJob(wg *sync.WaitGroup, client *redis.Client, BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	cntx := context.Background()
	payload, _ := json.Marshal(JobCertificate{BucketID: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	client.LPush(cntx, "queue", string(payload))
	wg.Done()
}
