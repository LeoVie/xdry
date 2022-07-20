package main

import (
	"flag"
	"log"
	"os"
	"path"
	"x-dry-go/src/internal/service"
)

func main() {
	var configPath string

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.StringVar(&configPath, "config", determineDefaultConfigPath(), "Path to xdry.json config")
	commandLine.Parse(os.Args[1:])

	os.Exit(service.Analyze(os.Stdout, configPath))
}

func determineDefaultConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(cwd, "xdry.json")
}
