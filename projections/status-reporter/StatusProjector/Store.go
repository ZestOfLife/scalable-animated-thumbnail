package main

type Store struct {
	BucketID       int    `gorm:"primaryKey"`
	VideoName      string `gorm:"primaryKey"`
	ExpectedFrames int
	Extracted      int
	Resized        int
	Compiled       int
}
