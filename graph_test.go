package main

import (
	"os"
	"path"
	"testing"
)

func getSamplePath(sampleNum string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(pwd, "test_samples/sample"+sampleNum), nil
}

func checkPkgDeps(t *testing.T, pkgDeps []*goPackage, expectedPkgDeps []string) bool {

	set := make(map[string]struct{})

	for _, item := range pkgDeps {
		set[item.name] = struct{}{}
	}

	for _, item := range expectedPkgDeps {
		if _, found := set[item]; !found {
			t.Errorf("couldn't find package %s in pkg2 deps, pkg2 deps are %v", item, set)
		}
	}

	return true

}
func TestSample1(t *testing.T) {
	sample1_path, err := getSamplePath("1")
	if err != nil {
		t.Error(err)
	}

	projectName, err := getProjectNameFromGoMod(sample1_path)
	if err != nil {
		t.Error(err)
	}

	if projectName != "github.com/samples/sample1" {
		t.Errorf("package name should be %s instead we got %s", "github.com/samples/sample1", projectName)
	}
	pkgMap, err := scanPackages(sample1_path, projectName)
	if err != nil {
		t.Error(err)
	}

	pkg := "github.com/samples/sample1/pkg1"
	checkPkgDeps(
		t,
		pkgMap[pkg],
		[]string{
			"github.com/samples/sample1/pkg2",
		})

	pkg = "github.com/samples/sample1/pkg2"
	checkPkgDeps(
		t,
		pkgMap[pkg],
		[]string{})

	pkg = "github.com/samples/sample1/main"
	checkPkgDeps(
		t,
		pkgMap[pkg],
		[]string{
			"github.com/samples/sample1/pkg1",
		})

}
