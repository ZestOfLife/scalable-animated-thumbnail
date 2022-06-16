package commands

type LogExtractFailure struct {
	BucketID       int
	VideoName      string
	FileName       string
	Timestamp      float32
	ExpectedFrames int
}
