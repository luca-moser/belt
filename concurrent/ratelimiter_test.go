package concurrent

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	s := time.Now()
	limiter := NewRateLimier(500, time.Duration(1)*time.Second)
	defer limiter.Exit()
	for i := 0; i < 1000; i++ {
		limiter.TryPass()
	}
	// should be 4, but we allow 100ms tolerance
	if time.Since(s).Seconds() > 2.1 {
		t.Error("rate limiter took too long")
	}
}
