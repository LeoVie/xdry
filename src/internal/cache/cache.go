package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"x-dry-go/src/internal/structs"
)

type FileCache struct {
	Path  string
	Items map[string]structs.File `json:"items"`
}

func InitOrReadCache(path string) (*FileCache, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return initCache(path)
	}

	return readFileCache(path)
}

func initCache(path string) (*FileCache, error) {
	cache := FileCache{
		Path:  path,
		Items: make(map[string]structs.File),
	}

	jsonFile, err := json.MarshalIndent(cache, "", " ")
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(cache.Path, jsonFile, 0644)
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func readFileCache(path string) (*FileCache, error) {
	var cache FileCache

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &cache)
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func Store(cache FileCache, file structs.File) error {
	cache.Items[file.Path] = file

	jsonFile, err := json.MarshalIndent(cache, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cache.Path, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Get(cache FileCache, key string) (*structs.File, error) {
	if file, ok := cache.Items[key]; ok {
		return &file, nil
	}

	return nil, errors.New("cache item does not exist")
}
