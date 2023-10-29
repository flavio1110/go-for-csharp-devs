package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	chPayload := make(chan payload, 1)

	go startWorkers(ctx, 1, chPayload)
	go listen(ctx, chPayload)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	x := <-ch
	fmt.Println(x)

	fmt.Println("Shuting down...")
	cancel()

	fmt.Println("Closing connections...")
	time.Sleep(time.Second * 2)

	fmt.Println("bye")
}

func listen(ctx context.Context, workToDo chan<- payload) {
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		log.Fatal("listen to localhost:9999")
	}

	conns := 1

	for {
		fmt.Println("waiting for connection...")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("fail to accept connection")
		}

		if ctx.Err() != nil {
			fmt.Println("connection skipped. Application is shutting down.")
			close(workToDo)
			break
		}

		fmt.Println("connection accepted")

		workToDo <- payload{
			id:   conns,
			conn: conn,
		}
		conns++
	}
}

func exchange(ctx context.Context, id int, conn net.Conn) {
	for i := 1; i < 10; i++ {
		if ctx.Err() != nil {
			writeToConnection("early closing. Bye o/", conn)
			break
		}

		if err := writeToConnection(fmt.Sprintf("%d - %s", id, time.Now().String()), conn); err != nil {
			break
		}
		time.Sleep(time.Second)
	}

	_ = writeToConnection("all done", conn)
	if err := conn.Close(); err != nil {
		fmt.Println("fail  to close connection", err)
	}
}

func writeToConnection(msg string, conn net.Conn) error {
	_, err := fmt.Fprintln(conn, msg)
	if err != nil {
		fmt.Println("fail  to write to connection", err)
		return err
	}

	return nil
}

func startWorkers(ctx context.Context, workersLimit int, workToDo <-chan payload) {
	for i := 0; i < workersLimit; i++ {
		go func() {
			for data := range workToDo {
				exchange(ctx, data.id, data.conn)
			}
		}()
	}
}

type payload struct {
	id   int
	conn net.Conn
}
