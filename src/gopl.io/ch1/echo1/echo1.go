package main

// gofmt 会按照字母顺序排序倒入的包
import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	// := 短变量声明，只可用于方法内。会根据初始化值给予合适的类型
	for i := 1; i < len(os.Args); i++ { // 左大括号必须和 i++ 在同一行
		s += sep + os.Args[i] // += 是赋值操作符
		sep = " "
	}
	fmt.Println(s)
}
