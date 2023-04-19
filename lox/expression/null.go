package expression

type NullExpr struct{}

func MakeNull() NullExpr {
	return NullExpr{}
}

func (n NullExpr) Expression() {}
