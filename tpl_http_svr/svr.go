package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	H "utils/http"
	J "utils/json"
	Mysql "utils/mysql"
	// "github.com/ghodss/yaml"
	Y "utils/yaml"
	// V "github.com/asaskevich/govalidator"
)

var db *Mysql.DB

// 一个模板参数的规则结构
type Para struct {
	Label      string                 `json:"label"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Desc       string                 `json:"description"`
	Optional   bool                   `json:"optional"`
	EnumValues map[string]interface{} `json:"enum_values"`
	Default    interface{}            `json:"default"`
}

type Paras []Para

// 检查模板id是否存在
func check_tpl_id_(id string) {

	// 从数据库查询 tpl_id 是否存在
	c := db.QryInt("select count(1) from template where id = ?", id)
	if c <= 0 {
		panic(fmt.Errorf("template id %s does not exist!", id))
	}

}

func get_params(id string) Paras {
	var paras Paras
	para_fname := path.Join("tpls", id, "para.yaml")
	Y.FromFileTo_(para_fname, &paras)
	return paras
}

// 校验values
func check_vals(vals interface{}, rules interface{}) {
}

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
	check_tpl_id_(q.TplId)

	// 读其para的配置文件
	ret.Data = get_params(q.TplId)

	log.Printf("%+v", q)

}

func gen_blueprint(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId  string `valid:"length(0|64)" json:"id"` // 模板id
		Values string `valid:"json"`
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
	check_tpl_id_(q.TplId)

	// 取到values
	vals := J.Str2Var(q.Values)

	// TODO
	// 验证values有效性

	// 渲染模板 !! json形式的对象，是不支持模板里的if的
	tpl_fname := path.Join("tpls", q.TplId, "tpl.yaml") // 找到模板的路径
	tpl, _ := template.ParseFiles(tpl_fname)            // 实例化模板对象

	bp := new(bytes.Buffer)
	tpl.Execute(bp, vals)
	ret.Data = bp.String()
	log.Printf("%+v", bp.String())
}

func init() {

	db = Mysql.Open_("mysql", "blueprint:ctg123@tcp(10.10.12.2:3306)/blueprint?charset=utf8")

	http.HandleFunc("/get_tpl_params", get_tpl_params)
	http.HandleFunc("/gen_blueprint", gen_blueprint)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", nil))
}
