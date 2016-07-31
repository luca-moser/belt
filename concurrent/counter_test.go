package concurrent

import (
	"testing"
)

const result = 5

func TestCounter(t *testing.T) {
	c := NewCounter()
	defer c.Exit()
	for i := 0; i < result; i++ {
		c.Incr()
	}
	r := c.Result()
	if r != result {
		t.Errorf("result was %d, expected %d", r, result)
	}
}

func TestSimpleCounter(t *testing.T) {
	c := NewSimpleCounter()
	for i := 0; i < 7; i++ {
		c.Incr()
	}
	c.Decr()
	c.Decr()
	if c.Val() != result {
		t.Errorf("result was %d, expected %d", c.Val(), result)
	}
}
