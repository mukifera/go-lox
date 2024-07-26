package main

import (
	"fmt"
	"math/big"
	"strings"
)

type ExpressionType int

var ExpressionTypeEnum = struct {
	LITERAL 	ExpressionType
	UNARY			ExpressionType
	BINARY		ExpressionType
	GROUPING	ExpressionType
}{
	LITERAL:	0,
	UNARY:		1,
	BINARY:		2,
	GROUPING:	3,
}

type Operator int

var OperatorEnum = struct {
	MINUS 					Operator
	PLUS 						Operator
	STAR 						Operator
	EQUAL 					Operator
	EQUAL_EQUAL 		Operator
	BANG 						Operator
	BANG_EQUAL 			Operator
	LESS 						Operator
	LESS_EQUAL 			Operator
	GREATER 				Operator
	GREATER_EQUAL 	Operator
	SLASH 					Operator
}{
	MINUS:					0,
	PLUS:						1,
	STAR:						2,
	EQUAL:					3,
	EQUAL_EQUAL:		4,
	BANG:						5,
	BANG_EQUAL:			6,
	LESS:						7,
	LESS_EQUAL:			8,
	GREATER:				9,
	GREATER_EQUAL:	10,
	SLASH:					11,
}

type Expression struct {
	expression_type ExpressionType
	operator Operator
	literal interface{}
	children []Expression
}

func (e *Expression) String() string {
	switch e.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return e.StringLiteral()
	case ExpressionTypeEnum.GROUPING:
		return fmt.Sprintf("(group %s)", e.children[0].String())
	}
	return ""
}

func (e *Expression) StringLiteral() string {
	formatted := ""
	switch t := e.literal.(type) {
	case int: return fmt.Sprintf("%d", t);
	case string: return t;
	case big.Float:
		formatted := t.String()
		if !strings.Contains(formatted, ".") {
			formatted += ".0"
		}
		return formatted
	case bool:
		if t {
			formatted += "true"
		} else {
			formatted += "false"
		}
		return formatted
	}
	return "nil"
}