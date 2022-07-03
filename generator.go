package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateFromTemplate(template Template, lang string, projectName string, features []string, settings Settings) {
	defaultStatementsExist := len(template.Features["default"].Statements) != 0
	anyFeaturesRequested := len(features) != 0
	standaloneFeatureCount := getRequestedStandaloneFeatureCount(features, template)

	if !defaultStatementsExist {
		if !anyFeaturesRequested {
			fmt.Printf("Nothing to be done. Did you forget to specify a feature?\n")
			return
		}
		if standaloneFeatureCount == 0 {
			fmt.Printf("None of the features requested are standalone and there are no default commands to do.\n")
			return
		}
	}

	if standaloneFeatureCount > 1 {
		fmt.Println("There can only be 1 standalone feature specified.")
		return
	}

	success := createDirectory(projectName)
	if !success {
		RemoveDir(projectName)
		return
	}

	success = runStatements(template.Features["default"].Statements, projectName, lang)
	if !success {
		RemoveDir(projectName)
		return
	}

	for _, feature := range features {
		if feat, exists := template.Features[feature]; exists {
			success = runStatements(feat.Statements, projectName, lang)
			if !success {
				RemoveDir(projectName)
				return
			}
		} else {
			fmt.Printf("Warning: unknown feature %s\n", feature)
		}
	}

	if settings.CreateVSCodeWorkspace {
		WriteFile(fmt.Sprintf("%s/%s.code-workspace", settings.VSCodeWorkspaceDir, projectName), fmt.Sprintf("{ \"folders\": [{ \"path\": \"%s\" }] }", projectName))
	}
}

func runStatements(statements []Statement, projectName string, lang string) bool {
	for _, statement := range statements {
		if statement.Type == STATEMENT_COMMAND {
			commandName, commandArgs := parseCommandStatement(statement.Args[0])
			for i := 0; i < len(commandArgs); i++ {
				commandArgs[i] = replaceBuiltinVariables(commandArgs[i], projectName)
			}

			success := runCommand(commandName, projectName, commandArgs...)
			if !success {
				return false
			}
		} else if statement.Type == STATEMENT_FILE {
			fileName, templateName := statement.Args[0], statement.Args[1]
			success := createFile(fmt.Sprintf("%s/%s", projectName, fileName), fmt.Sprintf("./templates/%s/%s", lang, templateName), projectName)
			if !success {
				return false
			}
		} else if statement.Type == STATEMENT_DIR {
			dirName := statement.Args[0]
			success := createDirectory(fmt.Sprintf("%s/%s", projectName, dirName))
			if !success {
				return false
			}
		} else if statement.Type == STATEMENT_RMFILE {
			fileName := statement.Args[0]
			success := RemoveFile(fmt.Sprintf("%s/%s", projectName, fileName))
			if !success {
				return false
			}
		} else if statement.Type == STATEMENT_RMDIR {
			dirName := statement.Args[0]
			success := RemoveDir(fmt.Sprintf("%s/%s", projectName, dirName))
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
	fmt.Printf("Running command '%s %s'...\n", name, strings.Join(args, " "))

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

func createFile(name string, templateName string, projectName string) bool {
	fmt.Printf("Creating a file %q\n", name)

	contents := ""

	if templateName != "" {
		file, success := ReadFile(templateName)
		if !success {
			return false
		}

		contents = replaceBuiltinVariables(file, projectName)
	}

	success := WriteFile(name, contents)
	return success
}

func createDirectory(name string) bool {
	fmt.Printf("Creating a directory %q\n", name)
	return CreateDir(name)
}

func replaceBuiltinVariables(str string, projectName string) string {
	result := str
	if strings.Contains(str, "{{PROJECT_NAME}}") {
		result = strings.ReplaceAll(result, "{{PROJECT_NAME}}", projectName)
	}

	return result
}

func getRequestedStandaloneFeatureCount(features []string, template Template) (result int) {
	for _, feat := range features {
		templateFeat, exists := template.Features[feat]
		if !exists {
			continue
		}

		if templateFeat.IsStandalone {
			result += 1
		}
	}

	return
}
