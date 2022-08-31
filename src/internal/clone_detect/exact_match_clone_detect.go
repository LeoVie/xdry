package clone_detect

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"x-dry-go/src/internal/cli"
	"x-dry-go/src/internal/compare"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/normalize"
	"x-dry-go/src/internal/structs"
)

type Clone = struct {
	A        string
	B        string
	Language string
	Matches  []compare.Match
}

type Pair = struct {
	AFile structs.File
	BFile structs.File
}

func DetectInDirectory(directory string, level int, levelNormalizers map[int][]config.Normalizer) (error, []Clone) {
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
		return err, []Clone{}
	}

	normalizeLevel := level
	if level > 2 {
		normalizeLevel = 2
	}

	normalizedFiles := normalizeFiles(
		normalizeLevel,
		levelNormalizers,
		filepaths,
	)

	var compareFunc func(string, string) []compare.Match
	if level == 1 || level == 2 {
		compareFunc = func(a string, b string) []compare.Match {
			return compare.FindExactMatches(a, b)
		}
	} else if level == 3 {
		compareFunc = func(a string, b string) []compare.Match {
			return compare.FindLongestCommonSubsequence(a, b)
		}
	} else {
		return fmt.Errorf("no compare function found for level %d", level), []Clone{}
	}

	clones := detectClones(normalizedFiles, compareFunc)

	return nil, clones
}

func detectClones(normalizedFiles map[string]structs.File, compareFunc func(a string, b string) []compare.Match) []Clone {
	pairs := make(map[string]Pair)

	for aPath, aFile := range normalizedFiles {
		for bPath, bFile := range normalizedFiles {
			if aPath == bPath {
				continue
			}

			firstFile, secondFile := orderFiles(aFile, bFile)

			hash := buildCloneHash(firstFile.Path, secondFile.Path)
			if _, ok := pairs[hash]; ok {
				continue
			}

			pairs[hash] = Pair{
				AFile: firstFile,
				BFile: secondFile,
			}
		}
	}

	var (
		clonesMutex sync.Mutex
		clones      []Clone
	)
	var clonesWg sync.WaitGroup

	for _, pair := range pairs {
		clonesWg.Add(1)

		go func(pair Pair) {
			defer clonesWg.Done()

			matches := compareFunc(pair.AFile.Content, pair.BFile.Content)

			if len(matches) == 0 {
				return
			}

			clonesMutex.Lock()
			clones = append(clones, Clone{
				A:        pair.AFile.Path,
				B:        pair.BFile.Path,
				Language: pair.AFile.Language,
				Matches:  matches,
			})
			clonesMutex.Unlock()
		}(pair)
	}

	clonesWg.Wait()
	return clones
}

func orderFiles(aFile structs.File, bFile structs.File) (structs.File, structs.File) {
	if aFile.Path < bFile.Path {
		return aFile, bFile
	}

	return bFile, aFile
}

func normalizeFiles(
	level int,
	levelNormalizers map[int][]config.Normalizer,
	filepaths []string,
) map[string]structs.File {
	var (
		normalizedFilesMutex sync.Mutex
		normalizedFiles      = make(map[string]structs.File)
	)

	normalizers, ok := levelNormalizers[level]
	if !ok {
		log.Printf("No normalizers configured for level %d\n", level)
		return normalizedFiles
	}

	mappedNormalizers := make(map[string]config.Normalizer)
	for _, normalizer := range normalizers {
		mappedNormalizers[normalizer.Extension] = normalizer
	}

	const max = 12
	semaphore := make(chan struct{}, max)
	wg := &sync.WaitGroup{}

	for _, path := range filepaths {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(path string, mappedNormalizers map[string]config.Normalizer) {
			defer wg.Done()

			err, normalizedFile := normalize.Normalize(path, mappedNormalizers, cli.NewCommandExecutor())

			if err != nil {
				fmt.Println(err)
			}

			normalizedFilesMutex.Lock()
			normalizedFiles[path] = normalizedFile
			normalizedFilesMutex.Unlock()

			<-semaphore
		}(path, mappedNormalizers)
	}
	wg.Wait()

	return normalizedFiles
}

func buildCloneHash(aPath string, bPath string) string {
	return aPath + "_" + bPath
}
