package main

import(
	"embed"
	"fmt"
	"os"
	"io"
	"bytes"
	"encoding/binary"
)

var selffs embed.FS
var magic = []byte("selffs")

func init(){
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
	numBytes:=len(magic)
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
	if (bytes.Compare(buffer,magic)==0){
		fmt.Print("文件系统存在,进行加载")
	}else{
		//初始化
		// _, err = file.Seek(0, io.SeekStart)
		// if err != nil {
		// 	fmt.Println("无法移动到指定位置：", err)
		// 	return
		// }
		// filebytes := make([]byte, offset)
		// _, err = file.Read(filebytes)
		// if err != nil {
		// 	fmt.Println("无法读取文件内容：", err)
		// 	return
		// }

		// //创建文件
		// file, err := os.Create("test")
		// if err != nil {
		// 	fmt.Println("无法创建或打开文件：", err)
		// 	return
		// }
		// defer file.Close()

		// err = binary.Write(file, binary.LittleEndian, filebytes)
		// if err != nil {
		// 	fmt.Println("无法写入", err)
		// 	return
		// }

		// 打开源文件
		sourceFile, err := os.Open(exePath)
		if err != nil {
			fmt.Println("无法打开源文件：", err)
			return
		}
		defer sourceFile.Close()

		// 创建或打开目标文件
		destinationFile, err := os.OpenFile("test", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

		err = binary.Write(destinationFile, binary.LittleEndian, magic)
		if err != nil {
			fmt.Println("无法写入", err)
			return
		}


	}

	fmt.Printf(exePath)
}