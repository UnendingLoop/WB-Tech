package main

import (
	"fmt"
	"sync"
)

/*
Разработать конвейер чисел.
Даны два канала:
- в первый пишутся числа x из массива,
- во второй – результат операции x*2.
После этого данные из второго канала должны выводиться в stdout.
То есть, организуйте конвейер из двух этапов с горутинами:
- генерация чисел и
- их обработка.
Убедитесь, что чтение из второго канала корректно завершается.
*/

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
	wg := sync.WaitGroup{}
	channel1 := make(chan int)
	channel2 := make(chan int)

	//producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(channel1)
		for _, v := range input {
			channel1 <- v
		}
	}()

	//consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(channel2)
		for i := range channel1 {
			channel2 <- i * 2
		}
	}()

	//printer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range channel2 {
			fmt.Println(i)
		}
	}()
	wg.Wait()
}
