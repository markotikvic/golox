package expression

import "golox/lox/token"

type Logical struct {
	Left     Expression
	Operator *token.Token
	Right    Expression
}

func NewLogical(left Expression, operator *token.Token, right Expression) *Logical {
	return &Logical{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (e *Logical) Expression() {}
