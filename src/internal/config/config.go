package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Config struct {
	Settings    Settings     `json:"settings"`
	Reports     []Report     `json:"reports"`
	Directories []string     `json:"directories"`
	Normalizers []Normalizer `json:"normalizers"`
	LogPath     string       `json:"logPath"`
}

type Settings struct {
	MinCloneLengths map[string]int `json:"minCloneLengths"`
}

type Report struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type Normalizer struct {
	Level     int      `json:"level"`
	Extension string   `json:"extension"`
	Language  string   `json:"language"`
	Command   string   `json:"command"`
	Args      []string `json:"args"`
}

func ParseConfig(configPath string, cwd string) (error, *Config) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return err, nil
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err, nil
	}

	hydrateReports(&config, configPath, cwd)
	hydrateDirectories(&config, configPath, cwd)
	hydrateLogPath(&config, configPath, cwd)

	return nil, &config
}

func hydrateReports(config *Config, configPath string, cwd string) {
	configDir := path.Dir(configPath)

	var hydratedReports []Report
	for _, report := range config.Reports {
		hydratedReports = append(
			hydratedReports,
			Report{
				Type: report.Type,
				Path: toAbsolutePath(report.Path, configDir, cwd),
			},
		)
	}
	config.Reports = hydratedReports
}

func hydrateDirectories(config *Config, configPath string, cwd string) {
	configDir := path.Dir(configPath)

	fmt.Printf("ConfigDir: %s\n", configDir)

	var hydratedDirectories []string
	for _, directory := range config.Directories {
		hydratedDirectories = append(
			hydratedDirectories,
			toAbsolutePath(directory, configDir, cwd),
		)
	}
	config.Directories = hydratedDirectories
}

func hydrateLogPath(config *Config, configPath string, cwd string) {
	if config.LogPath == "" {
		config.LogPath = "xdry.log"
	}

	configDir := path.Dir(configPath)
	config.LogPath = toAbsolutePath(config.LogPath, configDir, cwd)
}

func toAbsolutePath(directory string, configDir string, cwd string) string {
	if strings.HasPrefix(directory, "/") {
		return directory
	}

	if strings.HasPrefix(directory, "%pwd%") {
		return strings.ReplaceAll(directory, "%pwd%", cwd)
	}

	return path.Join(configDir, directory)
}
