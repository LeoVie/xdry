package clone_detect

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"x-dry-go/src/internal/cache"
	"x-dry-go/src/internal/cli"
	"x-dry-go/src/internal/compare"
	"x-dry-go/src/internal/config"
	"x-dry-go/src/internal/normalize"
	"x-dry-go/src/internal/structs"
)

type Clone struct {
	A        string
	B        string
	Language string
	Matches  []compare.Match
}

type filePair struct {
	AFile structs.File
	BFile structs.File
}

func (pair filePair) sort() filePair {
	if pair.AFile.Path < pair.BFile.Path {
		return pair
	}

	return filePair{
		AFile: pair.BFile,
		BFile: pair.AFile,
	}
}

func (pair filePair) hash() string {
	return pair.AFile.Path + "_" + pair.BFile.Path
}

func DetectInDirectory(
	directory string,
	cloneType int,
	levelNormalizers map[int][]config.Normalizer,
	configuration config.Config,
) (error, []Clone) {
	filepaths, err := findFilesInDir(directory)
	if err != nil {
		return err, []Clone{}
	}

	normalizedFiles := normalizeFiles(
		cloneTypeToNormalizeLevel(cloneType),
		levelNormalizers,
		filepaths,
		configuration,
	)

	err, compareFunc := getCompareFuncForLevel(cloneType)
	if err != nil {
		return err, []Clone{}
	}

	clones := detectClones(cloneType, normalizedFiles, compareFunc)

	return nil, clones
}

func findFilesInDir(directory string) ([]string, error) {
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

	return filepaths, err
}

func cloneTypeToNormalizeLevel(cloneType int) int {
	if cloneType > 2 {
		return 2
	}

	return cloneType
}

func getCompareFuncForLevel(level int) (error, func(a string, b string) []compare.Match) {
	if level == 1 || level == 2 {
		return nil, func(a string, b string) []compare.Match {
			return compare.FindExactMatches(a, b)
		}
	}

	if level == 3 {
		return nil, func(a string, b string) []compare.Match {
			return compare.FindLongestCommonSubsequence(a, b)
		}
	}

	return fmt.Errorf("no compare function found for level %d", level), nil
}

func detectClones(level int, normalizedFiles map[string]structs.File, compareFunc func(a string, b string) []compare.Match) []Clone {
	bar := createProgressbar(fmt.Sprintf("Detecting clones (level %d)", level), len(normalizedFiles)*len(normalizedFiles))

	var (
		clonesMutex sync.Mutex
		clones      []Clone
	)
	var clonesWg sync.WaitGroup

	pairs := buildFilePairs(normalizedFiles, bar)
	for _, pair := range pairs {
		clonesWg.Add(1)

		go func(pair filePair) {
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

func buildFilePairs(files map[string]structs.File, bar *progressbar.ProgressBar) map[string]filePair {
	pairs := make(map[string]filePair)

	for aPath, aFile := range files {
		for bPath, bFile := range files {
			bar.Add(1)

			if aPath == bPath {
				continue
			}

			pair := filePair{
				AFile: aFile,
				BFile: bFile,
			}.sort()

			hash := pair.hash()
			if _, ok := pairs[hash]; ok {
				continue
			}

			pairs[hash] = pair
		}
	}

	return pairs
}

func createProgressbar(description string, length int) *progressbar.ProgressBar {
	return progressbar.NewOptions(length,
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowDescriptionAtLineEnd(),
	)
}

func normalizeFiles(
	level int,
	levelNormalizers map[int][]config.Normalizer,
	filepaths []string,
	configuration config.Config,
) map[string]structs.File {
	cachePath := path.Join(configuration.Settings.CacheDirectory, "xdry-cache_level_"+strconv.Itoa(level)+".json")
	fileCache, err := cache.InitOrReadCache(cachePath)
	if err != nil {
		fmt.Println("Error reading cache")
	}

	var (
		normalizedFilesMutex sync.Mutex
		normalizedFiles      = make(map[string]structs.File)
	)
	var (
		cacheMutex       sync.Mutex
		mutexedFileCache = fileCache
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

	bar := createProgressbar(fmt.Sprintf("Normalizing files (level %d)", level), len(filepaths))

	const max = 12
	semaphore := make(chan struct{}, max)
	wg := &sync.WaitGroup{}

	for _, path := range filepaths {
		normalizedFile, err := cache.Get(*fileCache, path)
		if err == nil {
			normalizedFilesMutex.Lock()
			normalizedFiles[path] = *normalizedFile
			normalizedFilesMutex.Unlock()
			bar.Add(1)
			continue
		}

		semaphore <- struct{}{}
		wg.Add(1)

		go func(path string, mappedNormalizers map[string]config.Normalizer) {
			defer wg.Done()

			err, normalizedFile := normalize.Normalize(path, mappedNormalizers, cli.NewCommandExecutor())

			if err != nil {
				log.Println(err)
			}

			normalizedFilesMutex.Lock()
			normalizedFiles[path] = normalizedFile
			normalizedFilesMutex.Unlock()
			cacheMutex.Lock()
			cache.Store(*mutexedFileCache, normalizedFile)
			cacheMutex.Unlock()

			bar.Add(1)

			<-semaphore
		}(path, mappedNormalizers)
	}
	wg.Wait()

	return normalizedFiles
}
