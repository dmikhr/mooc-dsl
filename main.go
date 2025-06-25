package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

const (
	qsep        = "==="
	newSection  = "---"
	correctSym  = "*"
	oneOpt      = "o"
	multipleOpt = "[]"
)

type ErrWrap struct {
	Incorrect []Incorrect `json:"errors"`
}

type QuizWrap struct {
	Quiz Quiz `json:"quiz"`
}

type Quiz struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type Question struct {
	Text           string   `json:"text"`
	Multiple       bool     `json:"multiple"`
	Options        []Answer `json:"options"`
	Recommendation string   `json:"recommendation"`
}

type Answer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type Incorrect struct {
	LineNumber     int    `json:"lineNumber"`
	ErrDescription string `json:"errDescription"`
}

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

func trimMarker(text string) string {
	return strings.TrimSpace(text[strings.Index(text, " ")+1:])
}

func GetItem[T any](index int, elements *[]T) (T, error) {
	var empty T
	if index < len(*elements) {
		return (*elements)[index], nil
	} else {
		return empty, errors.New("element doesn't exist")
	}
}

type Block struct {
	start int
	end   int
}

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

func logErr(lnum int, description string, errInfo *[]Incorrect) {
	errItem := Incorrect{LineNumber: lnum, ErrDescription: description}
	slog.Warn(errItem.ErrDescription, "line:", errItem.LineNumber)
	*errInfo = append(*errInfo, errItem)
}

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

func Parse(lines *[]string) Quiz {
	// parse() runs after syntax check, so err check is more lenient here
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

		// can be 1 or 2 by specs
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

func saveJSON(data []byte, fname string, isErr bool) error {
	var newSuffix string

	if isErr {
		newSuffix = "_error.json"
	} else {
		newSuffix = ".json"
	}

	data = append(data, '\n')
	fname = strings.TrimSuffix(fname, ".txt") + newSuffix
	err := os.WriteFile(fname, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

var fpath string

func init() {
	flag.StringVar(&fpath, "fname", "", "path to file with test")
	flag.Parse()
}

func main() {
	if fpath == "" {
		fpath = "assets/sample.txt"
	}

	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")

	errInfo := SyntaxCheck(&lines)
	errJson, err := json.MarshalIndent(ErrWrap{Incorrect: errInfo}, "", "\t")
	if err != nil {
		log.Fatal("errInfo to JSON", err)
	}

	fmt.Println(string(errJson))
	if len(errInfo) == 0 {
		slog.Info("No errors found")
		q := Parse(&lines)
		var qJson []byte
		qJson, err = json.MarshalIndent(QuizWrap{Quiz: q}, "", "\t")
		if err != nil {
			log.Fatal("Quiz to JSON", err)
		}
		
		fmt.Println(string(qJson))
		err = saveJSON(qJson, fpath, false)
		if err != nil {
			log.Fatal("Err while saving file", err)
		}
	} else {
		err = saveJSON(errJson, fpath, true)
		if err != nil {
			log.Fatal("Err while saving file", err)
		}
	}
}
