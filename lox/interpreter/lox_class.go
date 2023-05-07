package interpreter

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{
		name: name,
	}
}

func (lc *LoxClass) Call(interp *Interpreter, args []interface{}) (interface{}, error) {
	instance := NewLoxInstance(lc)
	return instance, nil
}

func (lx *LoxClass) Arity() int {
	return 0
}

func (lc *LoxClass) String() string {
	return lc.name
}
