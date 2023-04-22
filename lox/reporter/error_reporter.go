package reporter

import "fmt"

type ErrorReporter struct {
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

func (r *ErrorReporter) Report(msg, file, sourceLine string, lineNumber, startColumn, endColumn int) error {
	err := fmt.Errorf("error: %s\nfile %s, line %d:\n\n\t%s", msg, file, lineNumber, sourceLine)
	fmt.Println(err)
	return err
}
