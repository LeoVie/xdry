package config

import (
	. "github.com/onsi/gomega"
	"log"
	"os"
	"path"
	"testing"
)

func TestParseConfig(t *testing.T) {
	g := NewGomegaWithT(t)

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
			LogPath:        path.Join(cwd, "test", "_testdata", "xdry.log"),
			CacheDirectory: path.Join(cwd, "test", "_testdata", "cache"),
		},
		Directories: []string{
			path.Join(cwd, "test", "_testdata", "php"),
			path.Join(cwd, "..", "..", "..", "_testdata", "test", "_testdata", "javascript"),
		},
		Normalizers: []Normalizer{
			{
				Level:     1,
				Extension: ".php",
				Language:  "php",
				Command:   "php",
				Args: []string{
					"%FILEPATH%",
				},
			},
		},
	}

	_, actual := ParseConfig(configPath, cwd)

	g.Expect(actual).To(Equal(&want))
}

func TestParseConfigMinimal(t *testing.T) {
	g := NewGomegaWithT(t)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := path.Join(cwd, "..", "..", "..", "_testdata", "xdry_minimal.json")

	want := Config{
		Settings: Settings{
			MinCloneLengths: nil,
			LogPath:         path.Join(cwd, "..", "..", "..", "_testdata", "xdry.log"),
			CacheDirectory:  path.Join(cwd, "..", "..", "..", "_testdata"),
		},
		Reports:     nil,
		Directories: nil,
		Normalizers: nil,
	}

	_, actual := ParseConfig(configPath, cwd)

	g.Expect(actual).To(Equal(&want))
}
