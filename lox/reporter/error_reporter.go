package reporter

import (
	"fmt"
	"golox/lox/token"
)

type ErrorReporter struct {
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

func (r *ErrorReporter) Report(msg string, tok *token.Token) error {
	err := fmt.Errorf("%s\n%s %d:\t%s", msg, tok.File, tok.Line, tok.Source)
	fmt.Println(err)
	return err
}
