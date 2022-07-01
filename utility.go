package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getExecutableDir() string {
	exePath, _ := os.Executable()
	return path.Dir(strings.ReplaceAll(exePath, "\\", "/"))
}
