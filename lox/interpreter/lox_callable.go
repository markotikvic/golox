package interpreter

type LoxCallable interface {
	Call(interp *Interpreter, args []interface{}) (interface{}, error)
	Arity() int
	String() string
}

type LoxCallableImpl struct {
	arity     int
	call      func(interp *Interpreter, args []interface{}) (interface{}, error)
	stringify func() string
}

func NewLoxCallable(arity int, callFunc func(interp *Interpreter, args []interface{}) (interface{}, error), strFunc func() string) LoxCallable {
	return &LoxCallableImpl{
		arity:     arity,
		call:      callFunc,
		stringify: strFunc,
	}
}

func (c *LoxCallableImpl) Call(interp *Interpreter, args []interface{}) (interface{}, error) {
	return c.call(interp, args)
}

func (c *LoxCallableImpl) Arity() int {
	return c.arity
}

func (c *LoxCallableImpl) String() string {
	return c.stringify()
}
