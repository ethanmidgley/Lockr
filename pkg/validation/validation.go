package validation

import (
	"errors"
	"os"
)

func NotEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("input cannot be empty")
	}
	return nil
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
