package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

func isEncrypted(filePath string) bool {
	file, _ := os.Open(filePath)
	buf := make([]byte, 24)
	file.Read(buf)
	// magic
	return hex.EncodeToString(buf) == "621423659d00630100000001452d536166654e6574000000"
}

func main() {
	args := os.Args
	if isEncrypted(args[1]) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
