package main

import (
	"fmt"
	"os"
)

type Args struct {
	Command     string
	CommandArgs []string
}

func runGenerateCommand(args []string, langs []string, settings Settings) bool {
	if len(args) < 2 {
		printUsage()
		return false
	}

	lang := args[0]
	projectName := args[1]
	var features []string = make([]string, 0)
	if len(args) > 2 {
		features = args[2:]
	}

	if IsInArray(langs, lang) {
		template, success := ParseTemplateFile(fmt.Sprintf("./templates/%s/template", lang))
		if !success {
			return false
		}

		success = GenerateFromTemplate(template, lang, projectName, features, settings)
		if !success {
			return false
		}
	} else {
		fmt.Printf("Unsupported lang %s\n", lang)
		return false
	}

	return true
}

func printUsage() {
	fmt.Println("Usage: <command> [arg1] [arg2] ...")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("	generate <lang> <project-name> [feature1] [feature2] ...")
	fmt.Println("	config")
	fmt.Println("	help <command>")
}

func getSupportedLanguages() (result []string, success bool) {
	dirs, success := ReadDir("./templates")
	if !success {
		return []string{}, false
	}

	for _, item := range dirs {
		if !item.IsDir() {
			continue
		}

		result = append(result, item.Name())
	}

	return result, true
}

func parseArgs() (result Args, success bool) {
	args := os.Args[1:]

	if len(args) < 1 {
		return Args{}, false
	}

	result = Args{
		Command:     args[0],
		CommandArgs: args[1:],
	}

	return result, true
}

func main() {
	settings, success := GetSettings()
	if !success {
		return
	}

	langs, success := getSupportedLanguages()
	if !success {
		return
	}

	args, success := parseArgs()
	if !success {
		printUsage()
		return
	}

	if args.Command == "generate" {
		runGenerateCommand(args.CommandArgs, langs, settings)
	} else if args.Command == "config" {
		runCommand("start", "", "settings.conf")
	} else if args.Command == "help" {
		if len(args.CommandArgs) != 1 {
			printUsage()
			return
		}

		command := args.CommandArgs[0]
		fmt.Printf("%s\n", command)
	}
}
