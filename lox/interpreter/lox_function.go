package interpreter

import (
	"golox/lox/environment"
	"golox/lox/expression"
	"golox/lox/statement"
)

type LoxFunction struct {
	declaration *statement.FunctionStmt
	closure     *environment.Environment
}

func NewLoxFunction(decl *statement.FunctionStmt, closure *environment.Environment) *LoxFunction {
	return &LoxFunction{
		declaration: decl,
		closure:     closure,
	}
}

func (f *LoxFunction) Call(interp *Interpreter, args []interface{}) (interface{}, error) {
	env := environment.NewEnvironment(f.closure)
	// Interpreter.evaluateCallExpr() already checks if the number of arguments match
	for i := 0; i < len(f.declaration.Params); i++ {
		env.Define(f.declaration.Params[i].Lexeme, args[i])
	}
	retval, err := interp.executeBlock(f.declaration.Body, env)
	if err != nil {
		return nil, err
	}
	if _, isNull := retval.(expression.NullExpr); isNull {
		return nil, nil
	}
	return retval, nil
}
func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}
func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
