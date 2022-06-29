package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateFromTemplate(template Template, lang string, projectName string, features []string) bool {
	success := createDirectory(projectName)
	if !success {
		return false
	}

	success = runStatements(template.Features["default"], projectName, lang)
	if !success {
		return false
	}

	for _, feature := range features {
		if statements, exists := template.Features[feature]; exists {
			success = runStatements(statements, projectName, lang)
			if !success {
				return false
			}
		} else {
			fmt.Printf("Warning: unknown feature %s\n", feature)
		}
	}

	return true
}

func runStatements(statements []Statement, projectName string, lang string) bool {
	for _, statement := range statements {
		if statement.Type == STATEMENT_COMMAND {
			commandName, commandArgs := parseCommandStatement(statement.Args[0])
			success := runCommand(commandName, projectName, commandArgs...)
			if !success {
				return false
			}

		} else if statement.Type == STATEMENT_FILE {
			fileName, templateName := statement.Args[0], statement.Args[1]
			success := createFile(fmt.Sprintf("%s/%s", projectName, fileName), fmt.Sprintf("./templates/%s/%s", lang, templateName))
			if !success {
				return false
			}
		} else if statement.Type == STATEMENT_DIR {
			dirName := statement.Args[0]
			success := createDirectory(fmt.Sprintf("%s/%s", projectName, dirName))
			if !success {
				return false
			}
		}
	}

	return true
}

func parseCommandStatement(str string) (name string, args []string) {
	split := strings.Split(str, " ")
	name = split[0]
	args = split[1:]

	return
}

func runCommand(name string, cwd string, args ...string) bool {
	fmt.Printf("Running commadn %s %s...\n", name, strings.Join(args, " "))

	var cmd = exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if cwd != "" {
		cmd.Dir = cwd
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}

	return true
}

func createFile(name string, templateName string) bool {
	contents := ""

	if templateName != "" {
		file, success := ReadFile(templateName)
		if !success {
			return false
		}

		contents = file
	}

	success := WriteFile(name, contents)
	return success
}

func createDirectory(name string) bool {
	return CreateDir(name)
}
