package compare

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"x-dry-go/internal/compare"
)

func TestFindExactMatches(t *testing.T) {
	datasets := []struct {
		name string
		a    string
		b    string
		want []compare.Match
	}{
		{"#1",
			"ABCDEFGHIJ",
			"KLMNOPQRST",
			[]compare.Match{},
		},
		{"#2",
			"ABCDEFGHIJ",
			"ABCDEFGHIJ",
			[]compare.Match{
				{
					Content: "ABCDEFGHIJ",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{"#3",
			"ABCDEFGHIJ",
			"ABCDE12345",
			[]compare.Match{
				{
					Content: "ABCDE",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{
			"#4",
			"ABCDEFGHIJ",
			"ABCDEF2345",
			[]compare.Match{
				{
					Content: "ABCDEF",
					IndexA:  0,
					IndexB:  0,
				},
			},
		},
		{
			"#5",
			"ABCDEFGHIJ",
			"_ABCDE1234",
			[]compare.Match{
				{
					Content: "ABCDE",
					IndexA:  0,
					IndexB:  1,
				},
			},
		},
		{
			"#6",
			"_ABCDE1234",
			"ABCDEFGHIJ",
			[]compare.Match{
				{
					Content: "ABCDE",
					IndexA:  1,
					IndexB:  0,
				},
			},
		},
		{
			"#7",
			"ABCDEABCDEFGHIJ",
			"ABCDEF___ABCDEFGHIJ",
			[]compare.Match{
				{
					Content: "ABCDE",
					IndexA:  0,
					IndexB:  0,
				},
				{
					Content: "ABCDEFGHIJ",
					IndexA:  5,
					IndexB:  9,
				},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.name, func(t *testing.T) {
			actual := compare.FindExactMatches(dataset.a, dataset.b)

			assert.Equal(t, dataset.want, actual)
		})
	}
}

func TestFindLongestCommonSubsequence(t *testing.T) {
	datasets := []struct {
		name string
		a    string
		b    string
		want []compare.Match
	}{
		{"#1",
			"ABCDEFGHIJ",
			"KLMNOPQRST",
			[]compare.Match{},
		},
		{"#2",
			"ABCDEFGHIJ",
			"ABCDEFGHIJ",
			[]compare.Match{
				{
					Content: "ABCDEFGHIJ",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{"#3",
			"ABCDEFGHIJ",
			"ABCDE12345",
			[]compare.Match{
				{
					Content: "ABCDE",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{
			"#4",
			"ABCDEABCDEFGHIJ",
			"ABCDEF___ABCDEFGHIJ",
			[]compare.Match{
				{
					Content: "ABCDEABCDEFGHIJ",
					IndexA:  0,
					IndexB:  0,
				},
			},
		},
		{
			"#5",
			"ABC",
			"CBA",
			[]compare.Match{
				{
					Content: "C",
					IndexA:  2,
					IndexB:  0,
				},
			},
		},
		{
			"#6",
			"AAABCDDDEE",
			"ABCXDDE",
			[]compare.Match{
				{
					Content: "ABCDDE",
					IndexA:  2,
					IndexB:  0,
				},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.name, func(t *testing.T) {
			actual := compare.FindLongestCommonSubsequence(dataset.a, dataset.b)

			assert.Equal(t, dataset.want, actual)
		})
	}
}
