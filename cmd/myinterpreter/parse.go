package main

import "fmt"

func Parse(tokens []Token) {
	for _, token := range tokens {
		switch token.token_type {
		case EOF: break;
		case TRUE: fmt.Println("true"); break;
		case FALSE: fmt.Println("false"); break;
		case NUMBER: fmt.Println(token.StringLiteral()); break;
		default: fmt.Println("nil"); break;
		}
	}
}