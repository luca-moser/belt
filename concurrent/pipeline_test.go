package concurrent

import (
	"testing"
	"time"
)

func TestPipline(t *testing.T) {
	pipeline := NewPipeline()
	feedChannel := make(chan interface{})
	for i := 0; i < 10; i++ {
		func(i int) {
			pipeline.AddPipe("", 100, func(input interface{}) interface{} {
				input = input.(int) + 1
				return input
			})
		}(i)
	}
	out := pipeline.Start(feedChannel)
	go func() {
		for i := 0; i < 10; i++ {
			feedChannel <- i
		}
	}()
	anyValueWrong := false
	anyPass := false
	go func() {
		for i := 0; i < 10; i++ {
			num := <-out
			if num != i+10 {
				anyValueWrong = true
			}
			anyPass = true
		}
	}()
	<-time.After(time.Duration(1) * time.Second)
	if anyValueWrong || !anyPass {
		t.Fail()
	}
}
