package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Config struct {
	Directories []string     `json:"directories"`
	Normalizers []Normalizer `json:"normalizers"`
}

type Normalizer struct {
	Level     int      `json:"level"`
	Extension string   `json:"extension"`
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

	err = hydrateDirectories(&config, configPath, cwd)
	if err != nil {
		return err, nil
	}

	return nil, &config
}

func hydrateDirectories(config *Config, configPath string, cwd string) error {
	configDir := path.Dir(configPath)

	var hydratedDirectories []string
	for _, directory := range config.Directories {
		hydratedDirectories = append(
			hydratedDirectories,
			convertDirectoryToAbsolutePath(directory, configDir, cwd),
		)
	}
	config.Directories = hydratedDirectories

	return nil
}

func convertDirectoryToAbsolutePath(directory string, configDir string, cwd string) string {
	if strings.HasPrefix(directory, "/") {
		return directory
	}

	if strings.HasPrefix(directory, "%pwd%") {
		return strings.ReplaceAll(directory, "%pwd%", cwd)
	}

	return configDir + string(os.PathSeparator) + directory
}
