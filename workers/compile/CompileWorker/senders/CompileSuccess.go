package senders

import (
	"CompileWorker/commands"
	"CompileWorker/queue"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func CompileSuccess(BucketID int, VideoName string, FileName string, ExpectedFrames int) {
	postBody, _ := json.Marshal(commands.LogCompileSuccess{BucketID: BucketID, VideoName: VideoName, FileName: FileName, ExpectedFrames: ExpectedFrames})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://event-store:8080/reportcompile", "application/json", responseBody)
	if err != nil {
		log.Println(err)
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "http://event-store:8080/reportcompile"})
		return
	}
	defer resp.Body.Close()

	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Println(err2)
		queue.PushDeadQueue(queue.DeadType{Message: postBody, URI: "http://event-store:8080/reportcompile"})
		return
	}
}
