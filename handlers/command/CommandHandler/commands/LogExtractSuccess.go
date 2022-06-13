package commands

type LogExtractSuccess struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
