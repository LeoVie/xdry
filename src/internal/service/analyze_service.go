package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"x-dry-go/src/internal/clone_detect"
	"x-dry-go/src/internal/compare"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/service/aggregate"
	"x-dry-go/src/internal/service/reporter"
)

const (
	CommandFailure = 1
	CommandSuccess = 0
)

func Analyze(out io.Writer, configPath string) int {
	configuration, err := readConfig(configPath)
	if err != nil {
		fmt.Println(out, err)

		return CommandFailure
	}

	logFile, err := os.OpenFile(configuration.Settings.LogPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(out, "error opening log file: %v", err)

		return CommandFailure
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile)

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

	clones := map[int][]clone_detect.Clone{
		1: relevantType1Clones,
		2: relevantType2Clones,
		3: relevantType3Clones,
	}

	cloneBundles := aggregate.AggregateCloneBundles(clones)

	cloneBundles = normalizeCloneBundles(cloneBundles)

	for _, report := range configuration.Reports {
		var err error

		if report.Type == "json" {
			err = reporter.WriteJsonReport(cloneBundles, report)
		} else if report.Type == "html" {
			err = reporter.WriteHtmlReport(cloneBundles, report)
		}
		if err != nil {
			fmt.Fprintln(out, err)

			return CommandFailure
		}
	}

	return CommandSuccess
}

func normalizeCloneBundles(cloneBundles []aggregate.CloneBundle) []aggregate.CloneBundle {
	for _, bundle := range cloneBundles {
		sort.Slice(bundle.AggregatedClones, func(i, j int) bool {
			return len(bundle.AggregatedClones[i].Content) < len(bundle.AggregatedClones[j].Content)
		})
	}
	sort.Slice(cloneBundles, func(i, j int) bool {
		return cloneBundles[i].CloneType < cloneBundles[j].CloneType
	})

	return cloneBundles
}

func readConfig(configPath string) (*config.Config, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not find config file '%s'", configPath)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absoluteConfigPath := convertConfigPathToAbsolutePath(configPath, cwd)
	err, configuration := config.ParseConfig(absoluteConfigPath, cwd)

	if err != nil {
		return nil, fmt.Errorf("error while parsing config under '%s'", absoluteConfigPath)
	}

	return configuration, nil
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
			A:        clone.A,
			B:        clone.B,
			Language: clone.Language,
			Matches:  filteredMatches,
		})
	}

	return filtered
}
