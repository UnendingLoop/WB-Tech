package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
Реализовать структуру-счётчик, которая будет инкрементироваться в конкурентной среде (т.е. из нескольких горутин).
По завершению программы структура должна выводить итоговое значение счётчика.

Подсказка: вам понадобится механизм синхронизации, например, sync.Mutex или sync/Atomic для безопасного инкремента.
*/

type secureCounter struct {
	cntAtomic atomic.Int64
	cntInt    int
	sync.Mutex
}

func (s *secureCounter) incrementAtomic() {
	s.cntAtomic.Add(1)
}
func (s *secureCounter) incrementInt() {
	s.Lock()
	s.cntInt++
	s.Unlock()
}
func (s *secureCounter) printValues() {
	fmt.Println("Atomic counter value is:", s.cntAtomic.Load())
	fmt.Println("Integer counter with mutex value is:", s.cntInt)
}

func main() {
	wg := sync.WaitGroup{}
	var counter secureCounter
	defer counter.printValues()

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.incrementAtomic()
			counter.incrementInt()
		}()
	}
	wg.Wait()
}
