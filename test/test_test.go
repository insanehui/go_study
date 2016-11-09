package main

import (
	"testing"
	"log"
	J "utils/json"
	V "github.com/asaskevich/govalidator"
	"path"
)

func init() {  
	V.SetFieldsRequiredByDefault(true)
}

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

func TestValidator(t *testing.T){
	{
		log.Println("=========================")
		var q  = struct {
			TplId string `valid:"length(1|5)"`// 模板id
		}{ "" }

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}

	{
		log.Println("=========================")
		var q  = struct {
			TplId string `valid:"length(0|3)"`// 模板id
		}{ "haha" }

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}
}

func Test_nil(t *testing.T){
	var s []string // slice的初始状态为nil
	log.Println(s == nil)
}

func Test_path(t *testing.T){
	log.Println(path.Join("a", "b", "c"))
}

func Test_mixin(t *testing.T){
	type A struct {
		I int
	}

	type B struct {
		J int
	}

	type C struct {
		Aa A 
		B
	}

	var c C
	c.Aa.I = 1 // 这里不能用 c.A.i
	c.B.J = 2 // 这里只能用c.B.i
	log.Printf("%+v", c)
	log.Printf("%+v", J.ToJson(c))
}
