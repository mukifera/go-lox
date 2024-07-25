package main

import (
	"fmt"
	"math/big"
	"strings"
)

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
	IDENTIFIER
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
		"IDENTIFIER",
	}[tt]
}

type Token struct {
	token_type TokenType
	lexeme string
	literal interface{}
}

func (t Token) String () string {
	token_string := fmt.Sprintf("%s %s ", t.token_type.String(), t.lexeme)
	switch literal := t.literal.(type) {
	case int:
		token_string += fmt.Sprintf("%d", literal); break;
	case string:
		token_string += literal
	case big.Float:
		formatted := literal.String()
		if !strings.Contains(formatted, ".") {
			formatted += ".0"
		}
		token_string += formatted;
		break;
	case nil:
		token_string += "null"; break;
	}
	return token_string
}