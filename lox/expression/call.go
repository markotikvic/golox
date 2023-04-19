package expression

import "golox/lox/token"

type Call struct {
	Callee Expression
	Paren  *token.Token
	Args   []Expression
}

func NewCall(callee Expression, paren *token.Token, args []Expression) *Call {
	return &Call{
		Callee: callee,
		Paren:  paren,
		Args:   args,
	}
}

func (c *Call) Expression() {}
