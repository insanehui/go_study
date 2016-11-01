package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	validator "github.com/asaskevich/govalidator"
	ms "github.com/mitchellh/mapstructure"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	filename = kingpin.Arg("filename", "the path of the yaml template file").Required().String()
)

type I map[interface{}]interface{}

type Para struct {
	Type string `yaml:"type" mapstructure:"type"`
	Desc string `yaml:"description" mapstructure:"description"`
}

func parse_parm(t interface{}) I {
	// log.Printf("t is:\n%+v", t)

	r := make(I)

	// now t is a map, get it as a map
	// use reflect
	v := reflect.ValueOf(t)
	switch v.Kind() {
	case reflect.Map:
		for _, rkey := range v.MapKeys() {
			key := reflect.ValueOf(rkey.Interface()).String()
			node := v.MapIndex(rkey).Interface()
			var para Para
			ms.Decode(node, &para)
			// log.Printf("%+v", para)
			fmt.Printf("please input [ %s (%s) ]: ", key, para.Type)

			// scan a line
		try_loop:
			for true {
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.Replace(input, "\n", "", -1)
				input = strings.Replace(input, "\r", "", -1)

				switch para.Type {
				case "int":
					if validator.IsInt(input) {
						r[key] = input
						break try_loop
					}
				case "email":
					if validator.IsEmail(input) {
						r[key] = input
						break try_loop
					}
				case "string":
					r[key] = input
					break try_loop
				}
				fmt.Printf("invalid! try again: ")
			}
		}
	default:
	}

	return r
}

func main() {

	kingpin.Parse()
	m := make(I)

	// read the parameters section
	data, err := ioutil.ReadFile(*filename)

	// data1 =

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		// log.Printf("%+v", err)
	}
	// log.Printf("%+v", m)

	// ask user to input params
	res := parse_parm(m["parameters"])
	// log.Printf("%+v", res)

	tpl, _ := template.ParseFiles(*filename)
	tpl.Execute(os.Stdout, res)

	// delete(m, "parameters")
	// log.Printf("===============")

	// res, _ := yaml.Marshal(m)
	// log.Printf("%s", res)

	// output the final yaml
}
