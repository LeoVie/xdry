package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"strings"
	"testing"
	"x-dry-go/src/internal/service"
)

func TestRun(t *testing.T) {
	cwd, _ := os.Getwd()
	projectPath := path.Join(cwd, "..", "..")

	generatedReportFile := path.Join(projectPath, "_testdata", "generated", "reports", "xdry_report.json")

	fmt.Println(generatedReportFile)

	os.Remove(generatedReportFile)

	assert.NoFileExists(t, generatedReportFile)

	wantFile := path.Join(projectPath, "_testdata", "expected_xdry_report.json")
	wantBytes, _ := os.ReadFile(wantFile)
	want := string(wantBytes)

	want = strings.ReplaceAll(want, "%PROJECT_PATH%", projectPath)

	configPath := path.Join(projectPath, "_testdata", "xdry_2.json")

	var out io.Writer = os.Stdout

	service.Analyze(out, configPath)

	actualBytes, _ := os.ReadFile(generatedReportFile)
	actual := string(actualBytes)

	assert.JSONEq(t, want, actual)
}