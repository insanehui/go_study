package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	J "utils/json"
	// Y "utils/yaml"
	"path"
	"github.com/satori/go.uuid"

	"net/http"
	"net/url"
	I "utils/io"

	V "github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
)

func init() {
	V.SetFieldsRequiredByDefault(true)
}

func Test_valid1(t *testing.T) {
	{
		log.Println("=========================")
		var q = struct {
			A string `valid:"ip"`
			B string `valid:"int"`
		}{"1.1.1.1", ""}

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}
}

func TestStructAssign(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}

	type B struct {
		A
		Company string
	}
	a := A{"aa", 13}
	var b B
	b.A = a // 可以将A作为B的一部分赋值，用组合的方式实现类似于继承的功能

	log.Printf("%+v", b)
	log.Printf("%+v", J.ToJson(b))
}

func TestValidator(t *testing.T) {
	{
		log.Println("=========================")
		var q = struct {
			TplId string `valid:"length(1|5)"` // 模板id
		}{""}

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}

	{
		log.Println("=========================")
		var q = struct {
			TplId string `valid:"length(0|3)"` // 模板id
		}{"haha"}

		// err表示具体错误的原因 result直接为bool类型
		result, err := V.ValidateStruct(&q)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		log.Println(result)
	}
}

func Test_nil(t *testing.T) {
	var s []string // slice的初始状态为nil
	log.Println(s == nil)
}

func Test_path(t *testing.T) {
	log.Println(path.Join("a", "b", "c"))
}

func Test_mixin(t *testing.T) {
	type A struct {
		I int
	}

	type B struct {
		J int
	}

	type C struct {
		Aa A
		B
	}

	var c C
	c.Aa.I = 1 // 这里不能用 c.A.i
	c.B.J = 2  // 这里只能用c.B.i
	log.Printf("%+v", c)
	log.Printf("%+v", J.ToJson(c))
}

func Test_yaml2tpl(t *testing.T) {
	str := `
Name:    1
Age:     30
Company: 2
`
	// var m map[string]interface{} // map的值为 interface{}的话，是不能对if生效的
	var m map[string]int // 定成具体类型之后，才可以

	yaml.Unmarshal([]byte(str), &m)
	log.Printf("%+v", m)

	tmpl, _ := template.New("xx").Parse(`

		my name is {{.Name}}
		i'm {{.Age}} years old
		{{if gt .Age 10}}
		i'm an adult
		{{else}}
		i'm a child
		{{end}}


		`)

	tmpl.Execute(os.Stdout, m)
}

func Test_booldef(t *testing.T) {
	var a bool
	log.Println(a == false)
	log.Println(a == true)
}

func Test_interface(t *testing.T) {
	var a interface{}
	a = false
	log.Println(a == nil)
	a = 1
	log.Println(a == true)
	log.Println(a.(int))
}

func Test_make(t *testing.T) {
	// var m map[string]int
	// m["aa"] = 1 // m为nil，这样执行会panic
}

func Test_strconv(t *testing.T) {
	var i interface{}

	defer log.Printf("hahahah")

	i = 1
	switch j := i.(type) { // 这里要赋给另一个值，为了减少变量，在不造成歧义的情况下，可这样写: i := i.(type)
	case int:
		log.Println(strconv.Itoa(j)) // 这里不能传入i
	}
}

func Test_print(t *testing.T) {
	i := 123
	fmt.Println(i)
}

type C struct {
}

func (me *C) F() {
	log.Printf("f")
}

func (me *C) G() {
	log.Printf("g")
}

func Test_intf_assert(t *testing.T) {
	type A interface {
		F()
	}
	type B interface {
		G()
	}

	var c C
	var i interface{}
	i = &c

	switch v := i.(type) {
	// 以下几个条件都符合
	// 所以谁在前面谁就被执行
	case B:
		log.Printf("bb")
		v.G()
	case *C:
		log.Printf("cc")
	case A:
		log.Printf("aa")
		v.F()
	default:
		log.Printf("......")
	}

	if _, ok := i.(B); ok {
		log.Printf("shit b")
	}

	if _, ok := i.(A); ok {
		log.Printf("shit a")
	}

	if _, ok := i.(*C); ok {
		log.Printf("shit c")
	}
}

func Test_postform(t *testing.T) {

	tpl := `
aaa
 bbb
  cccc
`

	resp, err := http.PostForm("http://localhost:8080/new_tpl",
		url.Values{"tpl_id": {"12345"}, "tpl": {tpl}})

	if err != nil {
		log.Printf("post err: %+v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response err: %+v", err)
	}

	fmt.Println(string(body))
}

func Test_uuid(t *testing.T) {
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)
}

func Test_reflect(t *testing.T) {
	type My int
	var i My
	i = 1
	rv := reflect.ValueOf(i)
	tp := reflect.TypeOf(i)
	log.Println(tp == rv.Type())            // true
	log.Println(rv.Type())                  // main.My
	log.Println(rv.Kind())                  // int
	log.Println(reflect.ValueOf(&i).Kind()) //
}

func Test_log(t *testing.T) {
	// 定义一个文件
	fileName := "ll.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	debugLog := log.New(logFile, "[Debug]", log.LstdFlags)
	debugLog.Println("A debug message here")
	//配置一个日志格式的前缀
	debugLog.SetPrefix("[Info]")
	debugLog.Println("A Info Message here ")
	//配置log的Flag参数
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	debugLog.Println("A different prefix")
}

func Test_delims(t *testing.T) {

	{
		data := map[string]interface{}{
			"Name":    "guanghui li",
			"Age":     30,
			"Company": "CloudToGo",
		}

		b := I.ReadFile_("tpl")
		tmpl, _ := template.New("test").Delims("<%", "%>").Parse(string(b))

		// tmpl, _ := template.New("test").Delims("<%", "%>").Parse(`
// my name is <% .Name %>
		// `)

		// 用以下方法，是不生效的！
		// tmpl, _ := template.New("test").Delims("<%", "%>").ParseFiles("tpl")

		tmpl.Execute(os.Stdout, data)
	}
}
