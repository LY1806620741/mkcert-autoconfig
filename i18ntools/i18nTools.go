package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/tools/go/ast/astutil"
)

// learn for https://github.com/nicksnyder/go-i18n/blob/main/v2/goi18n/extract_command.go

func main() {
	// view()
	execute()
}

func ContainsInSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

var messageMap map[string]string

func init() {
	messageMap = initMessageMap()

}
func execute() error {
	paths := []string{"../"}
	dist := "../"
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if ContainsInSlice(paths, path) {
					return nil
				} else {
					// 跳過子目錄
					return filepath.SkipDir
				}
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			// Don't extract from test files.
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			//不解析要生成的i18n.go
			if info.Name() == "i18n.go" {
				return nil
			}

			log.Println("解析" + path)
			_, err = extractMessages(path, &messageMap, dist+info.Name())
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	//autoi18n.go
	file, _ := os.OpenFile(dist+"i18n.go", os.O_CREATE|os.O_RDWR, os.ModePerm)

	var structList []string = []string{}

	var initStr string = ""
	for k, v := range messageMap {

		if strings.HasPrefix(k, "scan") {
			structList = append(structList, k)

			initStr += `
		i18nText.` + k + ` = localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "` + k + `",
				Other: ` + v + `,
			},
		})
		`
		}
	}

	sort.Strings(structList)
	file.WriteString(
		`package main

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18nText struct {
	` + strings.Join(structList, ",") + ` string
}

var localizer *i18n.Localizer

//go:embed active.*.toml
var LocaleFS embed.FS

var i18nText I18nText

func init() {
	userLanguage, _ := jibber_jabber.DetectLanguage()
	userTag := language.MustParse(userLanguage)
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "active.zh.toml")
	tag, _, _ := language.NewMatcher([]language.Tag{
		language.English,
		language.Chinese,
	}).Match(userTag)
	localizer = i18n.NewLocalizer(bundle, tag.String())
	//初始化自动收集的待翻译文本
	` + initStr + `
}
	`)

	return nil
}

// 加载已存在的信息
func initMessageMap() map[string]string {

	mapo := make(map[string]interface{})
	file, _ := os.Open("../active.en.toml")
	defer file.Close()
	toml.NewDecoder(file).Decode(&mapo)

	mapres := make(map[string]string)

	for k, v := range mapo {

		vl, ok := v.([]interface{})
		if ok {
			for k1, v1 := range vl {
				fmt.Println("group", k, k1, v1, "ignore")
			}
		}

		vs, ok1 := v.(string)
		if ok1 {
			mapres[k] = vs
		}
	}

	return mapres
}

func extractMessages(sourcepath string, messageMap *map[string]string, name string) (map[string]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, sourcepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {

		n := c.Node()
		//只处理调用

		if ce, ok := n.(*ast.CallExpr); ok {
			//区分标识调用还是对象调用
			switch t := ce.Fun.(type) {
			case *ast.SelectorExpr:
				if !isMessageType(t) {
					return true
				}
				extractMessage(ce)
			case *ast.Ident:
				if isMessageType(t) {
					return true
				}
				//处理方法调用
				extractMessage(ce)
			default:
			}
		}
		return true
	})
	file, _ := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()
	if err := format.Node(file, token.NewFileSet(), node); err != nil {
		log.Fatalln("Error:", err)
	}
	return nil, nil
}

func isMessageType(expr ast.Expr) bool {

	//处理对象函数调用
	x, ok := expr.(*ast.SelectorExpr)

	if ok {
		x1, ok1 := x.X.(*ast.Ident)
		if ok1 {
			// log 的方法调用
			if x1.Name == "log" {
				return true
			}
		}
	}

	//处理函数调用 CallExpr中嵌套SelectorExpr，由上面处理
	x3, ok3 := expr.(*ast.Ident)

	if ok3 {
		//额外的日志打印函数

		if x3.Name == "fatalIfErr" {
			return true
		}
	}

	return false
}

// 解析其中的文字内容
func extractMessage(cl *ast.CallExpr) {
	_, isIdent := cl.Fun.(*ast.Ident)

	for i, elt := range cl.Args {

		//处理log下的方法
		basiclin, basicok := elt.(*ast.BasicLit)

		if basicok {
			if isIdent {
				key := PutIfExistMessage(messageMap, basiclin.Value)
				newExpr, _ := parser.ParseExpr(`
					i18nText.` + key + `
				`)
				cl.Args[i] = newExpr
				log.Println("函数调用() " + basiclin.Value)

			} else {

				if basiclin.Kind == token.STRING {
					key := PutIfExistMessage(messageMap, basiclin.Value)
					// 修改参数
					newExpr, _ := parser.ParseExpr(`
										i18nText.` + key + `
									`)
					cl.Args[i] = newExpr
				} else {
					log.Println("忽略对象.函数调用()" + basiclin.Value)
				}
			}

		}

	}
}

func PutIfExistMessage(maps map[string]string, value string) string {
	for k, v := range maps {
		if v == value {
			log.Println("重复")
			return k
		}
	}
	key := fmt.Sprintf("scan%d", len(maps))
	maps[key] = value
	return key
}

func extractStringLiteral(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			return "", false
		}
		return s, true
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}
		x, ok := extractStringLiteral(v.X)
		if !ok {
			return "", false
		}
		y, ok := extractStringLiteral(v.Y)
		if !ok {
			return "", false
		}
		return x + y, true
	case *ast.Ident:
		if v.Obj == nil {
			return "", false
		}
		switch z := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if len(z.Values) == 0 {
				return "", false
			}
			s, ok := extractStringLiteral(z.Values[0])
			if !ok {
				return "", false
			}
			return s, true
		}
	case *ast.CallExpr:
		if fun, ok := v.Fun.(*ast.Ident); ok && fun.Name == "string" {
			return extractStringLiteral(v.Args[0])
		}
	}
	return "", false
}
