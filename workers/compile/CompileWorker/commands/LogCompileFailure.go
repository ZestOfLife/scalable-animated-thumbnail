package commands

type LogCompileFailure struct {
	BucketID       int
	VideoName      string
	FileName       string
	ExpectedFrames int
}
