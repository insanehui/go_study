package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func Test_new_tpl(t *testing.T) {

	tpl := `aaa
 bbb
  cccc
`
	id := "12345"

	resp, err := http.PostForm("http://localhost:8080/new_tpl",
		url.Values{"tpl_id": {id}, "tpl": {tpl}})

	if err != nil {
		log.Printf("post err: %+v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response err: %+v", err)
	}

	log.Println(string(body))
}

func Test_get_tpl_params(t *testing.T) {

	id := "234242"

	resp, err := http.PostForm("http://localhost:8080/get_tpl_params",
		url.Values{"tpl_id": {id}})

	if err != nil {
		log.Printf("post err: %+v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response err: %+v", err)
	}

	log.Println(string(body))
}
