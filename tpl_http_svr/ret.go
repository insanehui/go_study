// 定义返回类型、错误类型
package main

// 错误类型
type Err struct {
	Code int    `json:"code"` // 返回码
	Msg  string `json:"msg"`  // 错误消息
}
