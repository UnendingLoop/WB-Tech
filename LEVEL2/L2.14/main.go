package main

import (
	"fmt"
	"sync"
	"time"
)

func recursiveOr(channels ...<-chan interface{}) <-chan interface{} {
	var or func(channels ...<-chan interface{}) <-chan interface{} // если не объявить отдельно, то компилятор не даст возможность рекурсивного вызова

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		out := make(chan interface{})
		go func() {
			select {
			case <-channels[0]:
			case <-or(channels[1:]...):
			}
			defer close(out)
		}()

		return out
	}
	return or(channels...)
}

func nonRecursiveOr(channels ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	closeOnce := sync.Once{}
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	default:
		for _, v := range channels {
			go func(ch chan interface{}) {
				<-ch
				closeOnce.Do(func() { close(out) })
			}(v)
		}
	}
	return out
}

func main() {
	sig := func(after time.Duration) chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-recursiveOr(
		sig(5*time.Second),
		sig(7*time.Second),
		sig(9*time.Second),
		sig(11*time.Second),
		sig(13*time.Second),
	)
	fmt.Println("Recursive OR signal received at", time.Since(start))

	start = time.Now()
	<-nonRecursiveOr(
		sig(5*time.Second),
		sig(7*time.Second),
		sig(9*time.Second),
		sig(11*time.Second),
		sig(13*time.Second),
	)
	fmt.Println("Cyclic OR signal received at", time.Since(start))
}
