package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// 此方法不用像 echo2 中那样 每次通过旧的字符串、空格和下一个参数生成新的字符串，旧的字符串被垃圾回收
	fmt.Println(strings.Join(os.Args[1:], " "))
}
