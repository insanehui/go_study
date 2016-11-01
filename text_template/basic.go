package main

import (
	"html/template"
	"os"
)

func main() {

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
