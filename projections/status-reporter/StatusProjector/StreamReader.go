// Stream Reader
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/rs/xid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func streamReader(client *redis.Client, fetch int) {
	for fetch != -1 {
		cntx := context.Background()
		consumerId := xid.New().String()

		send_to, conn_err := gorm.Open(mysql.Open(DB_ADDR), &gorm.Config{})
		if conn_err != nil {
			log.Fatal("Unbale to connect to queue", conn_err)
		}

		var config redis.XReadGroupArgs
		if fetch == 1 {
			config = redis.XReadGroupArgs{
				Group:    "observer",
				Consumer: consumerId,
				Streams:  []string{"buckets", "0-0"},
				NoAck:    false,
			}
			fetch = -1
		} else {
			config = redis.XReadGroupArgs{
				Group:    "observer",
				Consumer: consumerId,
				Block:    0,
				Streams:  []string{"buckets", ">"},
				NoAck:    false,
			}
		}
		res, err := client.XReadGroup(cntx, &config).Result()
		if err != nil {
			log.Fatal(err)
		}

		buckets := make(map[int]Store)
		for i := 0; i < len(res[0].Messages); i++ {
			msgID := res[0].Messages[i].ID
			val := res[0].Messages[i].Values
			Event := fmt.Sprintf("%v", val["Event"])
			BucketID, _ := strconv.Atoi(fmt.Sprintf("%v", val["BucketID"]))
			VideoName := fmt.Sprintf("%v", val["VideoName"])
			ExpectedFrames, _ := strconv.Atoi(fmt.Sprintf("%v", val["ExpectedFrames"]))
			switch Event {
			case "NewRequest":
				buckets[BucketID] = Store{BucketID: BucketID, VideoName: VideoName, ExpectedFrames: ExpectedFrames, Extracted: 0, Resized: 0, Compiled: 0}
				rese := send_to.Create(&Store{BucketID: BucketID, VideoName: VideoName, ExpectedFrames: ExpectedFrames, Extracted: 0, Resized: 0, Compiled: 0})
				if rese.Error != nil {
					log.Fatal(rese.Error)
				}
			case "FrameExtracted":
				entry, a := buckets[BucketID]
				if a {
					entry.Extracted += 1
					buckets[BucketID] = entry
				} else {
					var data = Store{BucketID: BucketID, VideoName: VideoName}
					send_to.First(&data)
					data.Extracted += 1
					buckets[BucketID] = data
				}
			case "FrameResized":
				entry, a := buckets[BucketID]
				if a {
					entry.Resized += 1
					buckets[BucketID] = entry
				} else {
					var data = Store{BucketID: BucketID, VideoName: VideoName}
					send_to.First(&data)
					data.Resized += 1
					buckets[BucketID] = data
				}
			case "FrameCompiled":
				entry, a := buckets[BucketID]
				if a {
					entry.Compiled += 1
					buckets[BucketID] = entry
				} else {
					var data = Store{BucketID: BucketID, VideoName: VideoName}
					send_to.First(&data)
					data.Compiled += 1
					buckets[BucketID] = data
				}
			}
			client.XAck(cntx, "buckets", "observer", msgID)
		}
		for _, ele := range buckets {
			send_to.Save(&ele)
		}
		db, close_err := send_to.DB()
		if close_err != nil {
			log.Fatal(close_err)
		}
		db.Close()
	}
}
