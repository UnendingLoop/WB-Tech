package main

import (
	"fmt"
	"sync"
	"time"
)

func recursiveOr(channels ...<-chan any) <-chan any {
	var or func(channels ...<-chan any) <-chan any // если не объявить отдельно, то компилятор не даст возможность рекурсивного вызова
	// отсеиваем ниловые каналы
	filtered := make([]<-chan any, 0)
	for _, v := range channels {
		if v == nil {
			continue
		}
		filtered = append(filtered, v)
	}

	// запускаем рекурсию
	or = func(chans ...<-chan any) <-chan any {
		switch len(chans) {
		case 0:
			return nil
		case 1:
			return chans[0]
		}
		out := make(chan any)
		go func() {
			select {
			case <-chans[0]:
			case <-or(chans[1:]...):
			}
			defer close(out)
		}()

		return out
	}

	// возврат 'объединенного' канала
	return or(filtered...)
}

func nonRecursiveOr(channels ...<-chan any) <-chan any {
	out := make(chan any)

	// отсеиваем ниловые каналы
	filtered := make([]<-chan any, 0)
	for _, v := range channels {
		if v == nil {
			continue
		}
		filtered = append(filtered, v)
	}

	// готовим 'одноразовую' функцию
	closeOnce := sync.Once{}

	// запускаем цикл
	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		for _, v := range filtered {
			go func(ch <-chan any) {
				<-ch
				closeOnce.Do(func() { close(out) })
			}(v)
		}
	}
	return out
}

func main() {
	sig := func(after time.Duration) chan any {
		c := make(chan any)
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
		nil,
		sig(11*time.Second),
		sig(13*time.Second),
	)
	fmt.Println("Recursive OR signal received at", time.Since(start))

	start = time.Now()
	<-nonRecursiveOr(
		sig(5*time.Second),
		sig(7*time.Second),
		sig(9*time.Second),
		nil,
		sig(11*time.Second),
		sig(13*time.Second),
	)
	fmt.Println("Cyclic OR signal received at", time.Since(start))
}
