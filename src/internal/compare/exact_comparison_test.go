package compare

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestFindExactMatches(t *testing.T) {
	g := NewGomegaWithT(t)

	datasets := []struct {
		name string
		a    string
		b    string
		want []Match
	}{
		{"#1",
			"ABCDEFGHIJ",
			"KLMNOPQRST",
			[]Match{},
		},
		{"#2",
			"ABCDEFGHIJ",
			"ABCDEFGHIJ",
			[]Match{
				{
					Content: "ABCDEFGHIJ",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{"#3",
			"ABCDEFGHIJ",
			"ABCDE12345",
			[]Match{
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
			[]Match{
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
			[]Match{
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
			[]Match{
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
			[]Match{
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
			actual := FindExactMatches(dataset.a, dataset.b)

			g.Expect(actual).To(Equal(dataset.want))
		})
	}
}

func TestFindLongestCommonSubsequence(t *testing.T) {
	g := NewGomegaWithT(t)

	datasets := []struct {
		name string
		a    string
		b    string
		want []Match
	}{
		{"#1",
			"ABCDEFGHIJ",
			"KLMNOPQRST",
			[]Match{},
		},
		{"#2",
			"ABCDEFGHIJ",
			"ABCDEFGHIJ",
			[]Match{
				{
					Content: "ABCDEFGHIJ",
					IndexA:  0,
					IndexB:  0,
				},
			}},
		{"#3",
			"ABCDEFGHIJ",
			"ABCDE12345",
			[]Match{
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
			[]Match{
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
			[]Match{
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
			[]Match{
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
			actual := FindLongestCommonSubsequence(dataset.a, dataset.b)

			g.Expect(actual).To(Equal(dataset.want))
		})
	}
}
