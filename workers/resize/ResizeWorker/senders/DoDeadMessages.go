package senders

import (
	"ResizeWorker/queue"
	"bytes"
	"io/ioutil"
	"net/http"
)

func DoDeadMessages() {
	queue.DeadWG.Wait()
	queue.DeadWG.Add(1)
	count := queue.GetDeadLen()
	for i := 0; i < count; i++ {
		msg := queue.PopDeadQueue()
		responseBody := bytes.NewBuffer(msg.Message)
		resp, err := http.Post(msg.URI, "application/json", responseBody)
		if err != nil {
			queue.PushDeadQueue(msg)
			return
		}
		defer resp.Body.Close()

		_, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			queue.PushDeadQueue(msg)
			return
		}
	}
	queue.DeadWG.Done()
}
