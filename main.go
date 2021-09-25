package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		panic("引数の個数が正しくありません")
	}

	n, _ := strconv.Atoi(os.Args[1])
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("	mov rax, %d\n", n)
	fmt.Printf("	ret\n")
}
