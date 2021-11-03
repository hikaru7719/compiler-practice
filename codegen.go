package main

import "fmt"

var incrementNumber = 0

func UniqueNum() int {
	return incrementNumber
}

func Increment() {
	incrementNumber += 1
}

func GenLval(node *Node) {
	if node.Kind != ND_LVAR {
		Error("代入の左辺値が変数ではありません")
	}

	fmt.Printf("	mov rax, rbp\n")
	fmt.Printf("	sub rax, %d\n", node.Offset)
	fmt.Printf("	push rax\n")
}

func Gen(node *Node) {
	switch node.Kind {
	case ND_FUNCTION:
		fmt.Printf("	call %s\n", node.FunctionName)
		return
	case ND_BLOCK:
		for _, child := range node.Statements {
			Gen(child)
			// TODO: popする処理が必要？
			// ref https://www.sigbus.info/compilerbook#%E3%82%B9%E3%83%86%E3%83%83%E3%83%9713-%E3%83%96%E3%83%AD%E3%83%83%E3%82%AF
		}
		return
	case ND_IF_ELSE:
		lelse := UniqueNum()
		Increment()
		lend := UniqueNum()
		Increment()

		Gen(node.Compare)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .Lelse%d\n", lelse)
		Gen(node.Then)
		fmt.Printf("	jmp .Lend%d\n", lend)
		fmt.Printf(".Lelse%d:\n", lelse)
		Gen(node.Else)
		fmt.Printf(".Lend%d:\n", lend)
		return
	case ND_IF:
		lend := UniqueNum()
		Increment()

		Gen(node.Compare)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .Lend%d\n", lend)
		Gen(node.Then)
		fmt.Printf(".Lend%d:\n", lend)
		return
	case ND_WHILE:
		lbegin := UniqueNum()
		Increment()
		lend := UniqueNum()
		Increment()

		fmt.Printf(".Lbegin%d:\n", lbegin)
		Gen(node.Compare)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .Lend%d\n", lend)
		Gen(node.Then)
		fmt.Printf("	jmp .Lbegin%d\n", lbegin)
		fmt.Printf(".Lend%d:\n", lend)
		return
	case ND_FOR:
		lbegin := UniqueNum()
		Increment()
		lend := UniqueNum()
		Increment()

		if node.Init != nil {
			Gen(node.Init)
		}

		fmt.Printf(".Lbegin%d:\n", lbegin)

		if node.Compare != nil {
			Gen(node.Compare)
		}

		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .Lend%d\n", lend)

		Gen(node.Then)

		if node.After != nil {
			Gen(node.After)
		}

		fmt.Printf("	jmp .Lbegin%d\n", lbegin)
		fmt.Printf(".Lend%d:\n", lend)
		return
	case ND_RETURN:
		Gen(node.Lhs)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov rsp, rbp\n")
		fmt.Printf("	pop rbp\n")
		fmt.Printf("	ret\n")
		return
	case ND_NUM:
		fmt.Printf("	push %d\n", node.Val)
		return
	case ND_LVAR:
		GenLval(node)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov rax, [rax]\n")
		fmt.Printf("	push rax\n")
		return
	case ND_ASSIGN:
		GenLval(node.Lhs)
		Gen(node.Rhs)

		fmt.Printf("	pop rdi\n")
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov [rax], rdi\n")
		fmt.Printf("	push rdi\n")
		return
	}

	Gen(node.Lhs)
	Gen(node.Rhs)

	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		fmt.Printf("	add rax, rdi\n")
	case ND_SUB:
		fmt.Printf("	sub rax, rdi\n")
	case ND_MUL:
		fmt.Printf("	imul rax, rdi\n")
	case ND_DIV:
		fmt.Printf("	cqo\n")
		fmt.Printf("	idiv rdi\n")
	case ND_EQUAL:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	sete al\n")
		fmt.Printf("	movzb rax, al\n")
	case ND_LESS:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setl al\n")
		fmt.Printf("	movzb rax, al\n")
	case ND_LESS_EQUAL:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setle al\n")
		fmt.Printf("	movzb rax, al\n")
	case ND_NOT_EQUAL:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setne al\n")
		fmt.Printf("	movzb rax, al\n")
	}

	fmt.Printf("	push rax\n")
}
