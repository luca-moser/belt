package concurrent

import (
	"testing"
)

func TestCounter(t *testing.T) {
	c := NewCounter()
	defer c.Exit()
	const result = 5
	for i := 0; i < result; i++ {
		c.Incr()
	}
	r := c.Result()
	if r != result {
		t.Errorf("result was %d, expected %d", r, result)
	}
}
