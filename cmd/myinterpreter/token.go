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

	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
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
		"AND",
		"CLASS",
		"ELSE",
		"FALSE",
		"FOR",
		"FUN",
		"IF",
		"NIL",
		"OR",
		"PRINT",
		"RETURN",
		"SUPER",
		"THIS",
		"TRUE",
		"VAR",
		"WHILE",
	}[tt]
}

type Token struct {
	token_type TokenType
	lexeme     string
	literal    interface{}
}

func (t Token) StringLiteral() string {
	switch literal := t.literal.(type) {
	case int:
		return fmt.Sprintf("%d", literal)
	case string:
		return literal
	case big.Float:
		formatted := literal.String()
		if !strings.Contains(formatted, ".") {
			formatted += ".0"
		}
		return formatted
	}
	return "null"
}

func (t Token) String() string {
	token_string := fmt.Sprintf("%s %s ", t.token_type.String(), t.lexeme)
	token_string += t.StringLiteral()
	return token_string
}
