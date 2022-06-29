package main

import (
	"fmt"
	"io/fs"
	"os"
)

func ReadFile(path string) (result string, success bool) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", false
	}

	return string(bytes), true
}

func WriteFile(path string, content string) bool {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	defer file.Close()

	file.WriteString(content)
	file.Sync()

	return true
}

func ReadDir(path string) ([]fs.DirEntry, bool) {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return []fs.DirEntry{}, false
	}
	defer dir.Close()

	items, err := dir.ReadDir(-1)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return []fs.DirEntry{}, false
	}

	return items, true
}

func CreateDir(path string) bool {
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}

	return true
}
