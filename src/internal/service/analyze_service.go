package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"x-dry-go/src/internal/clone_detect"
	"x-dry-go/src/internal/compare"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/service/reporter"
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

	var type1Clones []clone_detect.Clone
	var type2Clones []clone_detect.Clone
	var type3Clones []clone_detect.Clone
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

		for _, clone := range type1ClonesInDir {
			type1Clones = append(type1Clones, clone)
		}
		for _, clone := range type2ClonesInDir {
			type2Clones = append(type2Clones, clone)
		}
		for _, clone := range type3ClonesInDir {
			type2Clones = append(type2Clones, clone)
		}
	}

	relevantType1Clones := filterClonesByLength(type1Clones, configuration.Settings.MinCloneLengths["level-1"])
	relevantType2Clones := filterClonesByLength(type2Clones, configuration.Settings.MinCloneLengths["level-2"])
	relevantType3Clones := filterClonesByLength(type3Clones, configuration.Settings.MinCloneLengths["level-3"])

	clones := map[string][]clone_detect.Clone{
		"TYPE 1": relevantType1Clones,
		"TYPE 2": relevantType2Clones,
		"TYPE 3": relevantType3Clones,
	}

	for _, report := range configuration.Reports {
		if report.Type == "json" {
			err := reporter.WriteJsonReport(clones, report)

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

	return path.Join(cwd, configPath)
}

func filterClonesByLength(clones []clone_detect.Clone, minLength int) []clone_detect.Clone {
	var filtered []clone_detect.Clone
	for _, clone := range clones {
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

		filtered = append(filtered, clone_detect.Clone{
			A:       clone.A,
			B:       clone.B,
			Matches: filteredMatches,
		})
	}

	return filtered
}
