package main

import (
	// main "jieshao/automkcert"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

//集成测试

// 根证书生成
// mockgen -source=prompt.go
func TestRootCertMake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockPrompt(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	m.
		EXPECT().
		GenRootCert().
		Return(true)

	(&Guide{m}).Run()

	_, err := os.Stat("automkcert-root")
	if err != nil {
		t.Fail()
	}

}

// 加载ca
func mockAndLoadCa(t *testing.T, i int) *Guide {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockPrompt(ctrl)

	// 生成子证书
	m.
		EXPECT().
		RootMenu().
		Return(i)

	caInit = true
	g := (&Guide{m})

	//加载证书
	file, err := os.Open("automkcert-root")
	if err != nil {
		fmt.Println("无法打开文件：", err)
		return nil
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("无法获取文件信息：", err)
		return nil
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
		return nil
	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("无法读取文件内容：", err)
		return nil
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadFS(*file)
		return g
	}

	return nil

}

func TestSubCertMake(t *testing.T) {
	//选择1选项

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockPrompt(ctrl)

	// 生成子证书
	m.
		EXPECT().
		RootMenu().
		Return(0)

	m.
		EXPECT().
		InputHost().
		Return([]string{"localhost", "127.0.1.1", "127.0.0.1"})

	caInit = true
	g := (&Guide{m})

	//加载证书
	file, err := os.Open("automkcert-root")
	if err != nil {
		fmt.Println("无法打开文件：", err)
		t.Fail()

	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("无法获取文件信息：", err)
		t.Fail()

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
		t.Fail()

	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("无法读取文件内容：", err)
		t.Fail()
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadFS(*file)
	}
	//初始化
	g.Run()
}

func TestExportRootCertMake(t *testing.T) {
	//选择2选项
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockPrompt(ctrl)

	// 生成子证书
	m.
		EXPECT().
		RootMenu().
		Return(1)

	caInit = true
	g := (&Guide{m})

	//加载证书
	file, err := os.Open("automkcert-root")
	if err != nil {
		fmt.Println("无法打开文件：", err)
		t.Fail()

	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("无法获取文件信息：", err)
		t.Fail()

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
		t.Fail()

	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("无法读取文件内容：", err)
		t.Fail()
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadFS(*file)
	}
	g.Run()
}

// 生成授信网页
func TestExportCertClient(t *testing.T) {
	//选择2选项
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockPrompt(ctrl)

	// 生成子证书
	m.
		EXPECT().
		RootMenu().
		Return(3)

	caInit = true
	g := (&Guide{m})

	//加载证书
	file, err := os.Open("automkcert-root")
	if err != nil {
		fmt.Println("无法打开文件：", err)
		t.Fail()

	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("无法获取文件信息：", err)
		t.Fail()

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
		t.Fail()

	}

	// 创建一个缓冲区来读取内容
	buffer := make([]byte, offset)

	// 读取内容
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("无法读取文件内容：", err)
		t.Fail()
	}

	//判断是否有魔术字
	if bytes.Equal(buffer, magic) {
		loadFS(*file)
		g.Run()
	} else {
		t.Fatalf("需要先生成客户端")
	}
}
