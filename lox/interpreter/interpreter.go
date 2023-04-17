package interpreter

import (
	"errors"
	"fmt"
	"golox/lox/environment"
	"golox/lox/expression"
	"golox/lox/reporter"
	"golox/lox/statement"
	"golox/lox/token"
)

type Interpreter struct {
	reporter *reporter.ErrorReporter
	env      *environment.Environment
	repl     bool
}

func NewInterpreter(reporter *reporter.ErrorReporter) *Interpreter {
	return &Interpreter{
		reporter: reporter,
		env:      environment.NewEnvironment(nil),
		repl:     false,
	}
}

func (interp *Interpreter) Interpret(statements []statement.Stmt, repl bool) error {
	interp.repl = repl
	for _, stmt := range statements {
		err := interp.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execute(stmt statement.Stmt) error {
	var err error
	switch v := stmt.(type) {
	case *statement.PrintStmt:
		err = interp.executePrintStmt(v)
	case *statement.ExpressionStmt:
		err = interp.executeExprStmt(v)
	case *statement.VarStmt:
		err = interp.executeVarStmt(v)
	case *statement.BlockStmt:
		err = interp.executeBlockStmt(v)
	}
	return err
}

func (interp *Interpreter) executeVarStmt(stmt *statement.VarStmt) error {
	if _, defined := interp.env.Lookup(stmt.Name); defined {
		err := fmt.Errorf("variable named '%s' already exists", stmt.Name.Lexeme)
		interp.reporter.ReportAtLocation(err, "TODO", "", stmt.Name.Line, 0, 0)
		return err
	}

	var (
		val interface{}
		err error
	)
	if stmt.Initializer != nil {
		if val, err = interp.evaluate(stmt.Initializer); err != nil {
			return err
		}
	}
	interp.env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (interp *Interpreter) executePrintStmt(stmt *statement.PrintStmt) error {
	val, err := interp.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", val)
	return nil
}

func (interp *Interpreter) executeExprStmt(stmt *statement.ExpressionStmt) error {
	val, err := interp.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	if interp.repl {
		fmt.Println("eval:", val)
	}
	return nil
}

func (interp *Interpreter) executeBlockStmt(stmt *statement.BlockStmt) error {
	env := environment.NewEnvironment(interp.env)
	return interp.executeBlock(stmt.Statements, env)
}

func (interp *Interpreter) executeBlock(statements []statement.Stmt, env *environment.Environment) error {
	previous := interp.env
	defer interp.setEnvironment(previous)
	interp.env = env

	for _, s := range statements {
		if err := interp.execute(s); err != nil {
			return err
		}
	}

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
	case *expression.Variable:
		return interp.evaluateVariableExpr(v)
	case *expression.Assign:
		return interp.evaluateAssignExpr(v)
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
		return !isTruthy(right), nil
	case token.Minus:
		if err := checkNumberOperand(expr.Operator, right); err != nil {
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
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GreaterEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.Less:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LessEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.EqualEqual:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return isEqual(left, right), nil
	case token.BangEqual:
		return !isEqual(left, right), nil
	case token.Minus:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.Plus:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) + right.(float64), nil
	case token.Slash:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.Star:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.ReportAtLocation(err, "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.DotDot:
		if err := checkStringOperands(expr.Operator, left, right); err != nil {
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

func (interp *Interpreter) evaluateVariableExpr(expr *expression.Variable) (interface{}, error) {
	val, found := interp.env.Lookup(expr.Name)
	if !found {
		return nil, fmt.Errorf("undefined variable '%s'", expr.Name.Lexeme)
	}

	return val, nil
}

func (interp *Interpreter) evaluateAssignExpr(expr *expression.Assign) (interface{}, error) {
	val, err := interp.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if ok := interp.env.Assign(expr.Name, val); !ok {
		err = fmt.Errorf("undefined variable '%s'", expr.Name.Lexeme)
		return nil, err
	}
	return val, nil
}

func (interp *Interpreter) setEnvironment(env *environment.Environment) {
	interp.env = env
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if boolVal, ok := val.(bool); ok {
		return !boolVal
	}
	return true
}

func isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	}

	return checkBools(left, right) || checkFloats(left, right) || checkStrings(left, right)
}

func checkBools(left, right interface{}) bool {
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

func checkFloats(left, right interface{}) bool {
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

func checkStrings(left, right interface{}) bool {
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

func checkNumberOperand(operator *token.Token, operand interface{}) error {
	_, ok := operand.(float64)
	if !ok {
		return fmt.Errorf("operand for binary operator '%s' must be a number", operator.Type)
	}
	return nil
}

func checkNumberOperands(operator *token.Token, left, right interface{}) error {
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

func checkStringOperands(operator *token.Token, left, right interface{}) error {
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

func stringify(val interface{}) string {
	if val == nil {
		return "null"
	}
	return fmt.Sprintf("%#v", val)
}
