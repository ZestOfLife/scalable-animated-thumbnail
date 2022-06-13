package commands

type LogExtractFailure struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
