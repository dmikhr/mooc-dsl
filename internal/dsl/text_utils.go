package dsl

import (
	"log/slog"
	"strings"
)

// trimMarker returns text without marker
// marker is part of DSL syntax and must be removed when question text is extracted
func trimMarker(text string) string {
	return strings.TrimSpace(text[strings.Index(text, " ")+1:])
}

// isMultiple returns true if text contains marker for multiple choice option
func isMultiple(text string) bool {
	marker := text[:strings.Index(text, " ")]
	switch marker {
	case oneOpt:
		return false
	case multipleOpt:
		return true
	}
	return false
}

// logErr logs error and adds it to errInfo
func logErr(lnum int, description string, errInfo *[]Incorrect) {
	errItem := Incorrect{LineNumber: lnum, ErrDescription: description}
	slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
	*errInfo = append(*errInfo, errItem)
}
