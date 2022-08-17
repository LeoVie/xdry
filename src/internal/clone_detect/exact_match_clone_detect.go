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
)

type Clone = struct {
	A       string
	B       string
	Matches []compare.Match
}

type Pair = struct {
	APath    string
	AContent string
	BPath    string
	BContent string
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

	normalizedFileContents := normalizeFiles(
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

	clones := detectClones(normalizedFileContents, compareFunc)

	return nil, clones
}

func detectClones(normalizedFileContents map[string]string, compareFunc func(a string, b string) []compare.Match) []Clone {
	pairs := make(map[string]Pair)

	for aPath, aContent := range normalizedFileContents {
		for bPath, bContent := range normalizedFileContents {
			if aPath == bPath {
				continue
			}

			firstPath, secondPath, firstContent, secondContent := orderPathsAndContents(aPath, aContent, bPath, bContent)

			hash := buildCloneHash(firstPath, secondPath)
			if _, ok := pairs[hash]; ok {
				continue
			}

			pairs[hash] = Pair{
				APath:    firstPath,
				AContent: firstContent,
				BPath:    secondPath,
				BContent: secondContent,
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

		go func(aPath string, aContent string, bPath string, bContent string) {
			defer clonesWg.Done()

			matches := compareFunc(aContent, bContent)

			if len(matches) == 0 {
				return
			}

			clonesMutex.Lock()
			clones = append(clones, Clone{
				A:       aPath,
				B:       bPath,
				Matches: matches,
			})
			clonesMutex.Unlock()
		}(pair.APath, pair.AContent, pair.BPath, pair.BContent)
	}

	//for aPath, aContent := range normalizedFileContents {
	//	for bPath, bContent := range normalizedFileContents {
	//		clonesWg.Add(1)
	//
	//		go func(aPath string, aContent string, bPath string, bContent string) {
	//			defer clonesWg.Done()
	//
	//			if aPath == bPath {
	//				return
	//			}
	//
	//			firstPath, secondPath, firstContent, secondContent := orderPathsAndContents(aPath, aContent, bPath, bContent)
	//
	//			hash := buildCloneHash(firstPath, secondPath)
	//
	//			fmt.Println(firstPath + ", " + secondPath)
	//
	//			matches := compareFunc(firstContent, secondContent)
	//
	//			if len(matches) == 0 {
	//				return
	//			}
	//
	//			clonesMutex.Lock()
	//			clones[hash] = Clone{
	//				A:       firstPath,
	//				B:       secondPath,
	//				Matches: matches,
	//			}
	//			clonesMutex.Unlock()
	//		}(aPath, aContent, bPath, bContent)
	//	}
	//}
	clonesWg.Wait()
	return clones
}

func orderPathsAndContents(aPath string, aContent string, bPath string, bContent string) (string, string, string, string) {
	if aPath < bPath {
		return aPath, bPath, aContent, bContent
	}

	return bPath, aPath, bContent, aContent
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

	const max = 12
	semaphore := make(chan struct{}, max)
	wg := &sync.WaitGroup{}

	for _, path := range filepaths {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(path string, normalizers []config.Normalizer) {
			defer wg.Done()

			err, normalizedFileContent := normalize.Normalize(path, normalizers, cli.NewCommandExecutor())

			if err != nil {
				fmt.Println(err)
			}

			normalizedFileContentsMutex.Lock()
			normalizedFileContents[path] = normalizedFileContent
			normalizedFileContentsMutex.Unlock()

			<-semaphore
		}(path, normalizers)
	}
	wg.Wait()

	return normalizedFileContents
}

func buildCloneHash(aPath string, bPath string) string {
	return aPath + "_" + bPath
}
