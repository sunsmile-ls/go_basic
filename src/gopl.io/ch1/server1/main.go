package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)
// 防止并发的读写
var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	// %q 带引号的字符串
	fmt.Fprintf(w, "URL.Path = %q \n", req.URL.Path)
}

func counter(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "count %d \n", count)
	mu.Unlock()
}