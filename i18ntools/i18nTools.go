package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// learn for https://github.com/nicksnyder/go-i18n/blob/main/v2/goi18n/extract_command.go

func main() {
	// view()
	execute()
}

func execute() error {
	paths := []string{"../"}
	messages := []*i18n.Message{}
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
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

			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			log.Println("解析" + path)
			msgs, err := extractMessages(buf)
			if err != nil {
				return err
			}
			messages = append(messages, msgs...)
			return nil
		}); err != nil {
			return err
		}
	}
	messageTemplates := map[string]*i18n.MessageTemplate{}
	for _, m := range messages {
		if mt := i18n.NewMessageTemplate(m); mt != nil {
			messageTemplates[m.ID] = mt
		}
	}
	fmt.Println(messageTemplates)
	return nil
}

type extractor struct {
	messages []*i18n.Message
}

func (e *extractor) Visit(node ast.Node) ast.Visitor {
	e.extractMessages(node)
	return e
}

func extractMessages(buf []byte) ([]*i18n.Message, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	extractor := newExtractor(file)
	ast.Walk(extractor, file)
	return extractor.messages, nil
}

func newExtractor(file *ast.File) *extractor {
	return &extractor{}
}

func (e *extractor) extractMessages(node ast.Node) {
	switch t := node.(type) {
	case *ast.SelectorExpr:
		if !e.isMessageType(t) {
			return
		}
		// e.extractMessage(t)
	case *ast.CallExpr:
		if !e.isMessageType(t) {
			return
		}
		//处理方法调用
		fmt.Println(t)
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
	x3, ok3 := expr.(*ast.CallExpr)

	if ok3 {
		//额外的日志打印函数
		x4, ok4 := x3.Fun.(*ast.Ident)

		if ok4 {
			if x4.Name == "fatalIfErr" {
				return true
			}
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

func (e *extractor) extractMessage(cl *ast.CompositeLit) {
	data := make(map[string]string)
	for _, elt := range cl.Elts {
		kve, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kve.Key.(*ast.Ident)
		if !ok {
			continue
		}
		v, ok := extractStringLiteral(kve.Value)
		if !ok {
			continue
		}
		data[key.Name] = v
	}
	if len(data) == 0 {
		return
	}
	if messageID := data["MessageID"]; messageID != "" {
		data["ID"] = messageID
	}
	e.messages = append(e.messages, i18n.MustNewMessage(data))
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
