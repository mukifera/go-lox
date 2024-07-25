package main

import (
	"fmt"
	"math/big"
	"strings"
)

type Parser struct {
	tokens []Token
	expressions []Expression
}

func NewParser(tokens []Token) *Parser {
	var parser Parser
	parser.tokens = tokens
	return &parser
}

func (parser *Parser) Parse() {
	for _, token := range parser.tokens {
		var expr Expression
		switch token.token_type {
		case EOF: continue;
		case TRUE:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = true;
			break;
		case FALSE:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = false;
			break;
		case NIL:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = nil;
			break;
		case NUMBER, STRING:
			expr.expression_type = ExpressionTypeEnum.LITERAL;
			expr.literal = token.literal;
			break;
		}
		parser.expressions = append(parser.expressions, expr)
	}
}

func StringLiteral(literal interface{}) string {
	formatted := ""
	switch t := literal.(type) {
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

func (parser *Parser) StringifyExpressions() string {
	str := ""
	for _, expr := range parser.expressions {
		switch expr.expression_type {
		case ExpressionTypeEnum.LITERAL:
			str += StringLiteral(expr.literal)
		}
	}
	return str
}