package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)

	// 创建了一个扫描器，从程序标准输入中获取内容
	input := bufio.NewScanner(os.Stdin)

	// 读取下一行，并且去掉结尾换行符， Scan 读到新行的时候返回true，没有读取到更多内容的时候返回false
	for input.Scan() {
		if input.Text() == "exit" { // 添加此行是为了跳出循环
			break;
		}
		// input.Text() 获取当前行出入的内容
		counts[input.Text()]++
	}

	for line, n := range counts { // n 表示数量，重复的行数， line 表示当前行的内容
		if n>1 {
			fmt.Printf("%d \t %s", n, line)
		}
	}
}