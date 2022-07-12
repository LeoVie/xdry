package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"x-dry-go/internal/config"
)

func TestParseConfig(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := cwd + string(os.PathSeparator) + "xdry.json"

	want := config.Config{
		Directories: []string{
			cwd + string(os.PathSeparator) + "test/_testdata/php/",
			cwd + string(os.PathSeparator) + "./test/_testdata/javascript/",
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
