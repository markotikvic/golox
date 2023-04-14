package reporter

import "fmt"

type ErrorReporter struct {
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

func (r *ErrorReporter) ReportAtLocation(err error, file, source string, lineNr, startCol, endCol int) {
	fmt.Printf("error: %s\nfile %s, line %d:\n\n\t%s\n", err, file, lineNr, source)
}

func (r *ErrorReporter) Report(err error) {
	fmt.Printf("error: %s\n", err)
}
