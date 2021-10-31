package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		// 检查前缀
		if !strings.HasPrefix(url, "http:") {
			fmt.Fprintf(os.Stderr, "fetch:请输入 http:// 前缀")
			os.Exit(1)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n",err)
			os.Exit(1)
		}
		fmt.Println(resp.Status)
		// 从 resp.Body 写到 os.Stdout中
		_, err= io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %s: %v \n",url, err)
			os.Exit(1)
		}
	}
}
