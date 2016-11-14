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

func get_tpl_params_old(w http.ResponseWriter, r *http.Request) {

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
			} else if e, ok := p.(string); ok {
				ret.Msg = e
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

type IHttpRet interface {
	FromPanic(interface{})
}

// para为自定义的传入参数
// ret为自定义的返回参数
func http_json(w http.ResponseWriter, r *http.Request, para interface{}, ret IHttpRet, fn func()) {

	defer func() {
		if p := recover(); p != nil {
			ret.FromPanic(p)
		}
		H.WriteJson(w, ret)
	}()

	H.Checkout_(r, para)
	fn()
}

func get_tpl_params(w http.ResponseWriter, r *http.Request) {

	// 传入参数
	var q struct {
		TplId string `valid:"length(0|64)" json:"id"` // 模板id
	}

	// 返回
	var ret struct {
		Err
		Data interface{} `json:"data"`
	}

	http_json(w, r, &q, &ret, func() {

		// check_tpl_id_(q.TplId)

		// 读其para的配置文件
		ret.Data = get_params(q.TplId)

	})

}

// 更通用的（不一定返回json）
// end(p) 将在 defer里执行, 参数p 为panic
func http_do(w http.ResponseWriter, r *http.Request, para interface{}, ret interface{}, fn func(), end func(p interface{})) {

	defer func() {
		p := recover()
		end(p)
	}()

	H.Checkout_(r, para)
	fn()
}

func gen_blueprint_old(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId  string `valid:"length(0|64)"` // 模板id
		Values string `valid:"json"`
		Op     string `valid:"-"` // 可选参数, 暂时支持 get_yaml
	}

	// 定义返回结构
	var ret struct {
		Err
		Data string `json:"data"`
	}

	defer func() {
		if p := recover(); p != nil {
			// 先支持 error 格式
			if e, ok := p.(error); ok {
				log.Printf("haha: %+v", e.Error())
				ret.FromError(e)
			} else if e, ok := p.(string); ok {
				ret.Msg = e
			}
		}

		if q.Op == "get_yaml" {
			io.WriteString(w, ret.Data)
		} else {
			H.WriteJson(w, ret)
		}
	}()

	H.Checkout_(r, &q)
	check_tpl_id_(q.TplId)

	// 取到values
	// vals := J.Str2Var(q.Values)
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
}

func gen_blueprint(w http.ResponseWriter, r *http.Request) {

	var q struct {
		TplId  string `valid:"length(0|64)"` // 模板id
		Values string `valid:"json"`
		Op     string `valid:"-"` // 可选参数, 暂时支持 get_yaml
	}

	// 定义返回结构
	var ret struct {
		Err
		Data string `json:"data"`
	}

	http_do(w, r, &q, &ret, func() {

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
		func(p interface{}) {
			if p != nil {
				// 先支持 error 格式
				if e, ok := p.(error); ok {
					log.Printf("haha: %+v", e.Error())
					ret.FromError(e)
				} else if e, ok := p.(string); ok {
					ret.Msg = e
				}
			}

			if q.Op == "get_yaml" {
				io.WriteString(w, ret.Data)
			} else {
				H.WriteJson(w, ret)
			}
		})

}

func init() {

	db = Mysql.Open_("mysql", "blueprint:ctg123@tcp(10.10.12.2:3306)/blueprint?charset=utf8")

	http.HandleFunc("/get_tpl_params", get_tpl_params)
	http.HandleFunc("/gen_blueprint", gen_blueprint)
	// http.HandleFunc("/test", test)
}

func main() {
	log.Println("running...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
