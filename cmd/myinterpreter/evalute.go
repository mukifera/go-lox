package main

type Evaluator struct {
	expressions []Expression
	current int
}

func NewEvaluator(expressions []Expression) *Evaluator {
	var evaluator Evaluator
	evaluator.expressions = expressions
	evaluator.current = 0
	return &evaluator
}

func (evaluator *Evaluator) Evaluate() []interface{} {
	var values []interface{}
	for _, expression := range evaluator.expressions {
		values = append(values, evaluator.evaluateExpression(expression))
	}
	return values
}

func (evaluator *Evaluator) evaluateExpression(expression Expression) interface{} {
	switch expression.expression_type {
	case ExpressionTypeEnum.LITERAL:
		return expression.literal
	}
	return nil
}