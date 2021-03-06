package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("引数の個数が正しくありません\n")
		os.Exit(1)
	}
	userInput = os.Args[1]
	token = Tokenize(userInput)
	Program()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	fmt.Printf("	push rbp\n")
	fmt.Printf("	mov rbp, rsp\n")
	fmt.Printf("	sub rsp, 208\n")

	for i := 0; code[i] != nil; i++ {
		Gen(code[i])
		fmt.Printf("	pop rax\n")
	}

	fmt.Printf("	mov rsp, rbp\n")
	fmt.Printf("	pop rbp\n")
	fmt.Printf("	ret\n")
	os.Exit(0)
}
