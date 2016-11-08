package main

import (
	"testing"
	"log"
	J "utils/json"
	V "github.com/asaskevich/govalidator"
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

