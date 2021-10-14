package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type TokenKind int

const (
	TK_RESERVED TokenKind = iota + 1
	TK_NUM
	TK_EOF
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  int
	Str  string
	Pos  int
	Len  int
}

type NodeKind int

const (
	ND_ADD        NodeKind = iota + 1 // +
	ND_SUB                            // -
	ND_MUL                            // *
	ND_DIV                            // /
	ND_NUM                            // number
	ND_EQUAL                          // ==
	ND_NOT_EQUAL                      // !=
	ND_LESS                           // <
	ND_LESS_EQUAL                     // <=
)

type Node struct {
	Kind NodeKind
	Lhs  *Node
	Rhs  *Node
	Val  int
}

var token *Token
var userInput string

func Error(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func ErrorAt(current int, format string, args ...interface{}) {
	fmt.Printf("%s\n", userInput)
	fmt.Printf("%*s", current, " ")
	fmt.Printf("^ ")
	fmt.Printf(format, args...)
	fmt.Printf("\n")
	os.Exit(1)
}

func Consume(op string) bool {
	if str := token.Str[:token.Len]; token.Kind != TK_RESERVED || len(op) != token.Len || str != op {
		return false
	}
	token = token.Next
	return true
}

func Expect(op string) {
	if str := token.Str[:token.Len]; token.Kind != TK_RESERVED || len(op) != token.Len || str != op {
		Error("%sではありません", op)
	}
	token = token.Next
}

func ExpectNumber() int {
	if token.Kind != TK_NUM {
		ErrorAt(token.Pos, "数ではありません")
	}
	val := token.Val
	token = token.Next
	return val
}

func AtEOF() bool {
	return token.Kind == TK_EOF
}

func NewToken(kind TokenKind, cur *Token, str string, current int) *Token {
	newToken := &Token{Kind: kind, Str: str, Pos: current, Len: len(str)}
	cur.Next = newToken
	return newToken
}

func Tokenize(p string) *Token {
	var head Token
	var cur *Token = &head

	current := 0
	for len(p) > current {
		s := rune(p[current])

		if unicode.IsSpace(s) {
			current++
			continue
		}

		if s == '=' || s == '<' || s == '>' || s == '!' {
			word := string(p[current : current+2])
			if word == "==" || word == "<=" || word == ">=" || word == "!=" {
				cur = NewToken(TK_RESERVED, cur, word, current)
				current += 2
				continue
			}
			cur = NewToken(TK_RESERVED, cur, string(s), current)
			current++
			continue
		}

		if s == '+' || s == '-' || s == '*' || s == '/' || s == '(' || s == ')' {
			cur = NewToken(TK_RESERVED, cur, string(s), current)
			current++
			continue
		}

		if unicode.IsDigit(s) {
			cur = NewToken(TK_NUM, cur, string(s), current)
			result, readed := Strtol(p, current)
			cur.Val = result
			current += readed
			continue
		}

		ErrorAt(current, "トークナイズできません")
	}
	NewToken(TK_EOF, cur, "$", current)
	return head.Next
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func NewNodeNum(val int) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
	}
}

func Expr() *Node {
	return Equality()
}

func Equality() *Node {
	node := Relational()
	for {
		if Consume("==") {
			node = NewNode(ND_EQUAL, node, Relational())
		} else if Consume("!=") {
			node = NewNode(ND_NOT_EQUAL, node, Relational())
		} else {
			return node
		}
	}
}

func Relational() *Node {
	node := Add()
	for {
		if Consume("<=") {
			node = NewNode(ND_LESS_EQUAL, node, Add())
		} else if Consume("<") {
			node = NewNode(ND_LESS, node, Add())
		} else if Consume(">=") {
			node = NewNode(ND_LESS_EQUAL, Add(), node)
		} else if Consume(">") {
			node = NewNode(ND_LESS, Add(), node)
		} else {
			return node
		}
	}
}

func Add() *Node {
	node := Mul()
	for {
		if Consume("+") {
			node = NewNode(ND_ADD, node, Mul())
		} else if Consume("-") {
			node = NewNode(ND_SUB, node, Mul())
		} else {
			return node
		}
	}
}

func Mul() *Node {
	node := Unary()

	for {
		if Consume("*") {
			node = NewNode(ND_MUL, node, Unary())
		} else if Consume("/") {
			node = NewNode(ND_DIV, node, Unary())
		} else {
			return node
		}
	}
}

func Unary() *Node {
	if Consume("+") {
		return Primary()
	} else if Consume("-") {
		return NewNode(ND_SUB, NewNodeNum(0), Primary())
	} else {
		return Primary()
	}
}

func Primary() *Node {
	if Consume("(") {
		node := Expr()
		Expect(")")
		return node
	}
	return NewNodeNum(ExpectNumber())
}

func Gen(node *Node) {
	if node.Kind == ND_NUM {
		fmt.Printf("	push %d\n", node.Val)
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

	Gen(node)

	fmt.Printf("	pop rax\n")
	fmt.Printf("	ret\n")
	os.Exit(0)
}

func Strtol(str string, current int) (result int, readed int) {
	for len(str) > current+readed {
		pop := str[current : current+readed+1]
		num, err := strconv.ParseInt(pop, 10, 64)
		if err != nil {
			if result == 0 {
				return 0, readed
			}
			return result, readed
		}
		// nolint
		result = int(num)
		readed++
	}
	return result, readed
}
