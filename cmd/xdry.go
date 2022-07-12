package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"x-dry-go/internal/clone_detect"
	"x-dry-go/internal/config"
)

var (
	configPath *string
)

func init() {
	configPath = flag.String("config", "xdry.json", "Path to xdry.json config")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*configPath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("could not find config file '%s'", *configPath)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	absoluteConfigPath := convertConfigPathToAbsolutePath(*configPath, cwd)
	err, configuration := config.ParseConfig(absoluteConfigPath, cwd)

	if err != nil {
		log.Fatalf("error while parsing config under '%s'", absoluteConfigPath)
	}

	var levelNormalizers = make(map[int][]config.Normalizer)

	for _, normalizer := range configuration.Normalizers {
		levelNormalizers[normalizer.Level] = append(levelNormalizers[normalizer.Level], normalizer)
	}

	type1Clones := make(map[string]clone_detect.Clone)
	type2Clones := make(map[string]clone_detect.Clone)
	for _, directory := range configuration.Directories {
		err, type1ClonesInDir := clone_detect.DetectInDirectory(directory, 1, levelNormalizers)
		if err != nil {
			log.Fatal(err)
		}
		err, type2ClonesInDir := clone_detect.DetectInDirectory(directory, 2, levelNormalizers)
		if err != nil {
			log.Fatal(err)
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

	fmt.Println("Type 1:")
	outputClones(relevantType1Clones)
	fmt.Println("Type 2:")
	outputClones(relevantType2Clones)

	fmt.Printf("\nFound %d type1Clones.\n", len(relevantType1Clones))
	fmt.Printf("\nFound %d type2Clones.\n", len(relevantType2Clones))
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

func outputClones(clones map[string]clone_detect.Clone) {
	for _, clone := range clones {
		fmt.Printf(
			"A: %s\nB: %s\nLength: %d, Match: %s\n\n",
			clone.A,
			clone.B,
			len(clone.Match),
			clone.Match,
		)
	}
}
