package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TokenKind int

const (
	TK_RESERVED TokenKind = iota + 1
	TK_IDENT
	TK_NUM
	TK_EOF
	TK_RETRUN
	TK_IF
	TK_ELSE
	TK_WHILE
	TK_FOR
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
	ND_ASSIGN                         // assign
	ND_LVAR                           // local variable
	ND_RETURN                         // retrun
	ND_IF                             // if
	ND_IF_ELSE                        // if else
	ND_WHILE                          // while
	ND_FOR                            // for
	ND_BLOCK                          // {}
	ND_FUNCTION                       // function
)

type Node struct {
	Kind NodeKind
	Lhs  *Node
	Rhs  *Node

	// if, while
	Compare *Node
	Then    *Node
	Else    *Node

	// for
	Init  *Node
	After *Node

	// block
	Statements []*Node

	// function
	FunctionName string

	Val    int
	Offset int
}

type LocalVar struct {
	Name   string
	Len    int
	Offset int
	Next   *LocalVar
}

var token *Token
var userInput string
var code [100]*Node
var locals *LocalVar

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

func ConsumeIdent() *Token {
	if token.Kind != TK_IDENT {
		return nil
	}
	returnToken := token
	token = token.Next
	return returnToken
}

func ConsumeKind(kind TokenKind) bool {
	if kind != token.Kind {
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

func IsAlnum(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || (c == '_')
}

func Ident(p string, current int) (string, int) {
	var result strings.Builder
	readed := 0
	for IsAlnum(rune(p[current+readed])) {
		result.WriteRune(rune(p[current+readed]))
		readed++
	}
	return result.String(), readed
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

		// return token
		if s == 'r' && len(p[current:]) >= 6 {
			if str := p[current : current+6]; str == "return" && !IsAlnum(rune(p[current+6])) {
				cur = NewToken(TK_RETRUN, cur, str, current)
				current += 6
				continue
			}
		}

		// if token
		if s == 'i' && len(p[current:]) >= 2 {
			if str := p[current : current+2]; str == "if" && !IsAlnum(rune(p[current+2])) {
				cur = NewToken(TK_IF, cur, str, current)
				current += 2
				continue
			}
		}

		// else token
		if s == 'e' && len(p[current:]) >= 4 {
			if str := p[current : current+4]; str == "else" && !IsAlnum(rune(p[current+4])) {
				cur = NewToken(TK_ELSE, cur, str, current)
				current += 4
				continue
			}
		}

		// while token
		if s == 'w' && len(p[current:]) >= 5 {
			if str := p[current : current+5]; str == "while" && !IsAlnum(rune(p[current+5])) {
				cur = NewToken(TK_WHILE, cur, str, current)
				current += 5
				continue
			}
		}

		// for token
		if s == 'f' && len(p[current:]) >= 3 {
			if str := p[current : current+3]; str == "for" && !IsAlnum(rune(p[current+3])) {
				cur = NewToken(TK_FOR, cur, str, current)
				current += 3
				continue
			}
		}

		if 'a' <= s && s <= 'z' {
			ident, readed := Ident(p, current)
			cur = NewToken(TK_IDENT, cur, ident, current)
			current += readed
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

		if s == '+' || s == '-' || s == '*' || s == '/' || s == '(' || s == ')' || s == ';' || s == '{' || s == '}' {
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

func FindLocalVar(token *Token) *LocalVar {
	for localVar := locals; localVar != nil; localVar = localVar.Next {
		if localVar.Len == token.Len && localVar.Name == token.Str {
			return localVar
		}
	}
	return nil
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func NewNodeIfElse(compare *Node, then *Node, el *Node) *Node {
	return &Node{
		Kind:    ND_IF_ELSE,
		Compare: compare,
		Then:    then,
		Else:    el,
	}
}

func NewNodeIf(compare *Node, then *Node) *Node {
	return &Node{
		Kind:    ND_IF,
		Compare: compare,
		Then:    then,
	}
}

func NewNodeWhile(compare *Node, then *Node) *Node {
	return &Node{
		Kind:    ND_WHILE,
		Compare: compare,
		Then:    then,
	}
}

func NewNodeBlock(statements []*Node) *Node {
	return &Node{
		Kind:       ND_BLOCK,
		Statements: statements,
	}
}

func NewNodeFunction(functionName string) *Node {
	return &Node{
		Kind:         ND_FUNCTION,
		FunctionName: functionName,
	}
}

func NewNodeFor(init *Node, compare *Node, after *Node, then *Node) *Node {
	return &Node{
		Kind:    ND_FOR,
		Init:    init,
		Compare: compare,
		After:   after,
		Then:    then,
	}
}

func NewNodeNum(val int) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
	}
}

func Program() {
	i := 0
	for !AtEOF() {
		code[i] = Stmt()
		i++
	}
	code[i] = nil
}

func Stmt() *Node {
	var node *Node
	if ConsumeKind(TK_RETRUN) {
		node = NewNode(ND_RETURN, Expr(), nil)
		if !Consume(";") {
			ErrorAt(token.Pos, "';'ではないトークンです")
		}
	} else if ConsumeKind(TK_IF) {
		if Consume("(") {
			compare := Expr()
			if Consume(")") {
				then := Stmt()
				if ConsumeKind(TK_ELSE) {
					el := Stmt()
					node = NewNodeIfElse(compare, then, el)
				} else {
					node = NewNodeIf(compare, then)
				}
			} else {
				ErrorAt(token.Pos, "')'ではないトークンです")
			}
		} else {
			ErrorAt(token.Pos, "'('ではないトークンです")
		}
	} else if ConsumeKind(TK_WHILE) {
		if Consume("(") {
			compare := Expr()
			if Consume(")") {
				then := Stmt()
				node = NewNodeWhile(compare, then)
			} else {
				ErrorAt(token.Pos, "')'ではないトークンです")
			}
		} else {
			ErrorAt(token.Pos, "'('ではないトークンです")
		}
	} else if ConsumeKind(TK_FOR) {
		var init, compare, after *Node
		if Consume("(") {
			for i := 0; i < 3; i++ {
				if Consume(";") {
					continue
				}
				if i == 0 {
					init = Expr()
				}
				if i == 1 {
					compare = Expr()
				}
				if i == 2 {
					after = Expr()
				}
				if Consume(";") {
					continue
				} else {
					ErrorAt(token.Pos, "';'ではないトークンです")
				}
			}
			if Consume(")") {
				then := Stmt()
				node = NewNodeFor(init, compare, after, then)
			} else {
				ErrorAt(token.Pos, "')'ではないトークンです")
			}
		} else {
			ErrorAt(token.Pos, "'('ではないトークンです")
		}
	} else if Consume("{") {
		statements := []*Node{}
		for !Consume("}") {
			statements = append(statements, Stmt())
		}
		node = NewNodeBlock(statements)
	} else {
		node = Expr()
		if !Consume(";") {
			ErrorAt(token.Pos, "';'ではないトークンです")
		}
	}
	return node
}

func Expr() *Node {
	return Assign()
}

func Assign() *Node {
	node := Equality()
	if Consume("=") {
		node = NewNode(ND_ASSIGN, node, Assign())
	}
	return node
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

	if tok := ConsumeIdent(); tok != nil {
		if Consume("(") {
			if Consume(")") {
				return NewNodeFunction(tok.Str)
			}
		}

		node := NewNode(ND_LVAR, nil, nil)
		if localVar := FindLocalVar(tok); localVar != nil {
			node.Offset = localVar.Offset
		} else if locals != nil {
			lVar := &LocalVar{
				Name:   tok.Str,
				Len:    tok.Len,
				Offset: locals.Offset + 8,
				Next:   locals,
			}
			node.Offset = locals.Offset + 8
			locals = lVar
		} else {
			lVar := &LocalVar{
				Name:   tok.Str,
				Len:    tok.Len,
				Offset: 8,
				Next:   locals,
			}
			node.Offset = 8
			locals = lVar
		}

		return node
	}

	return NewNodeNum(ExpectNumber())
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
