package senders

import (
	"ExtractWorker/queue"
	"sync"
)

func DoDeadMessages() {
	queue.DeadWG.Wait()
	queue.DeadWG.Add(1)
	count := queue.GetDeadLen()
	for i := 0; i < count; i++ {
		msg := queue.PopDeadQueue()
		responseBody := bytes.NewBuffer(msg.message)
		resp, err := http.Post(msg.uri, "application/json", responseBody)
		if err != nil {
			queue.PushDeadQueue(msg)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			queue.PushDeadQueue(msg)
			return
		}
	}
	queue.DeadWG.Done()
}
