package main

import (
	"testing"
	"log"
	J "utils/json"
)

func TestStructAssign(t *testing.T) {
	type A struct {
		Name string
		Age int
	}

	type B struct {
		A
		Company string
	}
	a := A{"aa", 13}
	var b B
	b.A = a // 可以将A作为B的一部分赋值，用组合的方式实现类似于继承的功能

	log.Printf("%+v", b)
	log.Printf("%+v", J.ToJson(b))
}

