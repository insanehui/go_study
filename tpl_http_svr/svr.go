package main

import (
	// "encoding/json"
	"path"
	"fmt"
	// "log"
	"log"
	"net/http"
	// "strconv"
	H "utils/http"
	// J "utils/json"
	Mysql "utils/mysql"
	// "github.com/ghodss/yaml"
	Y "utils/yaml"

	// _ "github.com/go-sql-driver/mysql"
	// V "github.com/asaskevich/govalidator"
)

var db *Mysql.DB

func get_tpl_params(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId string `valid:"length(0|64)" json:"id"` // 模板id
	}

	// 定义返回结构
	var ret struct {
		Err
		Data interface{} `json:"data"`
	}

	defer func() {
		if p := recover(); p != nil {
			// 先支持 error 格式
			if e, ok := p.(error); ok {
				log.Printf("haha: %+v", e.Error())
				ret.FromError(e)
			}
		}
		H.WriteJson(w, ret)
	}()

	H.Checkout_(r, &q)

	// 从数据库查询 tpl_id 是否存在
	res := db.Exist("template", q)
	if res == nil {
		panic(fmt.Errorf("template id %s does not exist!", q.TplId))
	}

	// 读其para的配置文件
	para_fname := path.Join("tpls", q.TplId, "para.yaml")
	ret.Data = Y.FromFile(para_fname)

	log.Printf("%+v", q)

}

func gen_blueprint(w http.ResponseWriter, r *http.Request) {
}

func init() {

	db = Mysql.Open_("mysql", "blueprint:ctg123@tcp(10.10.12.2:3306)/blueprint?charset=utf8")

	http.HandleFunc("/get_tpl_params", get_tpl_params)
	http.HandleFunc("/gen_blueprint", gen_blueprint)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", nil))
}
