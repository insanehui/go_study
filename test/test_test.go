package main

import (
	"html/template"
	"log"
	"os"
	"strconv"
	"testing"
	J "utils/json"
	// Y "utils/yaml"
	"path"

	V "github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
)

func init() {
	V.SetFieldsRequiredByDefault(true)
}

func TestStructAssign(t *testing.T) {
	type A struct {
		Name string
		Age  int
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

func TestValidator(t *testing.T) {
	{
		log.Println("=========================")
		var q = struct {
			TplId string `valid:"length(1|5)"` // 模板id
		}{""}

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}

	{
		log.Println("=========================")
		var q = struct {
			TplId string `valid:"length(0|3)"` // 模板id
		}{"haha"}

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}
}

func Test_nil(t *testing.T) {
	var s []string // slice的初始状态为nil
	log.Println(s == nil)
}

func Test_path(t *testing.T) {
	log.Println(path.Join("a", "b", "c"))
}

func Test_mixin(t *testing.T) {
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
	c.B.J = 2  // 这里只能用c.B.i
	log.Printf("%+v", c)
	log.Printf("%+v", J.ToJson(c))
}

func Test_yaml2tpl(t *testing.T) {
	str := `
Name:    1
Age:     30
Company: 2
`
	// var m map[string]interface{} // map的值为 interface{}的话，是不能对if生效的
	var m map[string]int // 定成具体类型之后，才可以

	yaml.Unmarshal([]byte(str), &m)
	log.Printf("%+v", m)

	tmpl, _ := template.New("xx").Parse(`

		my name is {{.Name}}
		i'm {{.Age}} years old
		{{if gt .Age 10}}
		i'm an adult
		{{else}}
		i'm a child
		{{end}}


		`)

	tmpl.Execute(os.Stdout, m)
}

func Test_booldef(t *testing.T) {
	var a bool
	log.Println(a == false)
	log.Println(a == true)
}

func Test_interface(t *testing.T) {
	var a interface{}
	a = false
	log.Println(a == nil)
	a = 1
	log.Println(a == true)
	log.Println(a.(int))
}

func Test_make(t *testing.T) {
	// var m map[string]int
	// m["aa"] = 1 // m为nil，这样执行会panic
}

func Test_strconv(t *testing.T) {
	var i interface{}
	
	defer log.Printf("hahahah")

	i = 1
	switch j := i.(type) { // 这里要赋给另一个值，为了减少变量，在不造成歧义的情况下，可这样写: i := i.(type)
	case int:
		log.Println(strconv.Itoa(j)) // 这里不能传入i
	}
}

