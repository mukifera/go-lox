package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func RunExpressions(exprs []Expression) error {
	return FRunExpressions(os.Stdout, exprs)
}

func RunExpression(expr Expression, context []map[string]interface{}) error {
	return FRunExpression(os.Stdout, expr, context)
}

func FRunExpressions(writer io.Writer, exprs []Expression) error {
	var err error = nil
	scope := make(map[string]interface{})
	context := []map[string]interface{}{scope}
	for _, expr := range exprs {
		sub_err := FRunExpression(writer, expr, context)
		err = errors.Join(err, sub_err)
		if err != nil {
			break
		}
	}
	return err
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

	case OperatorEnum.IF:
		condition := expr.children[0]
		body := expr.children[1]
		raw_value, err := EvaluateExpression(condition, context)
		if err != nil {
			return err
		}
		value := raw_value != false && raw_value != nil
		if value {
			err = FRunExpression(writer, body, context)
			if err != nil {
				return err
			}
		} else if len(expr.children) > 2 {
			err = FRunExpression(writer, expr.children[2], context)
			if err != nil {
				return err
			}
		}

	case OperatorEnum.WHILE:
		condition := expr.children[0]
		body := expr.children[1]
		for {
			raw_value, err := EvaluateExpression(condition, context)
			if err != nil {
				return err
			}
			value := raw_value != false && raw_value != nil
			if !value {
				break
			}
			err = FRunExpression(writer, body, context)
			if err != nil {
				return err
			}
		}

	case OperatorEnum.FOR:
		initial := expr.children[0]
		condition := expr.children[1]
		update := expr.children[2]
		body := expr.children[3]
		context = append(context, make(map[string]interface{}))
		err := FRunExpression(writer, initial, context)
		if err != nil {
			return err
		}
		for {
			raw_value, err := EvaluateExpression(condition, context)
			if err != nil {
				return err
			}
			value := raw_value != false && raw_value != nil
			if !value {
				break
			}
			err = FRunExpression(writer, body, context)
			if err != nil {
				return err
			}
			err = FRunExpression(writer, update, context)
			if err != nil {
				return err
			}
		}

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
