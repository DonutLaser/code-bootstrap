package main

import (
	"fmt"
	"strings"
)

type Settings struct {
	CreateVSCodeWorkspace bool
	VSCodeWorkspaceDir    string
}

func GetSettings() (Settings, bool) {
	file, success := ReadFile(fmt.Sprintf("%s/settings.conf", getExecutableDir()))
	if !success {
		return Settings{}, false
	}

	result := Settings{}

	lines := strings.Split(file, "\n")

	for index, line := range lines {
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		key, value := getKeyValue(strings.TrimSpace(line), " ")

		if key == "create_vscode_workspace" {
			result.CreateVSCodeWorkspace = value == "true"
		} else if key == "vscode_workspace_dir" {
			result.VSCodeWorkspaceDir = value
		} else {
			fmt.Printf("Error: unknown option %s\n at line %d", key, index)
			return Settings{}, false
		}
	}

	return result, true
}

func UpdateSettings(settings Settings, key string, value string) {
	if key == "create_vscode_workspace" {
		if !validateBool(value) {
			fmt.Printf("`%s` only supports values `true` and `false`, but received value is `%s`", key, value)
			return
		}

		settings.CreateVSCodeWorkspace = value == "true"
	} else if key == "vscode_workspace_dir" {
		settings.VSCodeWorkspaceDir = value
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("create_vscode_workspace %t\n", settings.CreateVSCodeWorkspace))
	sb.WriteString(fmt.Sprintf("vscode_workspace_dir %s\n", settings.VSCodeWorkspaceDir))

	WriteFile(fmt.Sprintf("%s/settings.conf", getExecutableDir()), sb.String())
}

func getKeyValue(str string, sep string) (key string, value string) {
	split := strings.Split(str, sep)
	return split[0], split[1]
}

func validateBool(value string) bool {
	return value == "true" || value == "false"
}
