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

func printFeatures(args []string, langs []string) bool {
	if len(args) != 1 {
		printUsage()
		return false
	}

	lang := args[0]

	if IsInArray(langs, lang) {
		template, success := ParseTemplateFile(fmt.Sprintf("./templates/%s/template", lang))
		if !success {
			return false
		}

		for k := range template.Features {
			if k != "default" {
				fmt.Printf("%s\n", k)
			}
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
	fmt.Println("  generate <type> <project-name> [feature1] [feature2] ...")
	fmt.Println("  features <type>")
	fmt.Println("  config")
	fmt.Println("  help <command>")
}

func printCommandHelp(command string, langs []string) {
	switch command {
	case "generate":
		fmt.Println("generate <type> <project-name> [feature1] [feature2] ...")
		fmt.Println("Generate a project with the specified name.")
		fmt.Println("Available project types:")
		for _, lang := range langs {
			fmt.Printf("  - %s\n", lang)
		}
	case "features":
		fmt.Println("features <type>")
		fmt.Println("List available features for the specified project type")
	case "config":
		fmt.Println("config")
		fmt.Println("Open the configuration file in the default text editor")
	default:
		fmt.Printf("Unknown command %q\n", command)
	}
}

func getSupportedLanguages() (result []string, success bool) {
	dirs, success := ReadDir(fmt.Sprintf("%s/templates", getExecutableDir()))
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
	} else if args.Command == "features" {
		printFeatures(args.CommandArgs, langs)
	} else if args.Command == "config" {
		OpenWithDefaultProgram("settings.conf")
	} else if args.Command == "help" {
		if len(args.CommandArgs) != 1 {
			printUsage()
			return
		}

		command := args.CommandArgs[0]
		printCommandHelp(command, langs)
	}
}
