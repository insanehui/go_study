package main

import (
	"encoding/json"
	"log"

	"testing"

	"gopkg.in/yaml.v2"
)

func TestY2J(t *testing.T) {

	m := make(map[interface{}]interface{})

	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	yaml.Unmarshal([]byte(data), &m)
	s1, _ := json.Marshal(m) // 这里居然不能解析出来？
	log.Printf("%+v", m)
	log.Printf("%+v", string(s1))

}
