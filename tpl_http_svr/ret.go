// 定义返回类型、错误类型
package main

// 错误类型
type Err struct {
	Code int    `json:"code"` // 返回码
	Msg  string `json:"msg"`  // 错误消息
}

func (me *Err) FromStr(s string) {
	me.Msg = s;
	me.Code = -1; // 缺省设为-1吧
}

func (me *Err) FromError(e error) {
	me.Msg = e.Error();
	me.Code = -1; // 缺省设为-1吧
}

func (me *Err) FromPanic(p interface{}) {
	// 先支持 error 格式
	if e, ok := p.(error); ok {
		me.FromError(e)
	} else if e, ok := p.(string); ok {
		me.FromStr(e)
	}
}

