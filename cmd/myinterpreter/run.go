package main

import (
	"fmt"
	"io"
	"os"
)

func RunExpressions(exprs []Expression) []error {
	return FRunExpressions(os.Stdout, exprs)
}

func RunExpression(expr Expression, scope map[string]interface{}) error {
	return FRunExpression(os.Stdout, expr, scope)
}

func FRunExpressions(writer io.Writer, exprs []Expression) []error {
	errs := make([]error, len(exprs))
	scope := make(map[string]interface{})
	for index, expr := range exprs {
		err := FRunExpression(writer, expr, scope)
		errs[index] = err
		if err != nil {
			break
		}
	}
	return errs
}

func FRunExpression(writer io.Writer, expr Expression, scope map[string]interface{}) error {
	switch expr.operator {
	case OperatorEnum.PRINT:
		value, err := EvaluateExpression(expr.children[0], scope)
		if err != nil {
			return err
		}
		str := StringifyEvaluationValue(value)
		fmt.Fprintln(writer, str)
	case OperatorEnum.VAR:
		var variable string
		var value interface{}
		var err error

		if expr.children[0].expression_type == ExpressionTypeEnum.IDENTIFIER {
			variable = expr.children[0].StringLiteral()
			value = nil
			err = nil
		} else {
			variable = expr.children[0].children[0].StringLiteral()
			value, err = EvaluateExpression(expr.children[0].children[1], scope)
			if err != nil {
				return err
			}
		}

		scope[variable] = value

	default:
		_, err := EvaluateExpression(expr, scope)
		if err != nil {
			return err
		}
	}
	return nil
}
