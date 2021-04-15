package util

import (
	"time"
)

// NewTicker ticker for f
func NewTicker(f func(), t time.Duration) {
	ticker := time.NewTicker(t)
	for range ticker.C {
		f()
	}
}

// NewTimer timer for f
func NewTimer(f func(), t time.Duration) {
	timer := time.NewTimer(t)
	<-timer.C
	f()
}
