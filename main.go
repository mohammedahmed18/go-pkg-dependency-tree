package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	projectPath := os.Args[1]

	if projectPath == "" {
		log.Fatal(errors.New("you must specify a project path"))
	}

	projectName, err := getProjectNameFromGoMod(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("===> %s\n", projectName)
	scanPackages(projectPath, projectName)
}

func readNonEmptyFirstLine(path string) (string, error) {
	content, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return "", err
	}
	contentWithoutWhiteSpace := strings.TrimSpace(string(content))
	lines := strings.Split(contentWithoutWhiteSpace, "\n")

	return lines[0], nil
}
