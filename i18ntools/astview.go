package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// demo_test.go
var srcCode = `
package main

import (
	"errors"
	"log"
)

func logtest() {
	log.Printf("测试")
	log.Println("测试")
	err := errors.New("test")
	fatalIfErr(err, "delete cert")
}

`

func view() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", srcCode, 0)
	if err != nil {
		panic(err.Error())
	}
	ast.Print(fset, f)
	file, _ := os.OpenFile("demo_test.ast", os.O_CREATE|os.O_RDWR, os.ModePerm)
	ast.Fprint(file, fset, f, ast.NotNilFilter)
}
