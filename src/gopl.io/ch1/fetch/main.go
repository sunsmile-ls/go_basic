package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)
// go语 fmt.Print 系列文章比较 https://juejin.cn/post/6948775660563202085
func main() {
	for _, url := range os.Args[1:] {
		// 返回可读取的数据流
		resp, err := http.Get(url)
		if err != nil {
			// Stderr是无缓冲，每个输出都会立即flush，Stdout是行缓冲的，要等到缓冲满了才flush,
			// 前者更符合作为日志的需要，不然你程序执行过程中core了，
			// 缓冲里的遗言可能就丢了，而丢掉的往往是最接近出问题的地方的。
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		// 读取整个响应结果到B
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // 释放资源，方式资源泄漏
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}
