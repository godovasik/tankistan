package main

import "time"

type User struct {
	Name       string
	LastUpdate time.Time
	KeepTrack  bool
}
