package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

var filePath = "large.csv"

func main() {
	opt := flag.String("o", "read", "read|readall|crate")

	flag.Parse()

	if *opt == "read" {
		readBigFile()
		return
	} else if *opt == "readall" {
		readAllBigFile()
		return
	} else if *opt == "write" {
		createBigFile()
		return
	} else {
		fmt.Println("Option not supported")
	}

}

func readBigFile() {
	fmt.Println("reading file...")
	PrintMemUsage()
	readFileByLine()
	PrintMemUsage()
}

func readAllBigFile() {
	fmt.Println("reading file...")
	PrintMemUsage()
	readAll()
	PrintMemUsage()
}

var allFile string

func readAll() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open file", err)
	}
	defer file.Close()
	t, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("read all file", err)
	}
	allFile = string(t)
}

func readFileByLine() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open file", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	ln := 0
	for {
		ln++
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("a real error happened here: %v\n", err)
		}

		handle(string(line))

		if ln%100_000 == 0 {
			fmt.Printf("%d processed\n", ln)
			PrintMemUsage()
			fmt.Println("=======================")
		}

	}
}

var text []string

func handle(line string) {
	text = append(text, line)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Printf("\tFrees = %v\n", m.Frees)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func createBigFile() {
	fmt.Println("writing file...")
	header := "id,name;age;country;id,name;age;country;id,name;age;country"

	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal("create file: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, header); err != nil {
		log.Fatal("write header", err)
	}

	for i := 0; i < 100_000_000; i++ {
		line := fmt.Sprintf("%d;josse;%d;NL;%d;josse;%d;NL;%d;josse;%d;NL", i, i+10, i, i+10, i, i+10)
		if _, err := fmt.Fprintln(f, line); err != nil {
			log.Fatal("write line", err)
		}
	}
}
