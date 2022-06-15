package queue

import (
	"sync"
)

var DEAD_MESSAGE_QUEUE = make([]DeadType, 0)
var DeadWG sync.WaitGroup

func PushDeadQueue(message DeadType) {
	DEAD_MESSAGE_QUEUE = append(DEAD_MESSAGE_QUEUE, message)
}

func PopDeadQueue() DeadType {
	removed := DEAD_MESSAGE_QUEUE[0]
	DEAD_MESSAGE_QUEUE = DEAD_MESSAGE_QUEUE[1:]
	return removed
}

func GetDeadLen() int {
	return len(DEAD_MESSAGE_QUEUE)
}
