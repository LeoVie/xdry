package clone_detect

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"x-dry-go/internal/compare"
	"x-dry-go/internal/config"
	"x-dry-go/internal/normalize"
)

type Clone = struct {
	A     string
	B     string
	Match string
}

func DetectInDirectory(directory string, level int, levelNormalizers map[int][]config.Normalizer) (error, map[string]Clone) {
	var filepaths []string
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filepaths = append(filepaths, path)
			}
			return nil
		})
	if err != nil {
		return err, map[string]Clone{}
	}

	normalizedFileContents := normalizeFiles(
		level,
		levelNormalizers,
		filepaths,
	)
	clones := detectClones(normalizedFileContents)

	return nil, clones
}

func detectClones(normalizedFileContents map[string]string) map[string]Clone {
	var (
		clonesMutex sync.Mutex
		clones      = make(map[string]Clone)
	)

	var clonesWg sync.WaitGroup

	for aPath, aContent := range normalizedFileContents {
		for bPath, bContent := range normalizedFileContents {
			clonesWg.Add(1)

			go func(aPath string, aContent string, bPath string, bContent string) {
				defer clonesWg.Done()

				if aPath == bPath {
					return
				}

				longestMatch := compare.FindLongestMatch(aContent, bContent)

				if longestMatch == "" {
					return
				}

				hash := calculateCloneHash(aPath, bPath)

				var first string
				var second string
				if aPath < bPath {
					first = aPath
					second = bPath
				} else {
					first = bPath
					second = aPath
				}

				clonesMutex.Lock()
				clones[hash] = Clone{
					A:     first,
					B:     second,
					Match: longestMatch,
				}
				clonesMutex.Unlock()
			}(aPath, aContent, bPath, bContent)
		}
	}
	clonesWg.Wait()
	return clones
}

func normalizeFiles(
	level int,
	levelNormalizers map[int][]config.Normalizer,
	filepaths []string,
) map[string]string {
	var (
		normalizedFileContentsMutex sync.Mutex
		normalizedFileContents      = make(map[string]string)
	)

	normalizers, ok := levelNormalizers[level]
	if !ok {
		log.Printf("No normalizers configured for level %d\n", level)
		return normalizedFileContents
	}

	var wg sync.WaitGroup
	for _, path := range filepaths {
		wg.Add(1)

		go func(path string, normalizers []config.Normalizer) {
			defer wg.Done()

			err, normalizedFileContent := normalize.Normalize(path, normalizers)

			if err != nil {
				fmt.Println(err)
			}

			normalizedFileContentsMutex.Lock()
			normalizedFileContents[path] = normalizedFileContent
			normalizedFileContentsMutex.Unlock()
		}(path, normalizers)
	}
	wg.Wait()

	return normalizedFileContents
}

func calculateCloneHash(aPath string, bPath string) string {
	if aPath < bPath {
		return aPath + "_" + bPath
	}

	return bPath + "_" + aPath
}
