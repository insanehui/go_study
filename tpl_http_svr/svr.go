package main

import (
	// "encoding/json"
	// "fmt"
	// "log"
	"log"
	"net/http"
	// "strconv"
	H "utils/http"
	// Mysql "utils/mysql"
	// _ "github.com/go-sql-driver/mysql"
	// V "github.com/asaskevich/govalidator"
)

func get_tpl_params(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId string `valid:"length(0|64)"`// 模板id
	}

	// 定义返回结构
	var ret struct {
		Err
		Data string `json:"data"`
	}

	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(error); ok { 
				log.Printf("haha: %+v", e.Error())
				ret.FromError(e)
			}
		}
		H.WriteJson(w, ret)
	}()

	H.Checkout_(r, &q)
	log.Printf("%+v", q)

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
