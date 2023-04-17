package statement

import "golox/lox/expression"

type IfStmt struct {
	Condition  expression.Expression
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition expression.Expression, thenBranch, elseBranch Stmt) *IfStmt {
	return &IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (ps *IfStmt) Stmt() {}
