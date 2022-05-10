package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)
	chProdSignal := produce(ch)
	chConsumeSignal := consume(ch)

	for {
		select {
		case <-chProdSignal:
			chProdSignal = nil
			fmt.Println("consumer finished")
		case <-chConsumeSignal:
			chConsumeSignal = nil
			fmt.Println("producer finished")
		case <-time.After(time.Duration(time.Millisecond * 200)):
			fmt.Println("waiting a bit")
		}

		if chProdSignal == nil && chConsumeSignal == nil {
			break
		}
	}

	fmt.Println("end")
}

func produce(ch chan<- int) chan int {
	signal := make(chan int, 1)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
		signal <- 1
	}()

	return signal
}

func consume(ch <-chan int) chan int {
	signal := make(chan int, 1)

	go func() {
		for v := range ch {
			fmt.Println(v)
			time.Sleep(500 * time.Millisecond)
		}
		signal <- 1
	}()

	return signal
}
