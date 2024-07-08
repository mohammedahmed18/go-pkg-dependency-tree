package main

import "strings"

func getProjectNameFromGoMod(path string) (string, error) {
	firstLine, err := readNonEmptyFirstLine(path)
	if err != nil {
		return "", err
	}

	return strings.Split(firstLine, " ")[1], nil
}
