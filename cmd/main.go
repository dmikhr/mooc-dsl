package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dmikhr/mooc-dsl/internal/config"
	"github.com/dmikhr/mooc-dsl/internal/dsl"
	"github.com/dmikhr/mooc-dsl/internal/storage"
)

var fpath string
var showErrors, verbose bool

func init() {
	flag.StringVar(&fpath, "fname", "", "path to file with test")
	flag.BoolVar(&showErrors, "showErrors", true, "print found errors")
	flag.BoolVar(&verbose, "verbose", true, "print results")
	flag.Parse()
}

func main() {
	if fpath == "" {
		fpath = config.DefaultSourceFile
	}

	log.Println("Opening:", fpath)
	lines := storage.LoadDSL(fpath)

	errInfo := dsl.SyntaxCheck(&lines)
	errJSON, err := json.MarshalIndent(dsl.ErrWrap{Incorrect: errInfo}, "", "\t")
	if err != nil {
		log.Fatal("errInfo to JSON", err)
	}

	if len(errInfo) == 0 {
		slog.Info("No errors found")
		q := dsl.Parse(&lines)
		var qJSON []byte
		qJSON, err = json.MarshalIndent(dsl.QuizWrap{Quiz: q}, "", "\t")
		if err != nil {
			log.Fatal("Quiz to JSON", err)
		}

		if verbose {
			fmt.Println(string(qJSON))
		}
		err = storage.SaveJSON(qJSON, fpath, false)
		if err != nil {
			log.Fatal("Err while saving file", err)
		}
	} else {
		if showErrors {
			fmt.Fprintln(os.Stderr, string(errJSON))
		}
		err = storage.SaveJSON(errJSON, fpath, true)
		if err != nil {
			log.Fatal("Err while saving file", err)
		}
	}
}
