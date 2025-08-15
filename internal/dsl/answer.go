package dsl

import (
	"fmt"
	"log/slog"
	"strings"
)

// Answer is a struct for storing answer data
type Answer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

// answCheck checks if answers are written using correct DSL syntax
func answCheck(qlines []string, offset int, errInfo *[]Incorrect) {
	answRange := GetBlocks(&qlines, newSection, 0)[0]
	answers := qlines[answRange.start : answRange.end+1]

	if len(answers) == 0 {
		logErr(offset+answRange.start, "No answers has been found", errInfo)
		return
	}

	var identifier string
	if strings.Index(answers[0], oneOpt) == 0 {
		identifier = oneOpt
		slog.Info(fmt.Sprintf(`Found identifier: "%s" for 1st option "%s"`, identifier, oneOpt))
	} else if strings.Index(answers[0], multipleOpt) == 0 {
		identifier = multipleOpt
		slog.Info(fmt.Sprintf(`Found identifier: "%s" for 1st option "%s"`, identifier, multipleOpt))
	} else {
		logErr(offset+answRange.start, "Answer identifier is not supported", errInfo)
		return
	}

	for i, answer := range answers[1:] {
		slog.Info("Checking", "option", i+2, "answer", answer)
		if strings.Index(answer, identifier) != 0 {
			logErr(offset+answRange.start+i,
				fmt.Sprintf("Inconsistent identifier for answer option %d", i), errInfo)

		}
	}

	var correctAnsw int
	for _, answer := range answers {
		if strings.Index(answer, correctSym) == len(answer)-1 {
			correctAnsw++
		}
	}

	if correctAnsw == 0 {
		logErr(offset+answRange.start, "No correct answers provided", errInfo)
	} else if correctAnsw == len(answers) {
		logErr(offset+answRange.start, "All answer options are correct", errInfo)
	} else if correctAnsw != 1 && string(identifier) == oneOpt {
		logErr(offset+answRange.start,
			fmt.Sprintf(`Test options with "%s" mark can have only one correct answer, found %d"`,
				oneOpt, correctAnsw), errInfo)
	} else {
		slog.Info("Check mark for correct answer... passed")
	}
}
