// 测试命令：
// <url:vimscript:!go run % --call=f shit -ca -cb> 演示长选项名的形式，以及重复选项
// <url:vimscript:!go run % -nf shit> 短选项名
// TODO: 更详细的用法今后再补充

package main

import (
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// 选项，带有"-"前缀
	n = kingpin.Flag("call", "my name").Short('n').String() // 单个值选项
	c = kingpin.Flag("collect", "my collection").Short('c').Strings() // 多个值选项

	// 参数：没有"-"前缀
	arg = kingpin.Arg("arg", "the arg").Required().String()
)

func main() {
	kingpin.Parse()
	log.Printf("%v, %s\n", *n, *arg)
	log.Printf("*c: %+v", *c)
}
