package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v9"
)

func queueJob(wg *sync.WaitGroup, client *redis.Client, BucketID int, VideoName string, ExpectedFrames int, FPS int, DurationAt int) {
	cntx := context.Background()
	for i := 0; i < ExpectedFrames; i++ {
		FileName := fmt.Sprintf("%04d.jpeg", i+1)
		payload, _ := json.Marshal(JobCertificate{BucketID: BucketID, VideoName: VideoName, Timestamp: float32(DurationAt) + (float32(i) * 1.0 / float32(FPS) * 1000.0), FileName: FileName, ExpectedFrames: ExpectedFrames})
		client.LPush(cntx, "queue", string(payload))
	}
	wg.Done()
}

func redoJob(wg *sync.WaitGroup, client *redis.Client, BucketID int, VideoName string, Timestamp float32, FileName string, ExpectedFrames int) {
	cntx := context.Background()
	payload, _ := json.Marshal(JobCertificate{BucketID: BucketID, VideoName: VideoName, Timestamp: Timestamp, FileName: FileName, ExpectedFrames: ExpectedFrames})
	client.LPush(cntx, "queue", string(payload))
	wg.Done()
}
