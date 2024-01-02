package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var magic = []byte("selffs")

var cert []byte

var caInit = false

func init() {
	exePath, _ := os.Executable()
	file, err := os.Open(exePath)
	if err != nil {
		fmt.Println(i18nText.
			scan99,

			err)
		return
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(i18nText.
			scan100,

			err)
		return
	}
	fileSize := fileInfo.Size()

	// 确定读取的位置
	numBytes := len(magic)
	offset := int64(numBytes)
	if fileSize < int64(numBytes) {
		offset = fileSize
	}
	seekPos := fileSize - offset

	// 移动到读取的位置
	_, err = file.Seek(seekPos, io.SeekStart)
	if err != nil {
		fmt.Println(i18nText.
			scan101,

			err)
		return
	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(i18nText.
			scan102,

			err)
		return
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadCA(*file)
	} else {
		// division()
	}

}

// 加载额外文件系统
func loadCA(file os.File) {

	// 移动到读取的位置
	_, _ = file.Seek(-int64(len(magic)+4), io.SeekCurrent)

	var num int
	var numBytes = make([]byte, 4)
	file.Read(numBytes)
	num = BytesToInt(numBytes)

	//读取
	// 移动到读取的位置
	_, _ = file.Seek(-int64(num+4), io.SeekCurrent)

	buffer := make([]byte, num)

	file.Read(buffer)

	cert = buffer
	caInit = true
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
