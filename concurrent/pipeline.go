package concurrent

import (
	"sync"
	"time"
)

// NewPipeline creates a new pipeline
func NewPipeline() *pipeline {
	p := &pipeline{
		make(chan interface{}),
		make(chan struct{}, 1), make(chan struct{}, 1), make(chan struct{}, 1),
		false, false, make([]*pipe, 0),
	}
	return p
}

type pipeline struct {
	output    chan interface{}
	stop      chan struct{}
	pause     chan struct{}
	resume    chan struct{}
	isStopped bool
	isPaused  bool
	pipes     []*pipe
}

// Adds a pipe to the end of the pipeline
func (p *pipeline) AddPipe(name string, buffer int, f func(interface{}) interface{}) *pipe {
	pipe := p.newpipe(name, buffer, f)
	pipe.next = p.output
	if len(p.pipes) > 0 {
		lastPipe := p.pipes[len(p.pipes)-1]
		lastPipe.next = pipe.receive
	}
	p.pipes = append(p.pipes, &pipe)
	return &pipe
}

// Start starts the pipeline for execution.
// It returns a channel on which the result of the end of the pipeline can be received.
func (p *pipeline) Start(input <-chan interface{}) <-chan interface{} {
	for i := range p.pipes {
		p.pipes[i].init()
	}
	go func() {
	exit:
		for {
			select {
			// prioritize stop channel
			case <-p.stop:
				for _, pipe := range p.pipes {
					pipe.stop <- struct{}{}
				}
				p.isStopped = true
				break exit
			default:
				// pause, resume or send input to first channel
				select {
				case <-p.pause:
					if p.isStopped || p.isPaused {
						panic("can't pause: pipeline is stopped")
					}
					for _, pipe := range p.pipes {
						pipe.pause <- struct{}{}
					}
					p.isPaused = true
				case <-p.resume:
					if p.isStopped {
						panic("can't resume: pipeline is stopped")
					}
					if !p.isPaused {
						continue
					}
					for _, pipe := range p.pipes {
						pipe.resume <- struct{}{}
					}
					p.isPaused = false
				case p.pipes[0].receive <- <-input:
				}
			}
		}
	}()
	return p.output
}

// Stop stops the execution of the pipeline.
// After calling this method, the pipeline becomes unusable.
func (p *pipeline) Stop() {
	p.stop <- struct{}{}
}

// Pause pauses the execution of the pipeline until Resume() is called.
// This function panics if the pipeline is stopped.
func (p *pipeline) Pause() {
	p.pause <- struct{}{}
}

// Resume resumes the pipeline.
// This function panics if the pipeline is stopped (not paused).
func (p *pipeline) Resume() {
	p.resume <- struct{}{}
}

func (p *pipeline) newpipe(name string, buffer int, f func(interface{}) interface{}) pipe {
	pi := pipe{
		f, make(chan interface{}, buffer), nil,
		make(chan struct{}, 1), make(chan struct{}, 1), make(chan struct{}, 1),
		name, false, make(chan pipeexecution, buffer), 0, sync.Mutex{},
	}
	return pi
}

type pipe struct {
	f           func(interface{}) interface{}
	receive     chan interface{}
	next        chan interface{}
	stop        chan struct{}
	pause       chan struct{}
	resume      chan struct{}
	name        string
	measure     bool
	measurement chan pipeexecution
	processed   int64
	mu          sync.Mutex
}

type pipeexecution struct {
	Name      string
	Processed int64
	Result    interface{}
	Delta     time.Duration
}

// Measure returns a buffered channel which receives the results of the pipe and passed worktime.
// The pipe's results are still sent to the next pipe in the pipeline.
func (p *pipe) Measure() <-chan pipeexecution {
	p.mu.Lock()
	p.measure = true
	p.mu.Unlock()
	return p.measurement
}

// StopMeasuring stops the pipe to send results to the channel given by Results()
func (p *pipe) StopMeasuring() {
	p.mu.Lock()
	p.measure = false
	p.mu.Unlock()
}

func (p *pipe) init() {
	go func() {
		for {
			select {
			case <-p.stop:
				break
			case <-p.pause:
				<-p.resume
			case val := <-p.receive:
				var res interface{}
				if p.measure {
					s := time.Now()
					res = p.f(val)
					p.processed++
					p.measurement <- pipeexecution{p.name, p.processed, res, time.Since(s)}
				} else {
					res = p.f(val)
					p.processed++
				}
				p.next <- res
			}
		}
	}()
}
