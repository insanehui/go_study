package main

import (
	"encoding/json"
	"log"

	// "gopkg.in/yaml.v2"
	"github.com/ghodss/yaml"
)

func main() {

	// xx := make(map[interface{}]interface{}) // 这种形式，两个yaml库都不能解析
	xx := make(map[string]interface{})

	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	yaml.Unmarshal([]byte(data), &xx) // 这里要用&xx，省掉"&"的话不能正常运行
	s1, err := json.Marshal(xx)       // 这里居然不能解析出来？
	if err != nil {
		log.Printf("err: %+v", err)
	}
	log.Printf("xx: %+v", xx)
	log.Printf("json: %+v", string(s1))

}
