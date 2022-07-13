package compare

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Match struct {
	Content string
	IndexA  int
	IndexB  int
}

func FindMatches(a string, b string) []Match {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(a, b, false)

	matches := []Match{}
	indexA := 0
	indexB := 0
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			indexB += len(text)
		case diffmatchpatch.DiffDelete:
			indexA += len(text)
		case diffmatchpatch.DiffEqual:
			matches = append(matches, Match{
				Content: text,
				IndexA:  indexA,
				IndexB:  indexB,
			})
			indexA += len(text)
			indexB += len(text)
		}
	}

	return matches
}
