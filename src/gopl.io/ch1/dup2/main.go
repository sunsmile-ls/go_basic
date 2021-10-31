package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 引用类型： 指针、slice 切片、 管道 channel、 接口 interface、 map、函数
	// 基本类型： int、float、bool、string、数组、struct
	// 本文件读取方式是流式读取文件内容，适合读取大文件
	counts := make(map[string]int)
	files := os.Args[1:]

	if len(files) == 0 { // 没有传入文件则标准输出
		countLines(os.Stdin, counts)
	} else { // 读取文件
		for _, arg := range files {
			// os.Open 返回两个值 1.打开的文件 2.err。打开成功的话err为 nil
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v", err)
				continue
			}
			countLines(f, counts)
			f.Close() // 关闭文件
		}	
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d \t %s \n", n, line)
		}
	}

}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	// 本示例忽略了input.Scan()产生的错误
}
