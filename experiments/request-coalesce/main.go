package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var processMap = sync.Map{}
var cacheMap = sync.Map{}
var ttlDuration = time.Second * 5

func main() {
	processMap = sync.Map{}
	cacheMap = sync.Map{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	fmt.Println("listening")
	http.ListenAndServe(":8080", mux)
}

func index(res http.ResponseWriter, req *http.Request) {

	start := time.Now()

	key := req.URL.Query().Get("key")

	cached, entry := tryGetFromCache(key)

	if cached {
		writeResult(res, start, key, entry.result)
		return
	}

	result := doSomeWork(key)
	cacheMap.LoadOrStore(key, cacheEntry{
		addedOn: time.Now(),
		key:     key,
		result:  result,
	})

	writeResult(res, start, key, result)
}

func writeResult(res http.ResponseWriter, start time.Time, key string, result string) {
	res.Write([]byte(fmt.Sprintf("served: \t %q \t result: \t %s \t duration: \t %v\n", key, result, time.Since(start))))
}

func tryGetFromCache(key string) (bool, cacheEntry) {
	v, ok := cacheMap.Load(key)

	if !ok {
		return false, cacheEntry{}
	}

	entry := v.(cacheEntry)

	if time.Since(entry.addedOn) > ttlDuration {
		cacheMap.Delete(key)
		return false, cacheEntry{}
	}

	return true, entry
}

func doSomeWork(key string) string {
	defer processMap.Delete(key)
	v, _ := processMap.LoadOrStore(key, startNewProcess(key))

	p := v.(*process)

	p.wait()

	return p.result
}

func startNewProcess(key string) *process {
	p := &process{
		key:  key,
		done: make(chan struct{}),
	}

	go p.doStuffInBackground()

	return p

}

type cacheEntry struct {
	addedOn time.Time
	key     string
	result  string
}

type process struct {
	key       string
	result    string
	done      chan struct{}
	completed bool
}

func (p *process) wait() {
	if p.completed {
		return
	}

	<-p.done
	p.completed = true
}

func (p *process) doStuffInBackground() {
	time.Sleep(time.Second * 2)
	p.result = fmt.Sprintf("result of %q", p.key)

	close(p.done)
}
