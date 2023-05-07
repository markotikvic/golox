package resolver

import (
	"fmt"
	"golox/lox/expression"
	"golox/lox/interpreter"
	"golox/lox/reporter"
	"golox/lox/statement"
	"golox/lox/token"
)

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunc
)

type Resolver struct {
	interp      *interpreter.Interpreter
	reporter    *reporter.ErrorReporter
	scopes      []map[string]bool
	currentFunc FunctionType
}

func New(interp *interpreter.Interpreter, reporter *reporter.ErrorReporter) *Resolver {
	return &Resolver{
		interp:      interp,
		reporter:    reporter,
		scopes:      make([]map[string]bool, 0),
		currentFunc: FunctionTypeNone,
	}
}

func (r *Resolver) Resolve(statements []statement.Stmt) error {
	_, err := r.resolveStmts(statements)
	return err
}

func (r *Resolver) resolveBlockStmt(stmt *statement.BlockStmt) (interface{}, error) {
	r.beginScope()
	if _, err := r.resolveStmts(stmt.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) resolveStmts(statements []statement.Stmt) (interface{}, error) {
	for _, s := range statements {
		if err := r.resolveStmt(s); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt statement.Stmt) error {
	switch v := stmt.(type) {
	case *statement.VarStmt:
		return r.resolveVarStmt(v)
	case *statement.FunctionStmt:
		return r.resolveFunctionStmt(v)
	case *statement.ExpressionStmt:
		return r.resolveExpressionStmt(v)
	case *statement.IfStmt:
		return r.resolveIfStmt(v)
	case *statement.PrintStmt:
		return r.resolvePrintStmt(v)
	case *statement.ReturnStmt:
		return r.resolveReturnStmt(v)
	case *statement.WhileStmt:
		return r.resolveWhileStmt(v)
	case *statement.ClassStmt:
		return r.resolveClassStmt(v)
	}
	return nil
}

func (r *Resolver) resolveVarStmt(stmt *statement.VarStmt) error {
	if err := r.declare(stmt.Name); err != nil {
		return err
	}
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return err
		}
	}
	return r.define(stmt.Name)
}

func (r *Resolver) resolveFunctionStmt(stmt *statement.FunctionStmt) error {
	if err := r.declare(stmt.Name); err != nil {
		return err
	}
	if err := r.define(stmt.Name); err != nil {
		return err
	}
	return r.resolveFunction(stmt, FunctionTypeFunc)
}

func (r *Resolver) resolveFunction(function *statement.FunctionStmt, funcType FunctionType) error {
	enclosingFunc := r.currentFunc
	r.currentFunc = funcType

	r.beginScope()
	for _, param := range function.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		if err := r.define(param); err != nil {
			return err
		}
	}
	if _, err := r.resolveStmts(function.Body); err != nil {
		return err
	}
	r.endScope()

	r.currentFunc = enclosingFunc

	return nil
}

func (r *Resolver) resolveExpressionStmt(stmt *statement.ExpressionStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) resolveIfStmt(stmt *statement.IfStmt) error {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return err
	}
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolvePrintStmt(stmt *statement.PrintStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) resolveReturnStmt(stmt *statement.ReturnStmt) error {
	if r.currentFunc == FunctionTypeNone {
		return r.reporter.Report("can't return from top-level code (outside of function)", stmt.Keyword)
	}
	if stmt.Value != nil {
		return r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) resolveWhileStmt(stmt *statement.WhileStmt) error {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return err
	}
	if err := r.resolveStmt(stmt.Body); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) resolveClassStmt(stmt *statement.ClassStmt) error {
	if err := r.declare(stmt.Name); err != nil {
		return err
	}
	return r.define(stmt.Name)
}

func (r *Resolver) declare(name *token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	inneermost := r.peekScope()
	if _, found := inneermost[name.Lexeme]; found {
		return r.reporter.Report(fmt.Sprintf("variable with name '%s' already exists", name.Lexeme), name)
	}
	inneermost[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	innermost := r.peekScope()
	innermost[name.Lexeme] = true
	return nil
}

func (r *Resolver) peekScope() map[string]bool {
	if len(r.scopes) == 0 {
		panic("trying to peek a scope that doesn't exist")
	}
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) resolveExpr(expr expression.Expression) error {
	switch v := expr.(type) {
	case *expression.Variable:
		return r.resolveVarExpr(v)
	case *expression.Assign:
		return r.resolveAssignExpr(v)
	case *expression.Binary:
		return r.resolveBinaryExpr(v)
	case *expression.Call:
		return r.resolveCallExpr(v)
	case *expression.Grouping:
		return r.resolveGroupingExpr(v)
	case *expression.Literal:
		return r.resolveLitteralExpr(v)
	case *expression.Logical:
		return r.resolveLogicalExpr(v)
	case *expression.Unary:
		return r.resolveUnaryExpr(v)
	}
	return nil
}

func (r *Resolver) resolveVarExpr(expr *expression.Variable) error {
	if len(r.scopes) > 0 && !r.peekScope()[expr.Name.Lexeme] {
		return r.reporter.Report(fmt.Sprintf("can't read local variable '%s' in its own initializer", expr.Name.Lexeme), expr.Name)
	}
	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) resolveAssignExpr(expr *expression.Assign) error {
	if err := r.resolveExpr(expr.Value); err != nil {
		return err
	}
	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) resolveLocal(expr expression.Expression, name *token.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interp.Resolve(expr, len(r.scopes)-1-i)
			return nil
		}
	}
	return nil
}

func (r *Resolver) resolveBinaryExpr(expr *expression.Binary) error {
	if err := r.resolveExpr(expr.Left); err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) resolveCallExpr(expr *expression.Call) error {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Args {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) resolveGroupingExpr(expr *expression.Grouping) error {
	return r.resolveExpr(expr.Expr)
}

func (r *Resolver) resolveLitteralExpr(expr *expression.Literal) error {
	return nil
}

func (r *Resolver) resolveLogicalExpr(expr *expression.Logical) error {
	if err := r.resolveExpr(expr.Left); err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) resolveUnaryExpr(expr *expression.Unary) error {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) beginScope() {
	scope := make(map[string]bool)
	r.scopes = append(r.scopes, scope)
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}
