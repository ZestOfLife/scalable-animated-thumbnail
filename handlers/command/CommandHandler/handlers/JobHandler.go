package handlers

import (
	"CommandHandler/commands"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v9"
)

func JobHandler(w http.ResponseWriter, req *http.Request) {
	var cmd LogVideoJob
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&cmd)
		if err != nil {
			http.Error(w, "Malformed request (check your payload sent to this server)", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Only POST request supported", http.StatusBadRequest)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("Unbale to connect to event store", err)
	}

	err := client.XAdd(&redis.XAddArgs{
		Stream:       "buckets",
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"Event": "NewRequest",
			"BucketID":     cmd.BucketID,
			"VideoName":   cmd.VideoName,
			"ExpectedFrames": cmd.ExpectedFrames
			"FPS": cmd.FPS
			"DurationAt": cmd.DurationAt
		},
	}).Err()
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
}
