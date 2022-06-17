package commands

type LogCompileFailure struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
