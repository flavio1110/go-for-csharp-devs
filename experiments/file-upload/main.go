package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const MAX_SIZE = 1024 * 1024 * 5 // 5MB
func upload(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Incoming request")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check the total size of the request
	// it's not exactly the size of the file, but prevents parsing the
	// full body to check the size
	log.Default().Println("Read body")
	r.Body = http.MaxBytesReader(w, r.Body, MAX_SIZE)
	log.Default().Println("Read file")
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("File not accepted: %q", err)))
		return
	}
	log.Default().Println("Read completed")
	defer file.Close()
	tmp, err := os.Create(handler.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to create file: %s", err)))
		return
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to copy file: %s", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Uploaded %q - Size: %d KB", handler.Filename, handler.Size/1024)))
}

func main() {
	http.HandleFunc("/upload", upload)
	log.Default().Println("Listing on port :9090")
	go func() {
		if err := http.ListenAndServe("127.0.0.1:9090", nil); err != nil {
			log.Default().Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Default().Println("Serer sttoped")
}
