package senders

import (
	"ExtractWorker/commands"
	"ExtractWorker/queue"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ExtractFailure(BucketID int, VideoName string, FileName string, Timestamp float32, ExpectedFrames int) {
	postBody, _ := json.Marshal(commands.LogExtractFailure{BucketID: BucketID, VideoName: VideoName, FileName: FileName, Timestamp: Timestamp, ExpectedFrames: ExpectedFrames})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://event-store:8080/reportextractfailure", "application/json", responseBody)
	if err != nil {
		log.Println(err)
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "http://event-store:8080/reportextractfailure"})
		return
	}
	defer resp.Body.Close()

	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Println(err2)
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "http://event-store:8080/reportextractfailure"})
		return
	}
}
