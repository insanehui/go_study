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

	U "utils"
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
	Name string `json:"name" mapstructure:"name"`
	Type string `json:"type" mapstructure:"type"`
	Desc string `json:"description" mapstructure:"description"`
}

func read_a_param(para *Para, r I) {
	key := para.Name
	fmt.Printf("please input [ %s (%s) ]: ", key, para.Type)
	for true {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)


		switch para.Type {
		case "int":
			if validator.IsInt(input) {
				r[key], _ = strconv.Atoi(input)
				return
			}
		case "email":
			if validator.IsEmail(input) {
				r[key] = input
				return
			}
		case "string":
			r[key] = input
			return
		}
		fmt.Printf("invalid! try again: ")
	}
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
			para.Name = key
			ms.Decode(node, &para)
			read_a_param(&para, r)
		}
	case reflect.Slice:
		{
			var paras []Para
			U.Conv(t, &paras)
			for _, para := range paras {
				read_a_param(&para, r)
			}
		}

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
