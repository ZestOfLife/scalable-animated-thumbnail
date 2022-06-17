package commands

type LogResizeSuccess struct {
	BucketType
	VideoName      string
	FileName       string
	ExpectedFrames int
}
