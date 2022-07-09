package main

import (
	"fmt"
	"log"
	"os"
	"x-dry-go/internal/clone_detect"
	"x-dry-go/internal/config"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("please pass directory")
	}
	directory := args[0]

	configPath := "/home/ubuntu/development/x-dry-go/xdry.json"
	err, configuration := config.ParseConfig(configPath)

	if err != nil {
		log.Fatalf("error while parsing config under '%s'", configPath)
	}

	var levelNormalizers = make(map[int][]config.Normalizer)

	for _, normalizer := range configuration.Normalizers {
		levelNormalizers[normalizer.Level] = append(levelNormalizers[normalizer.Level], normalizer)
	}

	err, type1Clones := clone_detect.DetectInDirectory(directory, 1, levelNormalizers)
	if err != nil {
		log.Fatal(err)
	}
	err, type2Clones := clone_detect.DetectInDirectory(directory, 2, levelNormalizers)
	if err != nil {
		log.Fatal(err)
	}

	for _, clone := range type1Clones {
		fmt.Printf(
			"%s;%s;%d\n",
			clone.A,
			clone.B,
			len(clone.Match),
		)
	}
	for _, clone := range type2Clones {
		fmt.Printf(
			"%s;%s;%d\n",
			clone.A,
			clone.B,
			len(clone.Match),
		)
	}

	fmt.Printf("\nFound %d type1Clones.\n", len(type1Clones))
	fmt.Printf("\nFound %d type2Clones.\n", len(type2Clones))
}
