package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 初始化说明文本
var execNameWithOutSuffix string

func init() {
	exePath, _ := os.Executable()
	_, execName := filepath.Split(exePath)
	exeSuffix := path.Ext(execName)
	execNameWithOutSuffix = strings.TrimSuffix(execName, exeSuffix)
}
