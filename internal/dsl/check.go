package dsl

import (
	"fmt"
	"log/slog"
)

// SyntaxCheck checks DSL syntax of quiz
func SyntaxCheck(lines *[]string) []Incorrect {
	var errInfo []Incorrect
	var sep1 string

	slog.Info("Syntax check: START")
	_, err := GetItem(0, lines)
	if err != nil {
		errItem := Incorrect{LineNumber: 1, ErrDescription: "No title"}
		slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
		errInfo = append(errInfo, errItem)
	} else {
		slog.Info("Syntax check: Title ok")
	}

	sep1, err = GetItem(1, lines)
	if err != nil {
		errItem := Incorrect{LineNumber: 2,
			ErrDescription: fmt.Sprintf("No separator")}
		slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
		errInfo = append(errInfo, errItem)
	} else if sep1 != newSection {
		errItem := Incorrect{LineNumber: 2,
			ErrDescription: fmt.Sprintf(`Second line must contain "%s" separator, found "%s" instead`,
				newSection, sep1)}
		slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
		errInfo = append(errInfo, errItem)
	} else {
		slog.Info("Syntax check: Separator after title ok")
	}

	qB := GetBlocks(lines, qsep, 2)
	slog.Info("Found questions", "total", len(qB))
	slog.Info(fmt.Sprintf("Question blocks (here indexes not line nums): %v", qB))

	if len(qB) == 0 {
		errItem := Incorrect{LineNumber: len(*lines) - 1,
			ErrDescription: "No questions found"}
		slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
		errInfo = append(errInfo, errItem)
	}

	for q := 0; q < len(qB); q++ {
		slog.Info("Checking question", "N", q+1)
		var secCount int
		for i := qB[q].start; i < qB[q].end; i++ {
			if (*lines)[i] == newSection {
				secCount++
			}
		}

		// can be 1 or 2 by DSL specs
		if secCount == 0 || secCount > 3 {
			slog.Warn("Number of separators error:", "line", qB[q].start+1, "current num",
				secCount, "should be", "1 or 2", "in question ", q+1)
		} else {
			slog.Info(fmt.Sprintf("Number of separators: %d in question %d ok", secCount, q+1))
			slog.Info("Checking answer options")
			answCheck((*lines)[qB[q].start:qB[q].end+1], qB[q].start, &errInfo)
		}
	}

	slog.Info("Syntax check: FINISHED")
	return errInfo
}
