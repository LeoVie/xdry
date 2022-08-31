package reporter

import (
	"bytes"
	"github.com/yosssi/gohtml"
	"os"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/service/aggregate"
	"x-dry-go/src/templates"
)

func WriteHtmlReport(cloneBundles []aggregate.CloneBundle, report config.Report) error {
	var buf bytes.Buffer
	templates.WriteReport(&buf, cloneBundles)

	err := os.WriteFile(report.Path, []byte(gohtml.Format(buf.String())), 0644)

	return err
}
