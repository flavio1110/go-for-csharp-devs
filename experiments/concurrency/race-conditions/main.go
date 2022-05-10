package main

import (
	"fmt"
	"sync"
	"time"
)

type resource struct {
	sync.Mutex
	data int
}

func (r *resource) Increment(v int) {
	r.Mutex.Lock()
	r.data += v
	r.Mutex.Unlock()
}

func main() {
	safeWithMutex(100_000)
	safeWithChannel(100_000)
	nonSafe(100_000)
}

func nonSafe(it int) {
	var wg sync.WaitGroup
	var now = time.Now()
	sharedResource := resource{}
	wg.Add(it)

	for i := 0; i < it; i++ {
		go func() {
			defer wg.Done()
			sharedResource.data += 2
			time.Sleep(time.Duration(400 * int(time.Millisecond)))
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(now), "nonsafe", sharedResource.data)
}

func safeWithMutex(it int) {
	var wg sync.WaitGroup
	var now = time.Now()
	sharedResource := resource{}
	wg.Add(it)

	for i := 0; i < it; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(400 * int(time.Millisecond)))
			sharedResource.Increment(2)
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(now), "safe w/ mutex", sharedResource.data)
}

func safeWithChannel(it int) {
	var wg sync.WaitGroup
	var now = time.Now()
	sharedResource := resource{}
	wg.Add(it)

	ch := make(chan int, it)

	for i := 0; i < it; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(400 * int(time.Millisecond)))
			ch <- 2
		}()
	}

	for i := 0; i < it; i++ {
		sharedResource.data += <-ch
	}

	wg.Wait()
	fmt.Println(time.Since(now), "safe w/ ch", sharedResource.data)
}
