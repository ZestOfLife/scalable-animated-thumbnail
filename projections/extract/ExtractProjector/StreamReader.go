// Stream Reader
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/rs/xid"
)

func streamReader(client *redis.Client) {
	for {
		cntx := context.Background()
		consumerId := xid.New().String()

		send_to := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		_, conn_err := client.Ping(cntx).Result()
		if conn_err != nil {
			log.Fatal("Unbale to connect to queue", conn_err)
		}

		res, err := client.XReadGroup(cntx, &redis.XReadGroupArgs{
			Group:    "job-observer",
			Consumer: consumerId,
			Streams:  []string{"buckets", ">"},
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		var wg2 sync.WaitGroup
		wg2.Add(len(res[0].Messages))
		for i := 0; i < len(res[0].Messages); i++ {
			msgID := res[0].Messages[i].ID
			val := res[0].Messages[i].Values
			Event := fmt.Sprintf("%v", val["Event"])
			if Event == "NewRequest" {
				BucketID, _ := strconv.Atoi(fmt.Sprintf("%v", val["BucketID"]))
				VideoName := fmt.Sprintf("%v", val["VideoName"])
				ExpectedFrames, _ := strconv.Atoi(fmt.Sprintf("%v", val["ExpectedFrames"]))
				FPS, _ := strconv.Atoi(fmt.Sprintf("%v", val["FPS"]))
				DurationAt, _ := strconv.Atoi(fmt.Sprintf("%v", val["DurationAt"]))
				go queueJob(&wg2, send_to, BucketID, VideoName, ExpectedFrames, FPS, DurationAt)
			} else if Event == "FrameExtractionFailure" {
				BucketID, _ := strconv.Atoi(fmt.Sprintf("%v", val["BucketID"]))
				VideoName := fmt.Sprintf("%v", val["VideoName"])
				FileName := fmt.Sprintf("%v", val["FileName"])
				Timestamp, _ := strconv.ParseFloat(fmt.Sprintf("%v", val["Timestamp"]), 32)
				ExpectedFrames, _ := strconv.Atoi(fmt.Sprintf("%v", val["ExpectedFrames"]))
				go redoJob(&wg2, send_to, BucketID, VideoName, float32(Timestamp), FileName, ExpectedFrames)
			} else {
				wg2.Done()
			}
			client.XAck(cntx, "buckets", "job-observer", msgID)
		}
		wg2.Wait()
		send_to.Close()
	}
}
