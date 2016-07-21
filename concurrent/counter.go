package concurrent

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
		for {
			select {
			case <-c.incr:
				c.val++
			case <-c.decr:
				c.val--
			case <-c.exit:
				break
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
