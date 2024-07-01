package main

import (
	"os"
	"path/filepath"
)

func getKadaiNum() string {
	curDir, _ := os.Getwd()
	return filepath.Base(curDir)
}
