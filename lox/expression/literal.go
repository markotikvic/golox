package expression

type Literal struct {
	Value interface{}
}

func NewLiteral(val interface{}) *Literal {
	return &Literal{
		Value: val,
	}
}

func (e *Literal) Expression() {}
