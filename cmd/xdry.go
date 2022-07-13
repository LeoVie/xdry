package main

import (
	"flag"
	"os"
	"x-dry-go/internal/service"
)

func main() {
	var configPath string

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.StringVar(&configPath, "config", "xdry.json", "Path to xdry.json config")
	commandLine.Parse(os.Args[1:])

	os.Exit(service.Analyze(os.Stdout, configPath))
}
