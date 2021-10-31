package main

import (
	"fmt"
	"os"
)

func main() {
	s, sep := "", ""
	// range 会产生一对键值对，索引和这个索引元素的值
	// _ 表示"空标识符"，可以表示任何语法需要变量名但是逻辑不需要的地方
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

/**
for {
	// 函数体 表示无线循环
}

for condition { // go 使用此方式替换 while
	// 循环体
}
*/
