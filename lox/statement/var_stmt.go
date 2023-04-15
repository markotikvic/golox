package statement

import (
	"golox/lox/expression"
	"golox/lox/token"
)

type VarStmt struct {
	Name        *token.Token
	Initializer expression.Expression
}

func NewVarStmt(name *token.Token, initalizer expression.Expression) *VarStmt {
	return &VarStmt{
		Name:        name,
		Initializer: initalizer,
	}
}

func (ps *VarStmt) Stmt() {}
