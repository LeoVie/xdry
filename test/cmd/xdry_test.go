package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"strings"
	"testing"
	"x-dry-go/internal/service"
)

func TestRun(t *testing.T) {
	generatedReportFile := "generated/reports/xdry_report.json"
	os.Remove(generatedReportFile)

	assert.NoFileExists(t, generatedReportFile)

	cwd, _ := os.Getwd()

	projectPath := path.Join(cwd, "..", "..")

	fmt.Println(projectPath)

	wantBytes, _ := os.ReadFile("expected_xdry_report.json")
	want := string(wantBytes)
	want = strings.ReplaceAll(want, "%PROJECT_PATH%", projectPath)

	configPath := path.Join(cwd, "xdry.json")

	var out io.Writer = os.Stdout

	service.Analyze(out, configPath)

	actualBytes, _ := os.ReadFile(generatedReportFile)
	actual := string(actualBytes)

	assert.JSONEq(t, want, actual)
}
