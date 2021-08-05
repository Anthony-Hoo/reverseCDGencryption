package main

import (
	"fmt"
	"os"
	"strings"
)

func impaction(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".block")
	os.Rename(filePath, newPath)
	fmt.Println(newPath + " 转换成功")
}

func main() {
	args := os.Args
	impaction(args[1])
}
