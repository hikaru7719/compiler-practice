package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("引数の個数が正しくありません\n")
		os.Exit(1)
	}

	p := os.Args[1]
	n, i := strtol(p, 0)
	pp := p[i:]
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("	mov rax, %d\n", n)

	current := 0
	for len(pp) > current {
		s := string(pp[current])
		if s == "+" {
			current++
			num, i := strtol(string(pp), current)
			fmt.Printf("	add rax, %d\n", num)
			current += i
			continue
		}
		if s == "-" {
			current++
			num, i := strtol(string(pp), current)
			fmt.Printf("	sub rax, %d\n", num)
			current += i
			continue
		}

		fmt.Printf("予期しない文字列です, %s", string(s))
		os.Exit(1)
	}

	fmt.Printf("	ret\n")
	os.Exit(0)
}

func strtol(str string, current int) (int, int) {
	result := 0
	readed := 0
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
