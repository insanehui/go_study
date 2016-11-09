package main

import (
	"html/template"
	"log"
	"os"
	J "utils/json"
)

func main() {

	{
		data := map[string]interface{}{
			"Name":    "guanghui",
			"Age":     30,
			"Company": "CloudToGo",
		}

		tmpl, _ := template.New("test").Parse(`

		my name is {{.Name}}
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
		log.Printf("%+v", data)

		tmpl, _ := template.New("xx").Parse(`

		my name is {{.Name}}
		i'm {{.Age}} years old
		{{if gt .Age 10}}
		i'm an adult
		{{else}}
		i'm a child
		{{end}}


		`)

		tmpl.Execute(os.Stdout, data)
	}

}
