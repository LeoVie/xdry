package normalize

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"x-dry-go/internal/config"
)

func Normalize(path string, normalizers []config.Normalizer) (error, string) {
	err, normalizer := findNormalizeImplementation(path, normalizers)

	if err != nil {
		return err, ""
	}

	var replacedArgs []string
	for _, arg := range normalizer.Args {
		replacedArgs = append(replacedArgs, strings.ReplaceAll(arg, "%FILEPATH%", path))
	}

	cmd := exec.Command(normalizer.Command, replacedArgs...)

	stdout, err := cmd.Output()

	if err != nil {
		return err, ""
	}

	return nil, string(stdout)
}

func findNormalizeImplementation(
	path string,
	normalizers []config.Normalizer,
) (error, *config.Normalizer) {
	fileExtension := filepath.Ext(path)

	for _, normalizer := range normalizers {
		if normalizer.Extension == fileExtension {
			return nil, &normalizer
		}
	}

	return fmt.Errorf("no normalizer found for file extension '%s'", fileExtension), nil
}
