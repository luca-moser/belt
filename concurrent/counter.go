package concurrent

import "sync"

// counter implements a thread-safe counter
type counter struct {
	val  int64
	incr chan struct{}
	decr chan struct{}
	exit chan struct{}
	res  chan int64
}

func (c *counter) init() {
	go func() {
	exit:
		for {
			select {
			case <-c.incr:
				c.val++
			case <-c.decr:
				c.val--
			case <-c.exit:
				break exit
			case c.res <- c.val:
			}
		}
	}()
}

func (c *counter) Incr() {
	c.incr <- struct{}{}
}

func (c *counter) Decr() {
	c.decr <- struct{}{}
}

func (c *counter) Exit() {
	c.exit <- struct{}{}
}

func (c *counter) Result() int64 {
	return <-c.res
}

// NewCounter creates a new counter
func NewCounter() *counter {
	c := &counter{
		incr: make(chan struct{}),
		decr: make(chan struct{}),
		res:  make(chan int64),
		exit: make(chan struct{}),
	}
	c.init()
	return c
}

type simplecounter struct {
	sync.Mutex
	num int64
}

func (sc *simplecounter) Incr() {
	sc.Lock()
	sc.num++
	sc.Unlock()
}

func (sc *simplecounter) Decr() {
	sc.Lock()
	sc.num--
	sc.Unlock()
}

func (sc *simplecounter) Val() int64 {
	return sc.num
}

func NewSimpleCounter() *simplecounter {
	return &simplecounter{}
}
