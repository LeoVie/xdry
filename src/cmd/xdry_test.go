package main

import (
	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/gomega"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"x-dry-go/src/internal/service"
)

func TestRun(t *testing.T) {
	g := NewGomegaWithT(t)

	cwd, _ := os.Getwd()
	projectPath := path.Join(cwd, "..", "..")
	testdataPath := path.Join(projectPath, "_testdata")
	reportsDir := path.Join(testdataPath, "generated", "reports")
	cacheDir := path.Join(testdataPath, "cache")

	deletedReports := clearReportsDir(reportsDir)
	deletedCaches := clearCacheDir(cacheDir)

	for _, deletedReport := range deletedReports {
		g.Expect(deletedReport).ShouldNot(BeAnExistingFile())
	}
	for _, deletedCache := range deletedCaches {
		g.Expect(deletedCache).ShouldNot(BeAnExistingFile())
	}

	wantFile := path.Join(testdataPath, "expected_xdry_report.json")
	wantBytes, _ := os.ReadFile(wantFile)
	want := string(wantBytes)

	want = strings.ReplaceAll(want, "%PROJECT_PATH%", projectPath)

	configPath := path.Join(testdataPath, "xdry_2.json")

	var out io.Writer = os.Stdout

	service.Analyze(out, configPath)

	generatedReportFile := path.Join(reportsDir, "xdry_report.json")
	actualBytes, _ := os.ReadFile(generatedReportFile)
	actual := string(actualBytes)

	g.Expect(actual).Should(
		MatchUnorderedJSON(want))
}

func clearReportsDir(reportsDir string) []string {
	files, err := filepath.Glob(path.Join(reportsDir, "xdry_report.*"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}

	return files
}

func clearCacheDir(cacheDir string) []string {
	files, err := filepath.Glob(path.Join(cacheDir, "xdry-cache_*.json"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}

	return files
}
