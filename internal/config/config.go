package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Normalizers []Normalizer `json:"normalizers"`
}

type Normalizer struct {
	Level     int      `json:"level"`
	Extension string   `json:"extension"`
	Command   string   `json:"command"`
	Args      []string `json:"args"`
}

func ParseConfig(configPath string) (error, *Config) {
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

	return err, &config
}
