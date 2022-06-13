// Stream Reader
package main

import (
	"fmt"
	"encoding/json"
	"sync"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/rs/xid"
)

type empty {}

func streamReader(client *redis.Client) {
	for {
		consumerId := xid.New().String()
		res, err := client.XRead(&redis.XReadArgs{
			Group: "job-observer",
			Consumer: consumerId,
			Streams: []string{"buckets", ">"},
			Block:   0,
			NoAck: false,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		send_to := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB:   0,
		})
		_, conn_err := client.Ping().Result()
		if conn_err != nil {
			log.Fatal("Unbale to connect to queue", err)
		}
		
		var wg2 sync.WaitGroup
		wg2.Add(len(res[0].Messages))
		for i := 0; i < len(res[0].Messages); i++ {
			msgID := res[0].Messages[i].ID
			val := res[0].Messages[i].Values
			Event := fmt.Sprintf("%v", val["Event"])
			if event == "NewRequest" {
				BucketID := fmt.Sprintf("%d", val["BucketID"])
				VideoName := fmt.Sprintf("%v", val["VideoName"])
				ExpectedFrames := fmt.Sprintf("%d", val["ExpectedFrames"])
				FPS := fmt.Sprintf("%d", val["FPS"])
				DurationAt := fmt.Sprintf("%d", val["DurationAt"])
				go queueJob(wg2, send_to, BucketID, VideoName, ExpectedFrames, FPS, DurationAt)
			} else {
				wg2.Done()
			}
			client.XAck("buckets", "job-observer", msgID)
		}
		wg2.Wait()
		send_to.Close()
	}
}
