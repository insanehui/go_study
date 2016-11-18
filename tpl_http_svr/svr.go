package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	U "utils"
	H "utils/http"
	J "utils/json"
	I "utils/io"
	Mysql "utils/mysql"
	// "github.com/ghodss/yaml"
	Y "utils/yaml"

	V "github.com/asaskevich/govalidator"
)

var db *Mysql.DB

// 一个模板参数的规则结构
type Para struct {
	Label      string                 `json:"label"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Desc       string                 `json:"description"`
	Optional   bool                   `json:"optional,omitempty"`
	EnumValues map[string]interface{} `json:"enum_values,omitempty"`
	Default    interface{}            `json:"default,omitempty"`
}

type Paras []Para

// 检查模板id是否存在
func check_tpl_id_(id string) {

	// 从数据库查询 tpl_id 是否存在
	c := db.QryInt("select count(1) from template where id = ?", id)
	if c <= 0 {
		log.Panicf("template id %s does not exist!", id)
	}
}

func get_params(id string) Paras {
	var paras Paras
	para_fname := path.Join("tpls", id, "para.yaml")
	Y.FromFileTo_(para_fname, &paras)
	return paras
}

// 校验values
func check_vals_(vals map[string]interface{}, paras Paras) {
	for _, para := range paras {

		name := para.Name
		v := vals[name]
		if v == nil {
			if para.Optional {
				continue // 可选值
			} else {
				log.Panicf("miss para: %s", name)
			}

		} else {
			sv := U.ToStr(v)
			log.Printf("%s : %s", name, sv)

			switch para.Type {
			case "int":
				if !V.IsInt(sv) {
					log.Panicf("para '%s' is not int", name)
				}
			case "email":
				if !V.IsEmail(sv) {
					log.Panicf("para '%s' is not email", name)
				}
			case "string":
			}
		}
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	panic("fuck")
}

func get_tpl_params(w http.ResponseWriter, r *http.Request) {

	// 传入参数
	var q struct {
		TplId string `valid:"length(1|64)" json:"id"` // 模板id
	}

	// 返回
	var ret struct {
		Err
		Data interface{} `json:"data"`
	}

	H.JsonDo(w, r, &q, &ret, func() {

		// check_tpl_id_(q.TplId)

		// 读其para的配置文件
		ret.Data = get_params(q.TplId)

	})

}

func gen_blueprint(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId  string `valid:"length(1|64)"` // 模板id
		Values string `valid:"json"`
		Op     string `valid:"-"` // 可选参数, 暂时支持 get_yaml
	}

	// 定义返回结构
	var ret struct {
		Err
		Data string `json:"data"`
	}

	H.Do(w, r, &q, &ret,
		func() {

			// check_tpl_id_(q.TplId)

			// 取到values
			var vals map[string]interface{}
			J.StrTo(q.Values, &vals)

			// 取paras
			paras := get_params(q.TplId)
			// 验证
			check_vals_(vals, paras)

			// 渲染模板 !! json形式的对象，是不支持模板里的if的
			tpl_fname := path.Join("tpls", q.TplId, "tpl.yaml") // 找到模板的路径
			tpl, _ := template.ParseFiles(tpl_fname)            // 实例化模板对象

			bp := new(bytes.Buffer)
			tpl.Execute(bp, vals)
			ret.Data = bp.String()
			log.Printf("%+v", bp.String())

		},
		func(p interface{}) { // in defer
			if p != nil {
				ret.FromPanic(p)
			}

			if q.Op == "get_yaml" {
				if ret.Ok() {
					io.WriteString(w, ret.Data)
				} else {
					http.Error(w, ret.Msg, http.StatusNotFound)
				}
			} else {
				H.WriteJson(w, ret)
			}
		})

}

// 新增一个模板
func new_tpl(w http.ResponseWriter, r *http.Request) {

	// 传入参数
	var q struct {
		TplId string `valid:"-"` // 模板id
		Paras string `valid:"-"`
		Tpl string `valid:"length(1|99999)"`
	}

	// 返回
	var ret struct {
		Err
		Data struct { 
			TplId string `json:"tpl_id"`
		} `json:"data"`
	}

	H.JsonDo(w, r, &q, &ret, func() {
		id := q.TplId

		TplInit_(id)

		if data := q.Tpl; data != "" {
			p := TplPath(id)
			I.WriteFile_(p, data)
		}

		if data := q.Paras; data != "" {
			p := TplParaPath(id)
			I.WriteFile_(p, data)
		}

		ret.Data.TplId = q.TplId
	})

}

func init() {

	db = Mysql.Open_("mysql", "blueprint:ctg123@tcp(10.10.12.2:3306)/blueprint?charset=utf8")

	http.HandleFunc("/get_tpl_params", get_tpl_params)
	http.HandleFunc("/gen_blueprint", gen_blueprint)
	http.HandleFunc("/new_tpl", new_tpl)
	http.HandleFunc("/test", test)
}

func main() {
	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
