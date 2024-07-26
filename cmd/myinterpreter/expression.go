package main

import (
	"fmt"
	"math/big"
	"strings"
)

type ExpressionType int

var ExpressionTypeEnum = struct {
	UNDEFINED  ExpressionType
	LITERAL 	 ExpressionType
	UNARY			 ExpressionType
	BINARY		 ExpressionType
	GROUPING	 ExpressionType
}{
	UNDEFINED: 0,
	LITERAL:	 1,
	UNARY:		 2,
	BINARY:		 3,
	GROUPING:	 4,
}

type Operator int

var OperatorEnum = struct {
	UNDEFINED				Operator
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
	UNDEFINED:			0,
	MINUS:					1,
	PLUS:						2,
	STAR:						3,
	EQUAL:					4,
	EQUAL_EQUAL:		5,
	BANG:						6,
	BANG_EQUAL:			7,
	LESS:						8,
	LESS_EQUAL:			9,
	GREATER:				10,
	GREATER_EQUAL:	11,
	SLASH:					12,
}

func (o *Operator) StringSymbol() string {
	switch *o {
	case OperatorEnum.BANG: return "!"
	case OperatorEnum.MINUS: return "-"
	case OperatorEnum.STAR: return "*"
	case OperatorEnum.SLASH: return "/"
	case OperatorEnum.PLUS: return "+"
	case OperatorEnum.LESS: return "<"
	case OperatorEnum.LESS_EQUAL: return "<="
	case OperatorEnum.GREATER: return ">"
	case OperatorEnum.GREATER_EQUAL: return ">="
	case OperatorEnum.EQUAL_EQUAL: return "=="
	case OperatorEnum.BANG_EQUAL: return "!="
	}
	return ""
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
	case ExpressionTypeEnum.UNARY:
		return fmt.Sprintf("(%s %s)", e.operator.StringSymbol(), e.children[0].String())
	case ExpressionTypeEnum.BINARY:
		return fmt.Sprintf("(%s %s %s)", e.operator.StringSymbol(), e.children[0].String(), e.children[1].String())
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