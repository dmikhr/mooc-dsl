package storage

import (
	"os"
	"strings"

	"github.com/dmikhr/mooc-dsl/internal/config"
)

// SaveJSON saves quiz to file
func SaveJSON(data []byte, fname string, isErr bool) error {
	var suffix string

	if isErr {
		suffix = config.ErrPrefix + ".json"
	} else {
		suffix = ".json"
	}

	data = append(data, '\n')
	fname = strings.TrimSuffix(fname, ".txt") + suffix
	err := os.WriteFile(fname, data, config.WritePermission)

	if err != nil {
		return err
	}

	return nil
}
