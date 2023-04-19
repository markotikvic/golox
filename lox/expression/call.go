package expression

import "golox/lox/token"

type Call struct {
	Calee Expression
	Paren *token.Token
	Args  []Expression
}

func NewCall(calee Expression, paren *token.Token, args []Expression) *Call {
	return &Call{
		Calee: calee,
		Paren: paren,
		Args:  args,
	}
}
