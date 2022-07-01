package main

import (
	"fmt"
	"strings"
)

type StatementType string

const (
	STATEMENT_COMMAND     StatementType = "CMD"
	STATEMENT_FILE        StatementType = "FILE"
	STATEMENT_DIR         StatementType = "DIR"
	STATEMENT_FEATURE     StatementType = "FEAT"
	STATEMENT_FEATURE_END StatementType = "ENDFEAT"
	STATEMENT_RMFILE      StatementType = "RMFILE"
	STATEMENT_RMDIR       StatementType = "RMDIR"
)

type Statement struct {
	Type StatementType
	Args []string
}

type Template struct {
	Features map[string][]Statement
}

func ParseTemplateFile(path string) (result Template, success bool) {
	file, success := ReadFile(path)
	if !success {
		return Template{}, false
	}

	result.Features = make(map[string][]Statement)

	activeFeature := "default"
	result.Features[activeFeature] = []Statement{}

	lines := strings.Split(strings.ReplaceAll(file, "\r\n", "\n"), "\n")

	for index, line := range lines {
		if len(line) == 0 {
			continue
		}

		// Comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		t, args := splitStatement(strings.TrimSpace(line))

		if t == string(STATEMENT_COMMAND) {
			result.Features[activeFeature] = append(result.Features[activeFeature], Statement{Type: STATEMENT_COMMAND, Args: []string{args}})
		} else if t == string(STATEMENT_DIR) {
			result.Features[activeFeature] = append(result.Features[activeFeature], Statement{Type: STATEMENT_DIR, Args: []string{args}})
		} else if t == string(STATEMENT_FILE) {
			result.Features[activeFeature] = append(result.Features[activeFeature], Statement{Type: STATEMENT_FILE, Args: strings.Split(args, " ")})
		} else if t == string(STATEMENT_RMFILE) {
			result.Features[activeFeature] = append(result.Features[activeFeature], Statement{Type: STATEMENT_RMFILE, Args: []string{args}})
		} else if t == string(STATEMENT_RMDIR) {
			result.Features[activeFeature] = append(result.Features[activeFeature], Statement{Type: STATEMENT_RMDIR, Args: []string{args}})
		} else if t == string(STATEMENT_FEATURE) {
			activeFeature = args
			result.Features[activeFeature] = []Statement{}
		} else if t == string(STATEMENT_FEATURE_END) {
			activeFeature = "default"
		} else {
			fmt.Printf("Error: unknown command %s at line %d\n", t, index)
			return Template{}, false
		}
	}

	if activeFeature != "default" {
		fmt.Println("Missing ENDFEAT")
	}

	return
}

func splitStatement(line string) (commandType string, args string) {
	commandType, args, _ = strings.Cut(line, " ")

	return
}
