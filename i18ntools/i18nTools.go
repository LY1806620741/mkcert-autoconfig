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

func execute() error {
	paths := []string{"../"}
	dist := "test/"
	messageMap := initMessageMap()
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
			if info.Name() == "autoi18n.go" {
				return nil
			}

			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			log.Println("解析" + path)
			_, err = extractMessages(buf, &messageMap, dist+info.Name())
			if err != nil {
				return err
			}
			// messages = append(messages, msgs...)
			return nil
		}); err != nil {
			return err
		}
	}

	//autoi18n.go
	file, _ := os.OpenFile(dist+"autoi18n.go", os.O_CREATE|os.O_RDWR, os.ModePerm)

	var structList []string = []string{}

	var initStr string = ""
	for k, v := range messageMap {
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
	` + initStr + `
	//初始化自动收集的待翻译文本
}
	`)

	// messageTemplates := map[string]*i18n.MessageTemplate{}
	// for _, m := range messages {
	// 	if mt := i18n.NewMessageTemplate(m); mt != nil {
	// 		messageTemplates[m.ID] = mt
	// 	}
	// }
	// fmt.Println(messageTemplates)
	return nil
}

type extractor struct {
	messageMap map[string]string
}

func (e *extractor) Visit(node ast.Node) ast.Visitor {
	e.extractMessages(node)
	return e
}

// 加载已存在的信息
func initMessageMap() map[string]string {
	// bundle := i18n.NewBundle(language.English)
	// bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// bundle.LoadMessageFile("../active.en.toml")

	// messageTemplates := reflect.ValueOf(bundle).Elem().FieldByName("messageTemplates").Interface().((map[language.Tag]map[string]*MessageTemplate))

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
			// fmt.Println(string(k), v)
		}
	}

	// fmt.Println(mapo)

	return mapres
}

func extractMessages(buf []byte, messageMap *map[string]string, name string) (map[string]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	extractor := newExtractor(messageMap)
	ast.Walk(extractor, node)
	file, _ := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()
	format.Node(file, token.NewFileSet(), node)
	return extractor.messageMap, nil
}

func newExtractor(messageMap *map[string]string) *extractor {
	return &extractor{messageMap: *messageMap}
}

func (e *extractor) extractMessages(node ast.Node) {
	//只处理调用
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	//区分标识调用还是对象调用
	switch t := ce.Fun.(type) {
	case *ast.SelectorExpr:
		if !e.isMessageType(t) {
			return
		}
		e.extractMessage(ce)
	case *ast.Ident:
		if !e.isMessageType(t) {
			return
		}
		//处理方法调用
		e.extractMessage(ce)
	default:
		// log.Println(t)
	}
}

func (e *extractor) isMessageType(expr ast.Expr) bool {

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

func unwrapIdent(e ast.Expr) *ast.Ident {
	switch et := e.(type) {
	case *ast.Ident:
		return et
	default:
		return nil
	}
}

func unwrapSelectorExpr(e ast.Expr) *ast.SelectorExpr {
	switch et := e.(type) {
	case *ast.SelectorExpr:
		return et
	case *ast.StarExpr:
		se, _ := et.X.(*ast.SelectorExpr)
		return se
	default:
		return nil
	}
}

// 解析其中的文字内容
func (e *extractor) extractMessage(cl *ast.CallExpr) {
	// data := make(map[string]string)
	_, isIdent := cl.Fun.(*ast.Ident)

	for i, elt := range cl.Args {

		//处理log下的方法
		basiclin, basicok := elt.(*ast.BasicLit)

		if basicok {
			if isIdent {
				key := PutIfExistMessage(e.messageMap, basiclin.Value)
				newExpr, _ := parser.ParseExpr(`
				i18nText.` + key + `
			`)
				cl.Args[i] = newExpr
				log.Println("函数调用() " + basiclin.Value)

			} else {

				if basiclin.Kind == token.STRING {
					key := PutIfExistMessage(e.messageMap, basiclin.Value)
					//修改参数
					newExpr, _ := parser.ParseExpr(`
										i18nText.` + key + `
									`)
					cl.Args[i] = newExpr
					// log.Println("对象.函数调用()" + basiclin.Value)
				} else {
					log.Println("忽略对象.函数调用()" + basiclin.Value)
				}
			}

		}
		continue

		// kve, ok := elt.(*ast.KeyValueExpr)
		// if !ok {
		// 	continue
		// }
		// key, ok := kve.Key.(*ast.Ident)
		// if !ok {
		// 	continue
		// }
		// v, ok := extractStringLiteral(kve.Value)
		// if !ok {
		// 	continue
		// }
		// data[key.Name] = v
	}
	// if len(data) == 0 {
	// 	return
	// }
	// if messageID := data["MessageID"]; messageID != "" {
	// 	data["ID"] = messageID
	// }
	// e.messages = append(e.messages, i18n.MustNewMessage(data))
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
