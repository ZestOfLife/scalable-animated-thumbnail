package commands

type LogVideoJob struct {
	BucketType
	VideoName      string
	ExpectedFrames int
	FPS            int
	DurationAt     int
}
