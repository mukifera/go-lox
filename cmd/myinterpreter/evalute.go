package main

import "fmt"

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
	}
	return nil
}

func (evaluator *Evaluator) StringifyValues() []string {
	var ret []string
	for _, value := range evaluator.values {
		str := ""
		switch value.(type)	{
		case bool:
			str = fmt.Sprintf("%v", value)
		case nil:
			str = fmt.Sprintf("nil")
		}
		ret = append(ret, str)
	}
	return ret
}