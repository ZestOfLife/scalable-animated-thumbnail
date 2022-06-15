package senders

import (
	"ExtractWorker/commands"
	"ExtractWorker/queue"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ExtractSuccess(BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	postBody, _ := json.Marshal(commands.LogExtractSuccess{BucketID: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("event-store:8080/reportextract", "application/json", responseBody)
	if err != nil {
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "event-store:8080/reportextract"})
		return
	}
	defer resp.Body.Close()

	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "event-store:8080/reportextract"})
		return
	}
}
