package main

import (
	// "encoding/json"
	// "fmt"
	// "log"
	// "net/http"
	// "strconv"
	// H "utils/http"
	// J "utils/json"
	// Mysql "utils/mysql"

	// _ "github.com/go-sql-driver/mysql"
)

func get_tpl_params(w http.ResponseWriter, r *http.Request) {
}

func gen_blueprint(w http.ResponseWriter, r *http.Request) {
}

func init() {

	// db, _ = Mysql.Open("mysql", "teamtalk:12345@tcp(115.29.233.2:3306)/PicPath?charset=utf8")

	http.HandleFunc("/get_tpl_params", get_tpl_params)
	http.HandleFunc("/gen_blueprint", gen_blueprint)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", nil))
}
