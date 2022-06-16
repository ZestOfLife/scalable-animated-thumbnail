package commands

type LogResizeFailure struct {
	BucketID       int
	VideoName      string
	FileName       string
	ExpectedFrames int
}
