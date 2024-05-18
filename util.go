package main

import (
	"errors"
	"os"
	"runtime"
)

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		return homeDrive + homePath
	}
	return os.Getenv("HOME")
}
