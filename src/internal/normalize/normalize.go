package normalize

import (
	"fmt"
	"path/filepath"
	"strings"
	"x-dry-go/src/internal/cli"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/structs"
)

func Normalize(path string, normalizers map[string]config.Normalizer, commandExecutor cli.CommandExecutor) (error, structs.File) {
	err, normalizer := findNormalizeImplementation(path, normalizers)
	if err != nil {
		return err, structs.File{}
	}

	commandOutput, err := commandExecutor.Execute(normalizer.Command, hydrateArgs(path, normalizer))
	if err != nil {
		return err, structs.File{}
	}

	return nil, structs.File{
		Path:     path,
		Content:  commandOutput,
		Language: normalizer.Language,
	}
}

func findNormalizeImplementation(
	path string,
	normalizers map[string]config.Normalizer,
) (error, *config.Normalizer) {
	fileExtension := filepath.Ext(path)

	if normalizer, ok := normalizers[fileExtension]; ok {
		return nil, &normalizer
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
