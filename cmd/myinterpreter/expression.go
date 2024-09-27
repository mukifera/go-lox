package main

import (
	"fmt"
	"math/big"
	"strings"
)

type ExpressionType int

var ExpressionTypeEnum = struct {
	UNDEFINED  ExpressionType
	LITERAL    ExpressionType
	UNARY      ExpressionType
	BINARY     ExpressionType
	GROUPING   ExpressionType
	IDENTIFIER ExpressionType
	BUILTIN    ExpressionType
}{
	UNDEFINED:  0,
	LITERAL:    1,
	UNARY:      2,
	BINARY:     3,
	GROUPING:   4,
	IDENTIFIER: 5,
	BUILTIN:    6,
}

type Operator int

var OperatorEnum = struct {
	UNDEFINED     Operator
	MINUS         Operator
	PLUS          Operator
	STAR          Operator
	EQUAL         Operator
	EQUAL_EQUAL   Operator
	BANG          Operator
	BANG_EQUAL    Operator
	LESS          Operator
	LESS_EQUAL    Operator
	GREATER       Operator
	GREATER_EQUAL Operator
	SLASH         Operator
	PRINT         Operator
	VAR           Operator
}{
	UNDEFINED:     0,
	MINUS:         1,
	PLUS:          2,
	STAR:          3,
	EQUAL:         4,
	EQUAL_EQUAL:   5,
	BANG:          6,
	BANG_EQUAL:    7,
	LESS:          8,
	LESS_EQUAL:    9,
	GREATER:       10,
	GREATER_EQUAL: 11,
	SLASH:         12,
	PRINT:         13,
	VAR:           14,
}

func (o *Operator) StringSymbol() string {
	switch *o {
	case OperatorEnum.BANG:
		return "!"
	case OperatorEnum.MINUS:
		return "-"
	case OperatorEnum.STAR:
		return "*"
	case OperatorEnum.SLASH:
		return "/"
	case OperatorEnum.PLUS:
		return "+"
	case OperatorEnum.LESS:
		return "<"
	case OperatorEnum.LESS_EQUAL:
		return "<="
	case OperatorEnum.GREATER:
		return ">"
	case OperatorEnum.GREATER_EQUAL:
		return ">="
	case OperatorEnum.EQUAL_EQUAL:
		return "=="
	case OperatorEnum.BANG_EQUAL:
		return "!="
	case OperatorEnum.EQUAL:
		return "="
	case OperatorEnum.PRINT:
		return "print"
	case OperatorEnum.VAR:
		return "var"
	}
	return ""
}

type Expression struct {
	expression_type ExpressionType
	operator        Operator
	literal         interface{}
	children        []Expression
}

func NewBinaryExpression(left Expression, right Expression, operator Operator) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.BINARY
	ret.operator = operator
	ret.literal = nil
	ret.children = []Expression{left, right}
	return ret
}

func NewUnaryExpression(expr Expression, operator Operator) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.UNARY
	ret.operator = operator
	ret.literal = nil
	ret.children = []Expression{expr}
	return ret
}

func NewLiteralExpression(literal interface{}) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.LITERAL
	ret.operator = OperatorEnum.UNDEFINED
	ret.literal = literal
	ret.children = []Expression{}
	return ret
}

func NewGroupingExpression(children ...Expression) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.GROUPING
	ret.operator = OperatorEnum.UNDEFINED
	ret.literal = nil
	ret.children = children
	return ret
}

func NewIdentifierExpression(lexeme string) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.IDENTIFIER
	ret.operator = OperatorEnum.UNDEFINED
	ret.literal = lexeme
	ret.children = []Expression{}
	return ret
}

func NewUndefinedExpression() Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.UNDEFINED
	ret.operator = OperatorEnum.UNDEFINED
	ret.literal = nil
	ret.children = []Expression{}
	return ret
}

func NewBuiltinExpression(expr Expression, operator Operator) Expression {
	var ret Expression
	ret.expression_type = ExpressionTypeEnum.BUILTIN
	ret.operator = operator
	ret.literal = nil
	ret.children = []Expression{expr}
	return ret
}

func (e *Expression) String() string {
	switch e.expression_type {
	case ExpressionTypeEnum.UNDEFINED:
		return ""
	case ExpressionTypeEnum.LITERAL, ExpressionTypeEnum.IDENTIFIER:
		return e.StringLiteral()
	case ExpressionTypeEnum.GROUPING:
		return fmt.Sprintf("(group %s)", e.children[0].String())
	default:
		str := fmt.Sprintf("(%s", e.operator.StringSymbol())
		for _, child := range e.children {
			str = str + " " + child.String()
		}
		str = str + ")"
		return str
	}
}

func (e *Expression) StringLiteral() string {
	formatted := ""
	switch t := e.literal.(type) {
	case int:
		return fmt.Sprintf("%d", t)
	case string:
		return t
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
