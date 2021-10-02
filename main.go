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
}

var token *Token

func Error(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func Consume(op string) bool {
	if str := token.Str[:1]; token.Kind != TK_RESERVED || str != op {
		return false
	}
	token = token.Next
	return true
}

func Expect(op string) {
	if str := token.Str[:1]; token.Kind != TK_RESERVED || str != op {
		Error("%sではありません", op)
	}
	token = token.Next
}

func ExpectNumber() int {
	if token.Kind != TK_NUM {
		Error("数ではありません")
	}
	val := token.Val
	token = token.Next
	return val
}

func atEOF() bool {
	return token.Kind == TK_EOF
}

func NewToken(kind TokenKind, cur *Token, str string) *Token {
	newToken := &Token{Kind: kind, Str: str}
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

		if s == '+' || s == '-' {
			cur = NewToken(TK_RESERVED, cur, string(s))
			current++
			continue
		}

		if unicode.IsDigit(s) {
			cur = NewToken(TK_NUM, cur, string(s))
			result, readed := strtol(p, current)
			cur.Val = result
			current += readed
			continue
		}

		Error("トークナイズできません")
	}
	NewToken(TK_EOF, cur, "$")
	return head.Next
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("引数の個数が正しくありません\n")
		os.Exit(1)
	}

	token = Tokenize(os.Args[1])

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("	mov rax, %d\n", ExpectNumber())

	for !atEOF() {

		if Consume("+") {
			fmt.Printf("	add rax, %d\n", ExpectNumber())
			continue
		}

		Expect("-")
		fmt.Printf("	sub rax, %d\n", ExpectNumber())
	}

	fmt.Printf("	ret\n")
	os.Exit(0)
}

func strtol(str string, current int) (result int, readed int) {
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
