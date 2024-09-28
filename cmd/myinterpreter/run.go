package main

import (
	"fmt"
	"io"
	"os"
)

func RunExpressions(exprs []Expression) []error {
	return FRunExpressions(os.Stdout, exprs)
}

func RunExpression(expr Expression, context []map[string]interface{}) error {
	return FRunExpression(os.Stdout, expr, context)
}

func FRunExpressions(writer io.Writer, exprs []Expression) []error {
	errs := make([]error, len(exprs))
	scope := make(map[string]interface{})
	context := []map[string]interface{}{scope}
	for index, expr := range exprs {
		err := FRunExpression(writer, expr, context)
		errs[index] = err
		if err != nil {
			break
		}
	}
	return errs
}

func FRunExpression(writer io.Writer, expr Expression, context []map[string]interface{}) error {
	switch expr.operator {
	case OperatorEnum.PRINT:
		value, err := EvaluateExpression(expr.children[0], context)
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
			value, err = EvaluateExpression(expr.children[0].children[1], context)
			if err != nil {
				return err
			}
		}

		context[len(context)-1][variable] = value

	default:
		if expr.expression_type == ExpressionTypeEnum.SCOPE {
			newScope := make(map[string]interface{})
			context = append(context, newScope)
			for _, child := range expr.children {
				err := FRunExpression(writer, child, context)
				if err != nil {
					return err
				}
			}
			context = context[:len(context)-1]
			break
		}

		_, err := EvaluateExpression(expr, context)
		if err != nil {
			return err
		}
	}
	return nil
}
