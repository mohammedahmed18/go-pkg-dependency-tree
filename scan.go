package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func scanPackages(path string, projectName string) {
	goFilesPaths, _ := getGoFiles(path)
	gofiles := make([]*goFile, 0)
	for _, path := range goFilesPaths {
		f, err := parseGoFile(path, projectName)
		if err == nil {
			gofiles = append(gofiles, f)
		}
	}

	generatePackagesGraph(gofiles)
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
		// Process the line
		line := scanner.Text()
		line = strings.TrimSpace(line)

		lineNotEmpty := line != ""
		packageLine := strings.HasPrefix(line, "package")
		importProjectPackageLine := func(l string) bool { return strings.HasPrefix(l, "\""+projectName) }
		commentLine := strings.HasPrefix(line, "//")

		if commentLine {
			continue
		}

		singleImportLine := strings.HasPrefix(line, "import") && !strings.Contains(line, "(")

		if lineNotEmpty && !foundImportLines && !strings.HasPrefix(line, "package") &&
			!strings.HasPrefix(line, "import") {
			break
		} else if foundImportLines && strings.HasPrefix(line, ")") {
			foundImportLines = false
			break
		} else if packageLine {
			packageName := strings.Split(line, " ")[1]
			gFile.goPackage = &goPackage{
				name: projectName + "/" + packageName,
			}
		} else if singleImportLine {
			importedPackageName := strings.Split(line, " ")[1]
			isProjectPackageImport := importProjectPackageLine(importedPackageName)
			if isProjectPackageImport {
				importedPackageName = strings.TrimPrefix(importedPackageName, "\"")
				importedPackageName = strings.TrimSuffix(importedPackageName, "\"")
				depsPackages = append(depsPackages, &goPackage{
					name: importedPackageName,
				})
			}
		} else if foundImportLines && lineNotEmpty {
			if importProjectPackageLine(line) {

				importedPackageName := strings.TrimPrefix(line, "\"")
				importedPackageName = strings.TrimSuffix(importedPackageName, "\"")
				depsPackages = append(depsPackages, &goPackage{
					name: importedPackageName,
				})
			} else if strings.Contains(line, " ") {
				importedPackageName := strings.Split(line, " ")[1]
				isProjectPackageImport := importProjectPackageLine(importedPackageName)
				if isProjectPackageImport {
					importedPackageName = strings.TrimPrefix(importedPackageName, "\"")
					importedPackageName = strings.TrimSuffix(importedPackageName, "\"")
					depsPackages = append(depsPackages, &goPackage{
						name: importedPackageName,
					})
				}
			}
		} else if strings.HasPrefix(line, "import") {
			// multi line imports
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
