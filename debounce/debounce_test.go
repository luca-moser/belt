package debounce

import (
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	timesExecuted := 0
	debFunc := Debounce(func() {
		timesExecuted++
	}, 600)

	for x := 0; x < 50; x++ {
		go debFunc()
	}
	<-time.After(time.Duration(650) * time.Millisecond)

	if timesExecuted != 1 {
		t.Errorf("result was %d, expected %d", timesExecuted, 1)
	}
}
