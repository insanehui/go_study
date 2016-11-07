package main

import (
	"log"

	"testing"

	// "gopkg.in/yaml.v2"
	j "utils/json"

	"github.com/ghodss/yaml"
)

type Para struct {
	Name string `yaml:"name" mapstructure:"name"`
	Type string `yaml:"type" mapstructure:"type"`
	Desc string `yaml:"description" mapstructure:"description"`
}

func TestY2J(t *testing.T) {

	{
		m := make(map[string]interface{})

		var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
		yaml.Unmarshal([]byte(data), &m)
		// s1, _ := json.Marshal(m) // 这里居然不能解析出来？
		s1 := j.ToJson(m)
		log.Printf("%+v", m)
		log.Printf("%+v", string(s1))
	}

	{
		var m []Para
		var data = `
- name: name
  type: string
  description: this's your name
- name: age
  type: int
  description: how old are you
- name: company
  type: string
  description: where you work
- name: email
  type: email
  description: Arbitrary key/value metadata
		`
		err := yaml.Unmarshal([]byte(data), &m)
		if err != nil {
			log.Printf("err: %+v", err)
		}

	}

}
