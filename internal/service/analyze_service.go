package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"x-dry-go/internal/clone_detect"
	"x-dry-go/internal/config"
)

const (
	CommandFailure = 1
	CommandSuccess = 0
)

func Analyze(out io.Writer, configPath string) int {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(out, "could not find config file '%s'\n", configPath)

		return CommandFailure
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(out, err)

		return CommandFailure
	}

	absoluteConfigPath := convertConfigPathToAbsolutePath(configPath, cwd)
	err, configuration := config.ParseConfig(absoluteConfigPath, cwd)

	if err != nil {
		fmt.Fprintf(out, "error while parsing config under '%s'", absoluteConfigPath)

		return CommandFailure
	}

	var levelNormalizers = make(map[int][]config.Normalizer)

	for _, normalizer := range configuration.Normalizers {
		levelNormalizers[normalizer.Level] = append(levelNormalizers[normalizer.Level], normalizer)
	}

	type1Clones := make(map[string]clone_detect.Clone)
	type2Clones := make(map[string]clone_detect.Clone)
	for _, directory := range configuration.Directories {
		err, type1ClonesInDir := clone_detect.DetectInDirectory(directory, CommandFailure, levelNormalizers)
		if err != nil {
			fmt.Fprintln(out, err)

			return CommandFailure
		}
		err, type2ClonesInDir := clone_detect.DetectInDirectory(directory, 2, levelNormalizers)
		if err != nil {
			fmt.Fprintln(out, err)

			return CommandFailure
		}

		for key, clone := range type1ClonesInDir {
			type1Clones[key] = clone
		}
		for key, clone := range type2ClonesInDir {
			type2Clones[key] = clone
		}
	}

	relevantType1Clones := filterClonesByLength(type1Clones, 20)
	relevantType2Clones := filterClonesByLength(type2Clones, 20)

	clones := map[string]map[string]clone_detect.Clone{
		"TYPE 1": relevantType1Clones,
		"TYPE 2": relevantType2Clones,
	}

	for _, report := range configuration.Reports {
		if report.Type == "json" {
			jsonStr, err := json.Marshal(clones)

			if err != nil {
				fmt.Fprintln(out, err)

				return CommandFailure
			}

			err = os.WriteFile(report.Path, jsonStr, 0644)

			if err != nil {
				fmt.Fprintln(out, err)

				return CommandFailure
			}
		}
	}

	return CommandSuccess
}

func convertConfigPathToAbsolutePath(configPath string, cwd string) string {
	if strings.HasPrefix(configPath, "/") {
		return configPath
	}

	return cwd + string(os.PathSeparator) + configPath
}

func filterClonesByLength(clones map[string]clone_detect.Clone, minLength int) map[string]clone_detect.Clone {
	filtered := make(map[string]clone_detect.Clone)

	for key, clone := range clones {
		if len(clone.Match) < minLength {
			continue
		}

		filtered[key] = clone
	}

	return filtered
}
