package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := path.Join(cwd, "..", "..", "..", "_testdata", "xdry_1.json")

	want := Config{
		Settings: Settings{
			MinCloneLengths: map[string]int{
				"level-1": 10,
				"level-2": 20,
			},
		},
		Directories: []string{
			path.Join(cwd, "test", "_testdata", "php"),
			path.Join(cwd, "..", "..", "..", "_testdata", "test", "_testdata", "javascript"),
		},
		Normalizers: []Normalizer{
			{
				Level:     1,
				Extension: ".php",
				Command:   "php",
				Args: []string{
					"%FILEPATH%",
				},
			},
		},
	}

	_, actual := ParseConfig(configPath, cwd)

	assert.Equal(t, &want, actual)
}