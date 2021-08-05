package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// 调用test模块以检测文件是否被加密过
func execEncryptionTest(filePath string, testPath string) bool {
	command := exec.Command(testPath, filePath)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err := command.Run()
	if err != nil {
		panic(err)
	}
	output := command.Stdout.(*bytes.Buffer).String()
	if output[:4] == "true" {
		return true
	} else {
		fmt.Println("文件" + filePath + "未被加密，无需解密")
		return false
	}
}

// 对被检测到已加密的文件，进行一次读取并写入，解密后的文件加上 .block 后缀防止被加密
func reverseFile(fileName string, testPath string, impactPath string, content []byte) {
	if execEncryptionTest(fileName, testPath) {
		fmt.Println("正在解密：" + fileName)
		fmt.Println(impactPath)
		file, _ := os.Open(fileName)
		content, _ = ioutil.ReadAll(file)
		if content != nil {
			// 输出到带.block扩展名
			ioutil.WriteFile(fileName+".block", content, 777)
			os.Remove(fileName)
			// 用不在监控名单中的程序把扩展名改回去
			command := exec.Command(impactPath, fileName+".block")
			command.Stdout = &bytes.Buffer{}
			command.Stderr = &bytes.Buffer{}
			err := command.Run()
			if err != nil {
				panic(err)
			}
			output := command.Stdout.(*bytes.Buffer).String()
			fmt.Println(output)
		}
	}
}

func main() {
	args := os.Args
	var content []byte
	pwd, _ := os.Getwd()
	fmt.Println("已指定test模块路径为：" + args[1])
	// 对单个文件操作
	if len(args) != 3 {
		reverseFile(pwd+args[3], args[1], args[2], content)
	} else {
		// 对整个当前目录操作
		filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
			reverseFile(path, args[1], args[2], content)
			return nil
		})
	}
}
