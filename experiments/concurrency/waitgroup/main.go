package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go produce(ch, &wg)
	go consume(ch, &wg)

	wg.Wait()
	fmt.Println("end")
}

func produce(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
}

func consume(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Println(v)
	}
}
