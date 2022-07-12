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

	cmd := exec.Command(normalizer.Command, hydrateArgs(path, normalizer)...)

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

func hydrateArgs(path string, normalizer *config.Normalizer) []string {
	var hydratedArgs []string
	for _, arg := range normalizer.Args {
		hydratedArgs = append(hydratedArgs, strings.ReplaceAll(arg, "%FILEPATH%", path))
	}

	return hydratedArgs
}
