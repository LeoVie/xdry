package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"x-dry-go/internal/clone_detect"
	"x-dry-go/internal/compare"
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
	type3Clones := make(map[string]clone_detect.Clone)
	for _, directory := range configuration.Directories {
		err, type1ClonesInDir := clone_detect.DetectInDirectory(directory, 1, levelNormalizers)
		if err != nil {
			fmt.Fprintln(out, err)

			return CommandFailure
		}
		err, type2ClonesInDir := clone_detect.DetectInDirectory(directory, 2, levelNormalizers)
		if err != nil {
			fmt.Fprintln(out, err)

			return CommandFailure
		}
		err, type3ClonesInDir := clone_detect.DetectInDirectory(directory, 3, levelNormalizers)
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
		for key, clone := range type3ClonesInDir {
			type3Clones[key] = clone
		}
	}

	relevantType1Clones := filterClonesByLength(type1Clones, configuration.Settings.MinCloneLengths["level-1"])
	relevantType2Clones := filterClonesByLength(type2Clones, configuration.Settings.MinCloneLengths["level-2"])
	relevantType3Clones := filterClonesByLength(type3Clones, configuration.Settings.MinCloneLengths["level-3"])

	clones := map[string]map[string]clone_detect.Clone{
		"TYPE 1": relevantType1Clones,
		"TYPE 2": relevantType2Clones,
		"TYPE 3": relevantType3Clones,
	}

	for _, report := range configuration.Reports {
		if report.Type == "json" {
			err := writeJsonReport(clones, report)

			if err != nil {
				fmt.Fprintln(out, err)

				return CommandFailure
			}
		}
	}

	return CommandSuccess
}

func writeJsonReport(clones map[string]map[string]clone_detect.Clone, report config.Report) error {
	jsonStr, err := json.Marshal(clones)
	err = os.WriteFile(report.Path, jsonStr, 0644)

	return err
}

func convertConfigPathToAbsolutePath(configPath string, cwd string) string {
	if strings.HasPrefix(configPath, "/") {
		return configPath
	}

	return path.Join(cwd, configPath)
}

func filterClonesByLength(clones map[string]clone_detect.Clone, minLength int) map[string]clone_detect.Clone {
	filtered := make(map[string]clone_detect.Clone)

	for key, clone := range clones {
		filteredMatches := []compare.Match{}

		for _, match := range clone.Matches {
			if len(match.Content) < minLength {
				continue
			}

			filteredMatches = append(filteredMatches, match)
		}

		if len(filteredMatches) == 0 {
			continue
		}

		filtered[key] = clone_detect.Clone{
			A:       clone.A,
			B:       clone.B,
			Matches: filteredMatches,
		}
	}

	return filtered
}
