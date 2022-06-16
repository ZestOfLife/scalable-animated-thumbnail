package senders

import (
	"CompileWorker/commands"
	"CompileWorker/queue"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CompileSuccess(BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	postBody, _ := json.Marshal(commands.LogCompileSuccess{BucketID: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("event-store:8080/reportcompile", "application/json", responseBody)
	if err != nil {
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "event-store:8080/reportcompile"})
		return
	}
	defer resp.Body.Close()

	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "event-store:8080/reportcompile"})
		return
	}
}
