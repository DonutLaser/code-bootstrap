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

		key, value := getKeyValue(strings.TrimSpace(line))

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

func getKeyValue(str string) (key string, value string) {
	split := strings.Split(str, " ")
	return split[0], split[1]
}
