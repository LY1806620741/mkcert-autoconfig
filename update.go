package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

//像 https://github.com/sanbornm/go-selfupdate 目前还在测试
// https://blog.csdn.net/wangge20091126/article/details/128036546

func selfUpdate() {
	fileName, _ := os.Executable()
	oldFileName := fileName + ".old"
	os.Remove(oldFileName)
	err := os.Rename(fileName, oldFileName)
	if err != nil {
		fmt.Print(err)
		return
	}
	time.Sleep(5 * time.Second)
	os.Remove(oldFileName)

	// 到程序安装路径下去执行启动命令(预防相对路径方式启动)
	daemon := "timeout /T 3 & " + fileName + " 2>&1 &"
	_ = exec.Command("cmd.exe", "/C", daemon).Start()
}
