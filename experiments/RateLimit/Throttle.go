package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	t := SimpleThrottler{
		interval: time.Second,
	}
	t.Do(func () {
		writeHello(1)
	})
	t.Do(func () {
		writeHello(2)
	})
	time.Sleep(time.Second * 2)
	t.Do(func () {
		writeHello(3)
	})
}

func writeHello(n int){
	fmt.Printf("hello %d %s\n", n, time.Now())
}

type SimpleThrottler struct {
	interval time.Duration
	once sync.Once
	mu sync.Mutex
}

func (t *SimpleThrottler) Do(todo func()) {
	lockAcquired := t.mu.TryLock()
	if !lockAcquired {
		return
	}
	go func() {
		time.Sleep(t.interval)
		t.mu.Unlock()
	}()
	todo()
}
