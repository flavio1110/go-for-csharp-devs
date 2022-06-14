package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	exit := make(chan os.Signal, 1)
	stopJob := make(chan any)

	// fire job
	go scheduleJob(stopJob)
	// capture signals to interrupt and terminate
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Working til shutdown")
	<-exit                // blocked until signal is received
	stopJob <- struct{}{} // send a message to stop the job
	fmt.Println("Graceful shutdown...")
	time.Sleep(time.Second) // wait a second to wait the job to complete
}

func scheduleJob(stop chan any) {
	fmt.Println("Starting job!")
	defer fmt.Println("Exiting job!")
L:
	for {
		select {
		case <-time.After(time.Duration(time.Millisecond * 200)):
			fmt.Println("doing some work")
		case <-stop:
			fmt.Println("stop received. I can finally rest")
			break L // without the label it will break the select statement
		}
	}
}
