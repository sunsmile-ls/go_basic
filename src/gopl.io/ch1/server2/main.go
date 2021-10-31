package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", req.Method, req.URL, req.Proto )
	for key, val := range req.Header {
		fmt.Fprintf(w, "Hander[%q] = %q\n", key, val)
	}
	fmt.Fprintf(w, "Host = %q\n", req.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", req.RemoteAddr)

	// 解析 请求数据
	if err := req.ParseForm(); err != nil{
		log.Print(err)
	}
	// 打印请求的数据
	for key, val := range req.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", key, val)
	}
}