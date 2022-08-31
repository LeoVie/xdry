package reporter

import (
	"encoding/json"
	"os"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/service/aggregate"
)

func WriteJsonReport(cloneBundles []aggregate.CloneBundle, report config.Report) error {
	jsonStr, err := json.Marshal(cloneBundles)
	err = os.WriteFile(report.Path, jsonStr, 0644)

	return err
}
