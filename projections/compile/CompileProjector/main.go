// server
package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"
)

func main() {
	cntx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "event-store:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(cntx).Result()
	if err != nil {
		log.Fatal("Unable to connect to event store: ", err)
	}
	log.Println("Connected to event store")

	err = client.XGroupCreateMkStream(cntx, "buckets", "resize-observer", "0").Err()
	if err != nil {
		log.Println(err) // Possible that it was created beforehand
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go streamReader(client)
	wg.Wait()
}
