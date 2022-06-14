package commands

type LogExtractFailure struct {
	BucketType
	VideoName      string
	FileName       string
	Timestamp      float32
	ExpectedFrames int
}
