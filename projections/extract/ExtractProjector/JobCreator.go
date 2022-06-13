package main

import (
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v9"
)

func queueJob(wg *sync.WaitGroup, client *redis.Client, BucketID int, VideoName string, ExpectedFrames int, FPS int, DurationAt int) {
	for i := 0; i < ExpectedFrames; i++ {
		payload, _ := json.Marshal(JobCertificate{BucketID: BucketID, VideoName: VideoName, Timestamp: DurationAt + (i * 1 / FPS * 1000)})
		client.LPush("queue", string(payload))
	}
	wg.Done()
}
