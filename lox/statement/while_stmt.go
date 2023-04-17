package statement

import "golox/lox/expression"

type WhileStmt struct {
	Condition expression.Expression
	Body      Stmt
}

func NewWhileStmt(condition expression.Expression, body Stmt) *WhileStmt {
	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (ws *WhileStmt) Stmt() {}
