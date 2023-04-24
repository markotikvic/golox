package resolver

import (
	"fmt"
	"golox/lox/expression"
	"golox/lox/interpreter"
	"golox/lox/reporter"
	"golox/lox/statement"
	"golox/lox/token"
)

type Resolver struct {
	interp   *interpreter.Interpreter
	reporter *reporter.ErrorReporter
	scopes   []map[string]bool
}

func New(interp *interpreter.Interpreter, reporter *reporter.ErrorReporter) *Resolver {
	return &Resolver{
		interp:   interp,
		reporter: reporter,
		scopes:   make([]map[string]bool, 0),
	}
}

func (r *Resolver) resolveBlockStmt(stmt *statement.BlockStmt) (interface{}, error) {
	r.beginScope()
	for _, s := range stmt.Statements {
		if err := r.resolveStmt(s); err != nil {
			return nil, err
		}
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt statement.Stmt) error {
	switch v := stmt.(type) {
	case *statement.VarStmt:
		return r.resolveVarStmt(v)
	}
	return nil
}

func (r *Resolver) resolveVarStmt(stmt *statement.VarStmt) error {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return err
		}
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) declare(name *token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	inneermost := r.peekScope()
	inneermost[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	inneermost := r.peekScope()
	inneermost[name.Lexeme] = true
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
	}
	return nil
}

func (r *Resolver) resolveVarExpr(expr *expression.Variable) error {
	if len(r.scopes) > 0 && !r.peekScope()[expr.Name.Lexeme] {
		msg := fmt.Sprintf("can't read local variable '%s' in its own initializer", expr.Name.Lexeme)
		return r.reporter.Report(msg, "", "", expr.Name.Line, 0, 0)
	}
	return r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) resolveLocal(expr *expression.Variable, tok *token.Token) error {
	panic("TODO")
	return nil
}

func (r *Resolver) beginScope() {
	scope := make(map[string]bool)
	r.scopes = append(r.scopes, scope)
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}
