package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	J "utils/json"
	Y "utils/yaml"

	validator "github.com/asaskevich/govalidator"
	ms "github.com/mitchellh/mapstructure"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/ghodss/yaml"
	// "gopkg.in/yaml.v2"
)

type I map[interface{}]interface{}

var (
	filename = kingpin.Arg("filename", "the path of the yaml template file").Required().String()
	op       = kingpin.Arg("op", "operation: p / g").String()
	p        = kingpin.Flag("param", "path to param file").Short('p').String()

	m = make(I)
)

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
						r[key], _ = strconv.Atoi(input)
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
	case reflect.Slice:
		log.Printf("hahaha!")

	default:
	}

	return r
}

// 获取模板参数
func get_params() {
	log.Printf("params is: %+v", m["parameters"])
	// j, err := yaml.YAMLToJSON(m["parameters"])
	j := J.ToJson(m["parameters"])
	log.Printf("params: %+v", j)
}

func main() {

	kingpin.Parse()

	log.Printf("p is: %+v", *p)

	// 解析命令行参数
	if *p != "" { // 如果指定了参数文件
		m["parameters"] = Y.FromFile(*p)
	} else {
		data, err := ioutil.ReadFile(*filename)
		err = yaml.Unmarshal(data, &m)

		if err != nil {
			log.Printf("parse tpl err: %+v", err)
		}
	}

	if *op == "" {

		res := parse_parm(m["parameters"])
		// log.Printf("%+v", res)

		tpl, _ := template.ParseFiles(*filename)
		tpl.Execute(os.Stdout, res)

		// delete(m, "parameters")
		// log.Printf("===============")

		// res, _ := yaml.Marshal(m)
		// log.Printf("%s", res)

		// output the final yaml
	} else if *op == "p" {
		get_params()
	}

}
