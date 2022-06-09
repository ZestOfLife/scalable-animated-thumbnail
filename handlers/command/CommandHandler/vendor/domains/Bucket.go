package domains

import (
	"container/list"
)

type Bucket struct {
	id     string
	videos list.List // VideoFile list type
}
