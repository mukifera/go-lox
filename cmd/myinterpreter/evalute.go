package main

import (
	"math/big"
	"fmt"
)

type Evaluator struct {
	expressions []Expression
	current int
	values []interface{}
}

func NewEvaluator(expressions []Expression) *Evaluator {
	var evaluator Evaluator
	evaluator.expressions = expressions
	evaluator.current = 0
	return &evaluator
}

func (evaluator *Evaluator) Evaluate() {
	evaluator.values = nil
	for _, expression := range evaluator.expressions {
		evaluator.values = append(evaluator.values, evaluator.evaluateExpression(expression))
	}
}

func (evaluator *Evaluator) evaluateExpression(expression Expression) interface{} {
	switch expression.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return expression.literal
	case ExpressionTypeEnum.GROUPING:
		return evaluator.evaluateExpression(expression.children[0])
	case ExpressionTypeEnum.UNARY:
		return evaluator.evaluateUnaryExpression(expression)
	}
	return nil
}

func (evaluator *Evaluator) evaluateUnaryExpression(expression Expression) interface{} {
	value := evaluator.evaluateExpression(expression.children[0])
	switch expression.operator {
	case OperatorEnum.BANG:
		value = value == false || value == nil
	case OperatorEnum.MINUS:
		value = *evaluator.negateValue(value)
	}
	return value
}

func (evaluator *Evaluator) negateValue(value interface{}) *big.Float {
	switch literal := value.(type){
	case big.Float: return literal.Neg(&literal)
	default: return nil
	}
}

func (evaluator *Evaluator) StringifyValues() []string {
	var ret []string
	for _, value := range evaluator.values {
		str := ""
		switch literal := value.(type)	{
		case nil:		 		str = fmt.Sprintf("nil")
		case bool:	 		str = fmt.Sprintf("%v", literal)
		case string: 		str = fmt.Sprintf("%s", literal)
		case big.Float: str = fmt.Sprintf("%s", literal.String())
		}
		ret = append(ret, str)
	}
	return ret
}