package main

import (
	"fmt"
	"sync"
)

/*
Написать программу, которая конкурентно рассчитает значения квадратов чисел,
взятых из массива [2,4,6,8,10], и выведет результаты в stdout.
Подсказка: запусти несколько горутин, каждая из которых возводит число в квадрат.
*/

func main() {
	inputArr := []int{2, 4, 6, 8, 10}
	channel := make(chan int)

	wg := sync.WaitGroup{}
	for i := range 3 {
		wg.Add(1)
		go func(n int) {
			for {
				number, ok := <-channel
				if number != 0 {
					fmt.Printf("Routine #%d answer: %d\n", n, number*number)
					continue
				}
				if !ok {
					break
				}
			}
			defer wg.Done()
		}(i)
	}

	for _, v := range inputArr {
		channel <- v
	}
	close(channel)

	wg.Wait()
}
