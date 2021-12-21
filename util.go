package main

import (
	"errors"
	"os"
)

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
