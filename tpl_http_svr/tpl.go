package main

import (
	"path"
	"os"
)

type Tpl struct {
}

// 取模板所在目录
func TplDir(id string) string {
	return path.Join("tpls", id)
}

// 初始化
func TplInit_(id string)  {
	err := os.MkdirAll(TplDir(id), 0777)
	if err != nil {
		panic(err)
	}
}

// 取para路径
func TplParaPath(id string) string {
	p := path.Join(TplDir(id), "para.yaml")
	return p
}

// 取tpl的路径
func TplPath(id string) string {
	p := path.Join(TplDir(id), "tpl.yaml")
	return p
}
