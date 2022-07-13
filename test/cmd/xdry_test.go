package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"x-dry-go/internal/service"
)

func TestRun(t *testing.T) {
	generatedReportFile := "generated/reports/xdry_report.json"
	os.Remove(generatedReportFile)

	assert.NoFileExists(t, generatedReportFile)

	wantBytes, _ := os.ReadFile("expected_xdry_report.json")
	want := string(wantBytes)

	cwd, _ := os.Getwd()

	configPath := cwd + string(os.PathSeparator) + "xdry.json"

	var out io.Writer = os.Stdout

	service.Analyze(out, configPath)

	actualBytes, _ := os.ReadFile(generatedReportFile)
	actual := string(actualBytes)

	assert.JSONEq(t, want, actual)
}
