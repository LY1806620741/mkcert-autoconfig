package main

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

//go:generate sh client/gen.sh
//go:embed dist/*x64 dist/index.html
var staticFs embed.FS

var selffs FS
var magic = []byte("selffs")

var caInit = false

func init() {
	exePath, _ := os.Executable()
	file, err := os.Open(exePath)
	if err != nil {
		fmt.Println("无法打开文件：", err)
		return
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("无法获取文件信息：", err)
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
		fmt.Println("无法移动到指定位置：", err)
		return
	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("无法读取文件内容：", err)
		return
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadFS(*file)
	} else {
		// division()
	}

}

// 加载额外文件系统
func loadFS(file os.File) {

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

	dec := gob.NewDecoder(bytes.NewReader(buffer))

	err := dec.Decode(&selffs)

	if err != nil {
		fmt.Println("反序列化失败", err)
	}

	caInit = true
}

// 分裂,创建携带根证书公私钥的母客户端
func division() {

	f, _ := selffs.Open(rootName)

	if f == nil {
		return
	}

	exePath, _ := os.Executable()

	// 打开源文件
	sourceFile, err := os.Open(exePath)
	if err != nil {
		fmt.Println("无法打开源文件：", err)
		return
	}
	defer sourceFile.Close()

	// 创建或打开目标文件
	destinationFile, err := os.OpenFile(execNameWithOutSuffix+"-root", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("无法创建或打开目标文件：", err)
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		fmt.Println("无法复制文件：", err)
		return
	}

	var buf bytes.Buffer = selffs.ToGob()

	err = binary.Write(destinationFile, binary.LittleEndian, buf.Bytes())
	if err != nil {
		fmt.Println("无法写入", err)
		return
	}

	err = binary.Write(destinationFile, binary.LittleEndian, IntToBytes(buf.Len()))
	if err != nil {
		fmt.Println("无法写入", err)
		return
	}

	err = binary.Write(destinationFile, binary.LittleEndian, magic)
	if err != nil {
		fmt.Println("无法写入", err)
		return
	}

}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
