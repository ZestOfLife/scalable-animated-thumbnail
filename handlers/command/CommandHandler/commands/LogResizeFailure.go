package commands

type LogResizeFailure struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
