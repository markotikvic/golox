package expression

type Grouping struct {
	Expr Expression
}

func NewGrouping(expr Expression) *Grouping {
	return &Grouping{
		Expr: expr,
	}
}

func (e *Grouping) Expression() {}
