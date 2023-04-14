package reporter

import "fmt"

type ErrorReporter struct {
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

func (r *ErrorReporter) Report(file string, lineNumber int, line string, err error) {
	fmt.Printf("error: %s\nin file %s on line %d:\n\t%s\n", err, file, lineNumber, line)
}
