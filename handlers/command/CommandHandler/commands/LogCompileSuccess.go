package commands

type LogCompileSuccess struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
