package storage

import (
	"io"
	"log"
	"os"
	"strings"
)

func LoadDSL(fpath string) []string {
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
	return strings.Split(string(data), "\n")
}
