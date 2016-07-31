package concurrent

import "time"

type simpleratelimiter struct {
	rate     int
	duration time.Duration
	limiter  chan struct{}
	exit     chan struct{}
}

func (rl *simpleratelimiter) init() {
	go func() {
	exit:
		for {
			select {
			case <-time.After(rl.duration):
			case <-rl.exit:
				break exit
			}
			for i := 0; i < rl.rate; i++ {
				select {
				case <-rl.limiter:
				case <-rl.exit:
					break exit
				}
			}
		}
	}()
}

// Exit tells the rate limiter to stop listening for any further channel messages.
// Calls to TryPass() after calling this function will result in a deadlock.
func (rl *simpleratelimiter) Exit() {
	rl.exit <- struct{}{}
}

// TryPass tries to get a passthrough during the current cycle and blocks until it can pass.
func (rl *simpleratelimiter) TryPass() {
	rl.limiter <- struct{}{}
}

// NewRateLimier returns a rate limiter.
func NewRateLimier(rate int, duration time.Duration) *simpleratelimiter {
	r := &simpleratelimiter{rate, duration, make(chan struct{}), make(chan struct{})}
	r.init()
	return r
}
