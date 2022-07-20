package normalize

import (
	"fmt"
	"path/filepath"
	"strings"
	"x-dry-go/src/internal/cli"
	"x-dry-go/src/internal/config"
)

func Normalize(path string, normalizers []config.Normalizer, commandExecutor cli.CommandExecutor) (error, string) {
	err, normalizer := findNormalizeImplementation(path, normalizers)

	if err != nil {
		return err, ""
	}

	commandOutput, err := commandExecutor.Execute(normalizer.Command, hydrateArgs(path, normalizer))

	if err != nil {
		return err, ""
	}

	return nil, commandOutput
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
