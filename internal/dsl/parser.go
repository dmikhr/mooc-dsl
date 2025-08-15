package dsl

import (
	"errors"
	"log/slog"
	"strings"
)

// Parse parses quiz from lines
func Parse(lines *[]string) Quiz {
	var q Quiz
	q.Name = (*lines)[0]
	slog.Info("Parsing:", "quiz name", q.Name)

	q.Description = func() string {
		var descrLines []string
		for i := 2; i < len(*lines); i++ {
			if (*lines)[i] == qsep {
				break
			}
			descrLines = append(descrLines, (*lines)[i])
		}
		return strings.Join(descrLines, "\n")
	}()

	slog.Info("Parsing:", "quiz description", q.Description)
	qB := GetBlocks(lines, qsep, 2)

	for i := 0; i < len(qB); i++ {
		slog.Info("Parsing:", "question", i+1)
		question := func() Question {
			var qdescr, opts, qrec []string
			var qst Question
			var sectionN = 1 // 1 - question description, 2 - answer options, 3 - recommendations (optional)
			for j := qB[i].start; j <= qB[i].end; j++ {
				if (*lines)[j] == newSection {
					sectionN++
					slog.Info("sectionN", "", sectionN, "j", j, "qB[i].end", qB[i].end)
					continue
				}

				switch sectionN {
				case 1:
					qdescr = append(qdescr, (*lines)[j])
				case 2:
					opts = append(opts, (*lines)[j])
				case 3:
					qrec = append(qrec, (*lines)[j])
				}
			}

			qst.Text = strings.Join(qdescr, "\n")
			var answ Answer
			qst.Multiple = isMultiple(opts[0])

			for _, opt := range opts {
				if string(opt[len(opt)-1]) == correctSym {
					answ = Answer{
						Text:      trimMarker(opt[:len(opt)-1]),
						IsCorrect: true}
				} else {
					answ = Answer{
						Text:      trimMarker(opt),
						IsCorrect: false}
				}
				qst.Options = append(qst.Options, answ)
			}

			qst.Recommendation = strings.Join(qrec, "\n")
			return qst
		}()

		q.Questions = append(q.Questions, question)
	}

	return q
}

// GetItem returns element from slice by index
func GetItem[T any](index int, elements *[]T) (T, error) {
	var empty T
	if index < len(*elements) {
		return (*elements)[index], nil
	}
	return empty, errors.New("element doesn't exist")
}

// GetBlocks returns blocks of lines corresponding to one question in a quiz
// where each question is separated by sep string
func GetBlocks(lines *[]string, sep string, startIdx int) []Block {
	var b []Block
	for n := startIdx; n < len(*lines); n++ {
		if (*lines)[n] == sep {
			if len(b) > 0 {
				b[len(b)-1].end = n - 1
			}
			b = append(b, Block{start: n + 1})
		}
	}
	b[len(b)-1].end = len(*lines) - 1
	return b
}
