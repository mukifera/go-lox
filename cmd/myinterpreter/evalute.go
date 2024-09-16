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
	case ExpressionTypeEnum.BINARY:
		return evaluator.evaluateBinaryExpression(expression)
	}
	return nil
}

func (evaluator *Evaluator) evaluateUnaryExpression(expression Expression) interface{} {
	value := evaluator.evaluateExpression(expression.children[0])
	switch expression.operator {
	case OperatorEnum.BANG:
		return value == false || value == nil
	case OperatorEnum.MINUS:
		number := evaluator.assertNumber(value)
		return *number.Neg(&number)
	}
	return nil
}

func (evaluator *Evaluator) evaluateBinaryExpression(expression Expression) interface{} {
	left_value := evaluator.evaluateExpression(expression.children[0])
	right_value := evaluator.evaluateExpression(expression.children[1])
	switch expression.operator {
	case OperatorEnum.STAR:
		left_number := evaluator.assertNumber(left_value)
		right_number := evaluator.assertNumber(right_value)
		return *left_number.Mul(&left_number, &right_number)
	case OperatorEnum.SLASH:
		left_number := evaluator.assertNumber(left_value)
		right_number := evaluator.assertNumber(right_value)
		return *left_number.Quo(&left_number, &right_number)
	case OperatorEnum.PLUS:
		left_number := evaluator.assertNumber(left_value)
		right_number := evaluator.assertNumber(right_value)
		return *left_number.Add(&left_number, &right_number)
	case OperatorEnum.MINUS:
		left_number := evaluator.assertNumber(left_value)
		right_number := evaluator.assertNumber(right_value)
		return *left_number.Sub(&left_number, &right_number)
	}
	return nil
}

func (evaluator *Evaluator) assertNumber(value interface{}) big.Float {
	var ret big.Float
	switch literal := value.(type){
	case big.Float: return literal
	}
	return ret
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