package reporter

import "fmt"

type ErrorReporter struct {
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{}
}

type Location struct {
	file, source                       string
	lineNumber, startColumn, endColumn int
}

func NewLocation(file, source string, lineNr, startCol, endCol int) *Location {
	return &Location{
		file:        file,
		source:      source,
		lineNumber:  lineNr,
		startColumn: startCol,
		endColumn:   endCol,
	}
}

func (r *ErrorReporter) Report(loc *Location, err error) {
	fmt.Printf("error: %s\nfile %s, line %d:\n\n\t%s\n", err, loc.file, loc.lineNumber, loc.source)
}
