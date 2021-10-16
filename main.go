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
	token = Tokenize(os.Args[1])
	node := Expr()

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")

	fmt.Printf("	push rbo\n")
	fmt.Printf("	mov rbp, rsp\n")
	fmt.Printf("	sub rsp, 208\n")

	for i := 0; code[i] != nil; i++ {
		Gen(node)
		fmt.Printf("	pop rax\n")
	}

	fmt.Printf("	mov rsp\n")
	fmt.Printf("	pop rbp\n")
	fmt.Printf("	ret\n")
	os.Exit(0)
}
