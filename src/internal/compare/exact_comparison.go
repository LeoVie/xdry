package compare

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

type Match struct {
	Content string
	IndexA  int
	IndexB  int
}

func FindExactMatches(a string, b string) []Match {
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

func FindLongestCommonSubsequence(a string, b string) []Match {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(a, b, false)

	indexA := 0
	indexB := 0
	matchBuffer := strings.Builder{}
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			if matchBuffer.Len() == 0 {
				indexB += len(text)
			}
		case diffmatchpatch.DiffDelete:
			if matchBuffer.Len() == 0 {
				indexA += len(text)
			}
		case diffmatchpatch.DiffEqual:
			matchBuffer.WriteString(text)
		}
	}

	if matchBuffer.Len() == 0 {
		return []Match{}
	}

	return []Match{
		{
			Content: matchBuffer.String(),
			IndexA:  indexA,
			IndexB:  indexB,
		},
	}
}
