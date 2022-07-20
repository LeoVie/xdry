package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path"
	"testing"
	"x-dry-go/internal/config"
)

func TestParseConfig(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := path.Join(cwd, "xdry.json")

	want := config.Config{
		Settings: config.Settings{
			MinCloneLengths: map[string]int{
				"level-1": 10,
				"level-2": 20,
			},
		},
		Directories: []string{
			path.Join(cwd, "test", "_testdata", "php"),
			path.Join(cwd, ".", "test", "_testdata", "javascript"),
		},
		Normalizers: []config.Normalizer{
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

	_, actual := config.ParseConfig(configPath, cwd)

	assert.Equal(t, &want, actual)
}
