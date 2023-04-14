package expression

import "golox/lox/token"

type Binary struct {
	Left, Right Expression
	Operator    *token.Token
}

func NewBinary(left Expression, operator *token.Token, right Expression) *Binary {
	return &Binary{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (e *Binary) Expression() {}
