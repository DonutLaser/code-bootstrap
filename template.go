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
	STATEMENT_STANDALONE  StatementType = "STANDALONE"
	STATEMENT_FEATURE_END StatementType = "ENDFEAT"
	STATEMENT_RMFILE      StatementType = "RMFILE"
	STATEMENT_RMDIR       StatementType = "RMDIR"
)

type Statement struct {
	Type StatementType
	Args []string
}

type Feature struct {
	IsStandalone bool
	Statements   []Statement
}

type Template struct {
	Features map[string]Feature
}

func ParseTemplateFile(templatePath string) (result Template, success bool) {
	file, success := ReadFile(fmt.Sprintf("%s/%s", getExecutableDir(), templatePath))
	if !success {
		return Template{}, false
	}

	result.Features = make(map[string]Feature)

	activeFeature := "default"
	result.Features[activeFeature] = Feature{}

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
			feat := result.Features[activeFeature]
			feat.Statements = append(feat.Statements, Statement{Type: STATEMENT_COMMAND, Args: []string{args}})
			result.Features[activeFeature] = feat
		} else if t == string(STATEMENT_DIR) {
			feat := result.Features[activeFeature]
			feat.Statements = append(feat.Statements, Statement{Type: STATEMENT_DIR, Args: []string{args}})
			result.Features[activeFeature] = feat
		} else if t == string(STATEMENT_FILE) {
			feat := result.Features[activeFeature]
			feat.Statements = append(feat.Statements, Statement{Type: STATEMENT_FILE, Args: strings.Split(args, " ")})
			result.Features[activeFeature] = feat
		} else if t == string(STATEMENT_RMFILE) {
			feat := result.Features[activeFeature]
			feat.Statements = append(feat.Statements, Statement{Type: STATEMENT_RMFILE, Args: []string{args}})
			result.Features[activeFeature] = feat
		} else if t == string(STATEMENT_RMDIR) {
			feat := result.Features[activeFeature]
			feat.Statements = append(feat.Statements, Statement{Type: STATEMENT_RMDIR, Args: []string{args}})
			result.Features[activeFeature] = feat
		} else if t == string(STATEMENT_FEATURE) {
			activeFeature = args
		} else if t == string(STATEMENT_FEATURE_END) {
			activeFeature = "default"
		} else if t == string(STATEMENT_STANDALONE) {
			if activeFeature == "default" {
				fmt.Printf("Error line %d: STANDALONE can only be used inside FEAT block", index)
				return Template{}, false
			}

			feat := result.Features[activeFeature]
			feat.IsStandalone = true
			result.Features[activeFeature] = feat
		} else {
			fmt.Printf("Error line %d: unknown command %s\n", index, t)
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
