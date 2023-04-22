package interpreter

import (
	"fmt"
	"golox/lox/environment"
	"golox/lox/expression"
	"golox/lox/reporter"
	"golox/lox/statement"
	"golox/lox/token"
	"time"
)

type Interpreter struct {
	reporter *reporter.ErrorReporter
	env      *environment.Environment
	globals  *environment.Environment
	repl     bool
}

func NewInterpreter(reporter *reporter.ErrorReporter) *Interpreter {
	interp := &Interpreter{
		reporter: reporter,
		env:      environment.NewEnvironment(nil),
		repl:     false,
	}

	interp.globals = interp.env

	interp.globals.Define("clock", NewLoxCallable(
		0,
		func(intrp *Interpreter, args []interface{}) (interface{}, error) {
			return time.Now().Unix(), nil
		},
		func() string {
			return "<native fn>"
		}),
	)

	return interp
}

func (interp *Interpreter) Interpret(statements []statement.Stmt, repl bool) error {
	interp.repl = repl
	for _, stmt := range statements {
		_, err := interp.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execute(stmt statement.Stmt) (interface{}, error) {
	switch v := stmt.(type) {
	case *statement.PrintStmt:
		return interp.executePrintStmt(v)
	case *statement.ExpressionStmt:
		return interp.executeExprStmt(v)
	case *statement.FunctionStmt:
		return interp.executeFuncStmt(v)
	case *statement.VarStmt:
		return interp.executeVarStmt(v)
	case *statement.BlockStmt:
		return interp.executeBlockStmt(v)
	case *statement.IfStmt:
		return interp.executeIfStmt(v)
	case *statement.WhileStmt:
		return interp.executeWhileStmt(v)
	case *statement.ReturnStmt:
		return interp.executeReturnStmt(v)
	default:
		panic(fmt.Sprintf("unimplemented: %#v", stmt))
	}
}

func (interp *Interpreter) executeVarStmt(stmt *statement.VarStmt) (interface{}, error) {
	// this part forbids shadowing variable names
	if _, defined := interp.env.Lookup(stmt.Name); defined {
		err := interp.reporter.Report(fmt.Sprintf("variable named '%s' already exists", stmt.Name.Lexeme), "TODO", "", stmt.Name.Line, 0, 0)
		return nil, err
	}

	var (
		val interface{}
		err error
	)
	if stmt.Initializer != nil {
		if val, err = interp.evaluate(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	interp.env.Define(stmt.Name.Lexeme, val)
	return nil, nil
}

func (interp *Interpreter) executePrintStmt(stmt *statement.PrintStmt) (interface{}, error) {
	val, err := interp.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", val)
	return nil, nil
}

func (interp *Interpreter) executeExprStmt(stmt *statement.ExpressionStmt) (interface{}, error) {
	val, err := interp.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	if interp.repl {
		fmt.Println("eval:", val)
	}
	return nil, nil
}

func (interp *Interpreter) executeFuncStmt(stmt *statement.FunctionStmt) (interface{}, error) {
	function := NewLoxFunction(stmt, interp.env)
	interp.env.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (interp *Interpreter) executeBlockStmt(stmt *statement.BlockStmt) (interface{}, error) {
	env := environment.NewEnvironment(interp.env)
	return interp.executeBlock(stmt.Statements, env)
}

// TODO: Implement ReturnValue struct which holds an array of return values.
// We should also support multiple return values.
// Implement NewLine token, so that we can ommit ';' from most of the code base.
func (interp *Interpreter) executeBlock(statements []statement.Stmt, env *environment.Environment) (interface{}, error) {
	previous := interp.env
	defer interp.setEnvironment(previous)
	interp.env = env

	for _, s := range statements {
		retval, err := interp.execute(s)
		if err != nil {
			return nil, err
		}
		// TODO: Potential rework. Make return statements more elegant.
		if retval != nil {
			return retval, nil
		}
	}

	return nil, nil
}

func (interp *Interpreter) executeIfStmt(stmt *statement.IfStmt) (interface{}, error) {
	cond, err := interp.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		return interp.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return interp.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (interp *Interpreter) executeWhileStmt(stmt *statement.WhileStmt) (interface{}, error) {
	for {
		cond, err := interp.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if isTruthy(cond) {
			if _, err := interp.execute(stmt.Body); err != nil {
				return nil, err
			}
			continue
		}
		return nil, nil
	}
}

func (interp *Interpreter) executeReturnStmt(stmt *statement.ReturnStmt) (interface{}, error) {
	var (
		val interface{} = nil
		err error
	)

	if stmt.Value != nil {
		if val, err = interp.evaluate(stmt.Value); err != nil {
			return nil, err
		}
	}
	return val, nil
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
	case *expression.Logical:
		return interp.evaluateLogicalExpr(v)
	case *expression.Call:
		return interp.evaluateCallExpr(v)
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
	case token.Bang, token.Not:
		return !isTruthy(right), nil
	case token.Minus:
		if err := interp.checkNumberOperand(expr.Operator, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
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
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GreaterEqual:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.Less:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LessEqual:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.EqualEqual:
		if err := interp.checkEqualityOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return isEqual(left, right), nil
	case token.BangEqual:
		return !isEqual(left, right), nil
	case token.Minus:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.Plus:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) + right.(float64), nil
	case token.Slash:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.Star:
		if err := interp.checkNumberOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.DotDot:
		if err := interp.checkStringOperands(expr.Operator, left, right); err != nil {
			interp.reporter.Report(err.Error(), "TODO", "", expr.Operator.Line, 0, 0)
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
		return nil, interp.reporter.Report(fmt.Sprintf("undefined variable '%s'", expr.Name.Lexeme), "", "", expr.Name.Line, 0, 0)
	}

	return val, nil
}

func (interp *Interpreter) evaluateAssignExpr(expr *expression.Assign) (interface{}, error) {
	val, err := interp.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if ok := interp.env.Assign(expr.Name, val); !ok {
		return nil, interp.reporter.Report(fmt.Sprintf("undefined variable '%s'", expr.Name.Lexeme), "", "", expr.Name.Line, 0, 0)
	}
	return val, nil
}

func (interp *Interpreter) evaluateLogicalExpr(expr *expression.Logical) (interface{}, error) {
	leftVal, err := interp.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.Or {
		if isTruthy(leftVal) {
			return leftVal, nil
		}
	} else {
		if !isTruthy(leftVal) {
			return leftVal, nil
		}
	}

	rightVal, err := interp.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	return rightVal, nil
}

func (interp *Interpreter) evaluateCallExpr(expr *expression.Call) (interface{}, error) {
	callee, err := interp.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, 0)
	for _, arg := range expr.Args {
		val, err := interp.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		err = interp.reporter.Report(fmt.Sprintf("'%v' is not a callable function or a class", callee), "", "", expr.Paren.Line, 0, 0)
		return nil, err
	}

	if len(args) != function.Arity() {
		err = interp.reporter.Report(fmt.Sprintf("expect %d arguments but got %d", function.Arity(), len(args)), "", "", expr.Paren.Line, 0, 0)
		return nil, err
	}

	return function.Call(interp, args)

}

func (interp *Interpreter) setEnvironment(env *environment.Environment) {
	interp.env = env
}

func (interp *Interpreter) checkNumberOperand(operator *token.Token, operand interface{}) error {
	_, ok := operand.(float64)
	if !ok {
		return interp.reporter.Report(fmt.Sprintf("operand for binary operator '%s' must be a number", operator.Lexeme), "", "", operator.Line, 0, 0)
	}
	return nil
}

func (interp *Interpreter) checkNumberOperands(operator *token.Token, left, right interface{}) error {
	_, lok := left.(float64)
	if !lok {
		return interp.reporter.Report(fmt.Sprintf("left operand for binary operator '%s' must be a number", operator.Lexeme), "", "", operator.Line, 0, 0)
	}
	rval, rok := right.(float64)
	if !rok {
		return interp.reporter.Report(fmt.Sprintf("right operand for binary operator '%s' must be a number", operator.Lexeme), "", "", operator.Line, 0, 0)
	}

	// special case
	if operator.Type == token.Slash && rval == 0.0 {
		return interp.reporter.Report("division by zero", "", "", operator.Line, 0, 0)

	}
	return nil
}

func (interp *Interpreter) checkEqualityOperands(operator *token.Token, left, right interface{}) error {
	_, lfloat := left.(float64)
	_, rfloat := right.(float64)

	_, lstr := left.(string)
	_, rstr := right.(string)

	_, lbool := left.(bool)
	_, rbool := right.(bool)

	if lfloat != rfloat || lstr != rstr || lbool != rbool {
		return interp.reporter.Report(fmt.Sprintf("left and right operands for binary operator '%s' must be of same type", operator.Lexeme), "", "", operator.Line, 0, 0)
	}

	return nil
}

func (interp *Interpreter) checkStringOperands(operator *token.Token, left, right interface{}) error {
	_, lok := left.(string)
	if !lok {
		return interp.reporter.Report(fmt.Sprintf("left operand for binary operator '%s' must be a string", operator.Lexeme), "", "", operator.Line, 0, 0)
	}
	_, rok := right.(string)
	if !rok {
		return interp.reporter.Report(fmt.Sprintf("right operand for binary operator '%s' must be a string", operator.Lexeme), "", "", operator.Line, 0, 0)
	}
	return nil
}

func stringify(val interface{}) string {
	if val == nil {
		return "null"
	}
	return fmt.Sprintf("%#v", val)
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	if stringVal, ok := val.(string); ok {
		return stringVal != ""
	}
	if floatVal, ok := val.(float64); ok {
		return floatVal != 0.0
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

	return compareBools(left, right) || compareFloats(left, right) || compareStrings(left, right)
}

func compareBools(left, right interface{}) bool {
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

func compareFloats(left, right interface{}) bool {
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

func compareStrings(left, right interface{}) bool {
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
