package main

import "fmt"

type TokenType int

const (
	EOF = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	STAR
	SEMICOLON
)

func (tt TokenType) String() string {
	return [...]string{
		"EOF",
		"LEFT_PAREN", "RIGHT_PAREN",
		"LEFT_BRACE", "RIGHT_BRACE",
		"COMMA", "DOT", "MINUS", "PLUS", "STAR", "SEMICOLON",
	}[tt]
}

type Token struct {
	token_type TokenType
	lexeme string
	literal interface{}
}

func (t Token) String () string {
	literal_string := t.literal
	if t.literal == nil{
		literal_string = "null"
	}
	return fmt.Sprintf("%s %s %s", t.token_type.String(), t.lexeme, literal_string)
}