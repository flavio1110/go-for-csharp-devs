package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	maxProcs := runtime.GOMAXPROCS(0)

	fmt.Printf("Current GOMAXPROCS: %d\n", maxProcs)

	now := time.Now()
	timeout := time.Millisecond * 10000
	numWorkers := 8

	w := newWorkerPool(numWorkers, timeout)
	w.start()

	for i := 0; i < 16; i++ {
		_i := i
		w.enqueue(func(ctx context.Context) error {
			return aWork(ctx, _i)
		})
	}

	w.stop()
	fmt.Println(time.Since(now))
}

func aWork(ctx context.Context, id int) error {
	now := time.Now()

	for i := 0; i < 2_000_000_000; i++ {
		//busy waiting
	}
	fmt.Printf("=====>#%d - waited for %s\n", id, time.Since(now).String())

	return nil
}

type work func(ctx context.Context) error

type workerPool struct {
	queue              chan work
	stopSignal         chan struct{}
	workerssDoneSignal chan struct{}
	maxConcurrency     int
	timeOut            time.Duration
}

func newWorkerPool(maxConcurrency int, timeOut time.Duration) *workerPool {
	return &workerPool{
		queue:              make(chan work),
		stopSignal:         make(chan struct{}),
		workerssDoneSignal: make(chan struct{}, 2),
		maxConcurrency:     maxConcurrency,
		timeOut:            timeOut,
	}
}

func (w *workerPool) start() {
	for i := 0; i < w.maxConcurrency; i++ {
		_i := i
		go func() {
			w.runWorker(_i)
			w.workerssDoneSignal <- struct{}{}
		}()
	}
}

func (w *workerPool) runWorker(id int) {
	for {
		select {
		case todo, ok := <-w.queue:
			if !ok {
				fmt.Printf("Worker #%d - Done\n", id)
				return
			}

			w.doWork(todo)
		case <-w.stopSignal:
			fmt.Printf("Worker #%d - Forced stop\n", id)
			return
		}
	}
}

func (w *workerPool) doWork(todo work) {
	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		todo(ctx)
		close(done)
	}()

	select {
	case <-time.After(w.timeOut):
		fmt.Println("work timeout")
	case <-done:
	}

}

func (w *workerPool) enqueue(todo work) error {
	w.queue <- todo
	return nil
}

func (w *workerPool) stop() error {
	close(w.queue)
	fmt.Println("stopping...")

	fmt.Println("waiting for all the work to finish")
	w.waitForWorkers()

	fmt.Println("all work completed")
	return nil
}

func (w *workerPool) waitForWorkers() {
	timeoutSignal := time.After(w.timeOut)

	for i := 0; i < w.maxConcurrency; i++ {
		select {
		case <-timeoutSignal:
			close(w.stopSignal)
			return
		case <-w.workerssDoneSignal:
		}
	}
	close(w.workerssDoneSignal)
}
