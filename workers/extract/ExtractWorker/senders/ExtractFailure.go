package senders

import (
	"ExtractWorker/commands"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ExtractFailure(BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	postBody, _ := json.Marshal(commands.LogExtractFailure{BucketId: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("event-store:8080/reportextractfailure", "application/json", responseBody)
	if err != nil {
		queue.PushDeadQueue(queue.DeadType{message: postBody, uri: "event-store:8080/reportextractfailure"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		queue.PushDeadQueue(queue.DeadType{message: postBody, uri: "event-store:8080/reportextractfailure"})
		return
	}
}
