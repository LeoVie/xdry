package reporter

import (
	"encoding/json"
	"os"
	"x-dry-go/src/internal/clone_detect"
	"x-dry-go/src/internal/config"
)

func WriteJsonReport(clones map[string][]clone_detect.Clone, report config.Report) error {
	jsonStr, err := json.Marshal(clones)
	err = os.WriteFile(report.Path, jsonStr, 0644)

	return err
}
