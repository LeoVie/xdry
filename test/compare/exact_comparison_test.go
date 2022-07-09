package compare

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"x-dry-go/internal/compare"
)

func TestFindLongestMatch(t *testing.T) {
	datasets := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"#1", "ABCDEFGHIJ", "KLMNOPQRST", ""},
		{"#2", "ABCDEFGHIJ", "ABCDEFGHIJ", "ABCDEFGHIJ"},
		{"#3", "ABCDEFGHIJ", "ABCDE12345", "ABCDE"},
		{"#4", "ABCDEFGHIJ", "ABCDEF2345", "ABCDEF"},
		{"#5", "ABCDEFGHIJ", "_ABCDE1234", "ABCDE"},
		{"#6", "_ABCDE1234", "ABCDEFGHIJ", "ABCDE"},
		{"#7", "ABCDEABCDEFGHIJ", "ABCDEF___ABCDEFGHIJ", "ABCDEFGHIJ"},
	}

	for _, dataset := range datasets {
		t.Run(dataset.name, func(t *testing.T) {
			actual := compare.FindLongestMatch(dataset.a, dataset.b)

			assert.Equal(t, dataset.want, actual)
		})
	}
}
