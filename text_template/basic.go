package main

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
	J "utils/json"
)

func toint(i interface{}) int {
	if v, ok := i.(int); ok {
		return v
	} else if v, ok := i.(float64); ok {
		return int(v)
	}
	return 0
}

func kind(i interface{}) string {
	v := reflect.ValueOf(i)
	return fmt.Sprintf("%+v", v.Type())
}

func main() {

	funcMap := template.FuncMap{
		"toint": toint,
		"title": strings.Title,
		"kind":  kind,
	}

	{
		data := map[string]interface{}{
			"Name":    "guanghui li",
			"Age":     30,
			"Company": "CloudToGo",
		}

		tmpl, _ := template.New("test").Funcs(funcMap).Parse(`

		my name is {{title .Name}}
		my name is {{title "aa bb"}}
		i'm {{.Age}} years old
		{{if gt .Age 10}}
		i'm an adult
		{{else}}
		i'm a child
		{{end}}

		`)

		tmpl.Execute(os.Stdout, data)
	}

	{
		data := J.Str2Var(`{
			"Name":    "guanghui",
			"Age":     30,
			"Company": "CloudToGo"
		}`)

		tmpl, _ := template.New("xx").Funcs(funcMap).Parse(`

		my name is {{.Name}}
		i'm {{.Age}} years old
		{{$a := "123"}}{{/* 定义一个变量 */}}
		{{ toint .Age }}
		{{kind .Age}}
		{{title "aa bb cc"}}
		{{if gt .Age 10.0}}
		i'm an adult
		{{else}}
		i'm a child
		{{end}}


		`)

		tmpl.Execute(os.Stdout, data)
	}

}
