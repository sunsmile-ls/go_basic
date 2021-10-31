package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	// 记录开始时间
	start := time.Now()
	// 并发获取，需要通过channel
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	// 打印获取到的值
	for range os.Args[1:] { // 此处返回值全部不要
		fmt.Println(<-ch) // 从通道内获取值
	}
	// 获取间隔时间
	fmt.Printf("%.2fs elapsed", time.Since(start).Seconds())
}

func fetch(url string, ch chan<-string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // 发送到通道ch
		return
	}
	// 返回复制的字节数
	// ioutil.Discard 输出流进行丢弃
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // 释放资源
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s elapsed", secs, nbytes, url)
}