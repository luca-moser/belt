package tests

import (
	"testing"

	"github.com/luca-moser/belt/concurrent"
)

func TestCounter(t *testing.T) {
	c := concurrent.NewCounter()
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
