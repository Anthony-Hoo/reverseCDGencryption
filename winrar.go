package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// 调用test模块以检测文件是否被加密过
func execEncryptionTest(filePath string, testPath string) bool {
	command := exec.Command(testPath, filePath)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err := command.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(command.Stderr.(*bytes.Buffer).String())
	}
	output := command.Stdout.(*bytes.Buffer).String()
	if output[:4] == "true" {
		return true
	} else {
		fmt.Println("文件未被加密，无需解密")
		return false
	}
}

// 对被检测到已加密的文件，进行一次读取并写入，解密后的文件加上 .block 后缀防止被加密
func reverseFile(fileName string, testPath string, content []byte) {
	if execEncryptionTest(fileName, testPath) {
		fmt.Println("正在解密：" + fileName)
		file, _ := os.Open(fileName)
		content, _ = ioutil.ReadAll(file)
		if content != nil {
			ioutil.WriteFile(fileName+".block", content, 777)
		}
	}
}

func main() {
	args := os.Args
	var content []byte
	pwd, _ := os.Getwd()
	// 对单个文件操作
	if len(args) != 2 {
		fmt.Println("已指定test模块路径为：" + args[1])
		reverseFile(pwd+args[2], args[1], content)
	} else {
		// 对整个当前目录操作
		fmt.Println("已指定test模块路径为：" + args[1])
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			panic(err)
		}
		for i := range fileInfoList {
			reverseFile(fileInfoList[i].Name(), args[1], content)
		}
	}
}
