package statement

import "golox/lox/expression"

type IfStmt struct {
	Condition    expression.Expression
	ThenBranch   Stmt
	ElifBranches []Stmt
	ElseBranch   Stmt
}

func NewIfStmt(condition expression.Expression, thenBranch Stmt, elifBranches []Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{
		Condition:    condition,
		ThenBranch:   thenBranch,
		ElifBranches: elifBranches,
		ElseBranch:   elseBranch,
	}
}

func (ps *IfStmt) Stmt() {}
