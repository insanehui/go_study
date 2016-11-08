package main

import (
	"log"
	"reflect"

	validator "github.com/asaskevich/govalidator"
)

func f(t interface{}) {

	v := reflect.ValueOf(t)
	switch v.Kind() {

	case reflect.Map:
		log.Printf("Map")
		for _, rkey := range v.MapKeys() {
			log.Printf("rkey's kind is: %+v", rkey.Kind())

			iv := rkey.Interface()
			vv := reflect.ValueOf(iv)
			log.Printf("%+v", vv.String())
		}

	default:

	}

}

func main() {
	log.Printf("haha")

	m := make(map[interface{}]interface{})

	m["haha"] = 1
	m[123] = 2

	// log.Printf("%+v", m)

	log.Println(validator.IsInt("123"))

}
