package interpreter

import (
	"errors"
	"fmt"
	"golox/lox/expression"
	"golox/lox/reporter"
	"golox/lox/token"
)

type Interpreter struct {
	reporter *reporter.ErrorReporter
}

func NewInterpreter(reporter *reporter.ErrorReporter) *Interpreter {
	return &Interpreter{
		reporter: reporter,
	}
}

func (interp *Interpreter) Interpret(expr expression.Expression) error {
	val, err := interp.evaluate(expr)
	if err != nil {
		return err
	}
	fmt.Println(interp.stringify(val))
	return nil
}

func (interp *Interpreter) evaluate(expr expression.Expression) (interface{}, error) {
	switch v := expr.(type) {
	case *expression.Unary:
		return interp.evaluateUnaryExpr(v)
	case *expression.Binary:
		return interp.evaluateBinaryExpr(v)
	case *expression.Literal:
		return interp.evaluateLiteralExpr(v)
	case *expression.Grouping:
		return interp.evaluateGroupingExpr(v)
	default:
		fmt.Println("unknown expression type")
	}
	return nil, nil
}

func (interp *Interpreter) evaluateLiteralExpr(expr *expression.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (interp *Interpreter) evaluateUnaryExpr(expr *expression.Unary) (interface{}, error) {
	right, err := interp.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Bang:
		return !interp.isTruthy(right), nil
	case token.Minus:
		if err := interp.checkNumberOperand(expr.Operator, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			break
		}
		return -right.(float64), nil
	}

	return nil, nil
}

func (interp *Interpreter) evaluateBinaryExpr(expr *expression.Binary) (interface{}, error) {
	left, err := interp.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := interp.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Greater:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GreaterEqual:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.Less:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LessEqual:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.EqualEqual:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return interp.isEqual(left, right), nil
	case token.BangEqual:
		return !interp.isEqual(left, right), nil
	case token.Minus:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.Plus:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) + right.(float64), nil
	case token.Slash:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.Star:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.DotDot:
		if err := interp.checkStringOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(string) + right.(string), nil
	}

	return nil, nil
}

func (interp *Interpreter) evaluateGroupingExpr(expr *expression.Grouping) (interface{}, error) {
	return interp.evaluate(expr.Expr)
}

func (interp *Interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if boolVal, ok := val.(bool); ok {
		return !boolVal
	}
	return true
}

func (interp *Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	}

	return interp.checkBools(left, right) || interp.checkFloats(left, right) || interp.checkStrings(left, right)
}

func (interp *Interpreter) checkBools(left, right interface{}) bool {
	lval, lok := left.(bool)
	if !lok {
		return false
	}
	rval, rok := right.(bool)
	if !rok {
		return false
	}

	return lval == rval
}

func (interp *Interpreter) checkFloats(left, right interface{}) bool {
	lval, lok := left.(float64)
	if !lok {
		return false
	}
	rval, rok := right.(float64)
	if !rok {
		return false
	}

	return lval == rval
}

func (interp *Interpreter) checkStrings(left, right interface{}) bool {
	lval, lok := left.(string)
	if !lok {
		return false
	}
	rval, rok := right.(string)
	if !rok {
		return false
	}

	return lval == rval
}

func (interp *Interpreter) checkNumberOperand(operator *token.Token, operand interface{}) error {
	_, ok := operand.(float64)
	if !ok {
		return fmt.Errorf("operand for binary operator '%s' must be a number", operator.Type)
	}
	return nil
}

func (interp *Interpreter) checkNumberOperands(operator *token.Token, left, right interface{}) error {
	_, lok := left.(float64)
	if !lok {
		return fmt.Errorf("left operand for binary operator '%s' must be a number", operator.Type)
	}
	rval, rok := right.(float64)
	if !rok {
		return fmt.Errorf("right operand for binary operator '%s' must be a number", operator.Type)
	}

	// special case
	if operator.Type == token.Slash && rval == 0.0 {
		return errors.New("division by zero detected")

	}
	return nil
}

func (interp *Interpreter) checkStringOperands(operator *token.Token, left, right interface{}) error {
	_, lok := left.(string)
	if !lok {
		return fmt.Errorf("left operand for binary operator '%s' must be a string", operator.Type)
	}
	_, rok := right.(string)
	if !rok {
		return fmt.Errorf("right operand for binary operator '%s' must be a string", operator.Type)
	}
	return nil
}

func (interp *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "null"
	}
	return fmt.Sprintf("%#v", val)
}
