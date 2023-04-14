package expression

import "golox/lox/token"

type Unary struct {
	Right    Expression
	Operator *token.Token
}

func NewUnary(operator *token.Token, right Expression) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (e *Unary) Expression() {}
