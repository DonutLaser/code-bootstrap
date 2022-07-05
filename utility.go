package main

import (
	"os"
	"path"
	"strings"
)

func getExecutableDir() string {
	exePath, _ := os.Executable()
	return path.Dir(strings.ReplaceAll(exePath, "\\", "/"))
}
