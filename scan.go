package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type goFile struct {
	name        string
	linesNumber int64
	*goPackage
}

type goPackage struct {
	name string
	deps []*goPackage
}

func scanPackages(path string, projectName string) (map[string][]*goPackage, error) {
	goFilesPaths, err := getGoFiles(path)
	if err != nil {
		return nil, err
	}

	gofiles := make([]*goFile, 0)

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(goFilesPaths))

	for _, path := range goFilesPaths {
		go func(goFilePath string) {
			f, _ := parseGoFile(goFilePath, projectName)
			gofiles = append(gofiles, f)
			waitGroup.Done()
		}(path)
	}

	waitGroup.Wait()
	return generatePackagesGraph(gofiles), nil
}

func removeQuotesFromPkg(pkg string) string {
	pkg = strings.TrimPrefix(pkg, "\"")
	return strings.TrimSuffix(pkg, "\"")
}

func parseGoFile(path string, projectName string) (*goFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundImportLines := false

	var gFile goFile
	gFile.name = path

	depsPackages := make([]*goPackage, 0)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		lineNotEmpty := line != ""
		packageLine := strings.HasPrefix(line, "package")
		isImportProjectPackageLine := func(l string) bool { return strings.HasPrefix(l, "\""+projectName) }
		commentLine := strings.HasPrefix(line, "//")
		singleImportLine := strings.HasPrefix(line, "import") && !strings.Contains(line, "(")

		if commentLine {
			continue
		}

		if packageLine {
			packageName := strings.Split(line, " ")[1]
			gFile.goPackage = &goPackage{
				name: projectName + "/" + packageName,
			}
		} else if lineNotEmpty && !foundImportLines && !packageLine &&
			!strings.HasPrefix(line, "import") {
			break
		} else if foundImportLines && strings.HasPrefix(line, ")") {
			foundImportLines = false
		} else if singleImportLine {
			parts := strings.Split(line, " ")
			importedPackageName := parts[len(parts)-1]
			if isImportProjectPackageLine(importedPackageName) {
				depsPackages = append(depsPackages, &goPackage{
					name: removeQuotesFromPkg(importedPackageName),
				})
			}
		} else if foundImportLines && lineNotEmpty && isImportProjectPackageLine(line) {
			// parse multi line imports
			parts := strings.Split(line, " ")
			importedPackageName := parts[len(parts)-1]

			depsPackages = append(depsPackages, &goPackage{
				name: removeQuotesFromPkg(importedPackageName),
			})
		} else if strings.HasPrefix(line, "import") && !singleImportLine {
			// begin multi line imports
			foundImportLines = true
		}
	}
	gFile.deps = depsPackages
	return &gFile, nil
}

func getGoFiles(dir string) ([]string, error) {
	goFiles := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return goFiles, nil
}
