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
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	SLASH
	STRING
	NUMBER
)

func (tt TokenType) String() string {
	return [...]string{
		"EOF",
		"LEFT_PAREN", "RIGHT_PAREN",
		"LEFT_BRACE", "RIGHT_BRACE",
		"COMMA", "DOT", "MINUS", "PLUS", "STAR", "SEMICOLON",
		"EQUAL", "EQUAL_EQUAL",
		"BANG", "BANG_EQUAL",
		"LESS", "LESS_EQUAL",
		"GREATER", "GREATER_EQUAL",
		"SLASH",
		"STRING",
		"NUMBER",
	}[tt]
}

type Token struct {
	token_type TokenType
	lexeme string
	literal interface{}
}

func (t Token) String () string {
	token_string := fmt.Sprintf("%s %s ", t.token_type.String(), t.lexeme)
	switch t.literal.(type) {
	case int: token_string += fmt.Sprintf("%d", t.literal); break;
	case float64: token_string += fmt.Sprintf("%g", t.literal); break;
	case nil: token_string += "null"; break;
	}
	return token_string
}