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
			Group:    "resize-observer",
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
			if Event == "FrameResized" {
				BucketID, _ := strconv.Atoi(fmt.Sprintf("%v", val["BucketID"]))
				VideoName := fmt.Sprintf("%v", val["VideoName"])
				FileName := fmt.Sprintf("%v", val["FileName"])
				ExpectedFrames, _ := strconv.Atoi(fmt.Sprintf("%v", val["ExpectedFrames"]))
				go queueJob(&wg2, send_to, BucketID, VideoName, FileName, ExpectedFrames)
			} else if Event == "CompileFailure" {
				BucketID, _ := strconv.Atoi(fmt.Sprintf("%v", val["BucketID"]))
				VideoName := fmt.Sprintf("%v", val["VideoName"])
				FileName := fmt.Sprint("%v", val["FileName"])
				ExpectedFrames, _ := strconv.Atoi(fmt.Sprintf("%v", val["ExpectedFrames"]))
				go redoJob(&wg2, send_to, BucketID, VideoName, FileName, ExpectedFrames)
			} else {
				wg2.Done()
			}
			client.XAck(cntx, "buckets", "resize-observer", msgID)
		}
		wg2.Wait()
		send_to.Close()
	}
}
