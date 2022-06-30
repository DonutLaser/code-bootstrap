package main

import (
	"fmt"
	"log"
	"os"
)

type Args struct {
	Lang     string
	Name     string
	Features []string
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

func parseArgs() (result Args) {
	args := os.Args[1:]

	if len(args) < 2 {
		log.Fatal("Usage: <lang> <project-name> [feature-name-1] [feature-name-2] ...")
	}

	result = Args{
		Lang: args[0],
		Name: args[1],
	}

	if len(args) > 2 {
		result.Features = args[2:]
	}

	return
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

	args := parseArgs()

	if IsInArray(langs, args.Lang) {
		template, success := ParseTemplateFile(fmt.Sprintf("./templates/%s/template", args.Lang))
		if !success {
			return
		}

		success = GenerateFromTemplate(template, args.Lang, args.Name, args.Features, settings)
		if !success {
			return
		}
	} else {
		fmt.Println("Unsupported lang")
	}
}
