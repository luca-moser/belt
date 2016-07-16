package debounce

import (
	"sync"
	"time"
)

// Debounce returns a debounced version of the given function.
// A debounced function only executes if the function isn't called again
// for the specified debounce time.
func Debounce(f func(), debounce int) func() {
	var mtx sync.Mutex
	var callers int
	return func() {
		mtx.Lock()
		callers++
		mtx.Unlock()

		<-time.After(time.Duration(debounce) * time.Millisecond)

		mtx.Lock()
		callers--
		if callers == 0 {
			f()
		}
		mtx.Unlock()
	}
}
