package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
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

	typedClones := map[int][]clone_detect.Clone{}
	for _, directory := range configuration.Directories {
		for cloneType := 1; cloneType <= 3; cloneType++ {
			err, clonesInDir := clone_detect.DetectInDirectory(directory, cloneType, levelNormalizers, *configuration)
			if err != nil {
				fmt.Fprintln(out, err)

				return CommandFailure
			}

			clonesInDir = filterClonesByLength(
				clonesInDir,
				configuration.Settings.MinCloneLengths["level-"+strconv.Itoa(cloneType)],
			)

			typedClones[cloneType] = append(typedClones[cloneType], clonesInDir...)
		}
	}

	cloneBundles := aggregate.AggregateCloneBundles(typedClones)
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
			a := bundle.AggregatedClones[i].Content
			b := bundle.AggregatedClones[j].Content

			if len(a) == len(b) {
				return a < b
			}

			return len(a) < len(b)
		})
		for _, aggrClone := range bundle.AggregatedClones {
			sort.Slice(aggrClone.Instances, func(i, j int) bool {
				return aggrClone.Instances[i].Path < aggrClone.Instances[j].Path
			})
		}
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

	absoluteConfigPath := toAbsolutePath(configPath, cwd)
	err, configuration := config.ParseConfig(absoluteConfigPath, cwd)

	if err != nil {
		return nil, fmt.Errorf("error while parsing config under '%s'", absoluteConfigPath)
	}

	return configuration, nil
}

func toAbsolutePath(p string, cwd string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}

	return path.Join(cwd, p)
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
