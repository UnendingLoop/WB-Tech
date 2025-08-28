package main

import (
	"fmt"
	"sync"
)

/*Реализовать безопасную для конкуренции запись данных в структуру map.
Подсказка: необходимость использования синхронизации (например, sync.Mutex или встроенная concurrent-map).
Проверьте работу кода на гонки (util go run -race).*/

type mapa struct {
	Map map[int]struct{}
	sync.RWMutex
}

func main() {
	var secureMap sync.Map
	secureMapWithMTX := mapa{Map: make(map[int]struct{})}
	wg := sync.WaitGroup{}

	for i := range 5 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Printf("Goroutine #%d started working...\n", n)
			start := n * 1000
			end := start + 1000
			for v := start; v <= end; v++ {
				secureMap.Store(v, struct{}{})

				secureMapWithMTX.Lock()
				secureMapWithMTX.Map[v] = struct{}{}
				secureMapWithMTX.Unlock()
			}
			fmt.Printf("Goroutine #%d finished. Shutting down.\n", n)
		}(i)
	}
	wg.Wait()
	fmt.Println("App exiting. Bye")
}
