package compare

import (
	"sort"
	"strings"
)

func FindLongestMatch(a string, b string) string {
	return longestCommonSubstring(a, b)
}

func longestCommonSubstring(a string, b string) string {
	var bothStrings []string
	if len(a) < len(b) {
		bothStrings = []string{a, b}
	} else {
		bothStrings = []string{b, a}
	}

	var longestCommonSubstring []string
	shorterString := bothStrings[0]
	longerString := bothStrings[1]

	charactersOfShorterString := strings.Split(shorterString, "")

	for len(charactersOfShorterString) > 0 {
		longestCommonSubstring = append([]string{""}, longestCommonSubstring...)
		for _, char := range charactersOfShorterString {
			if strstr(longerString, longestCommonSubstring[0]+char) == "" {
				break
			}
			longestCommonSubstring[0] += char
		}
		charactersOfShorterString = charactersOfShorterString[1:]
	}

	sort.Slice(longestCommonSubstring, func(i, j int) bool {
		return len(longestCommonSubstring[i]) > len(longestCommonSubstring[j])
	})

	if len(longestCommonSubstring) == 0 {
		return ""
	}

	return longestCommonSubstring[0]
}

func strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}

	return haystack[idx+len([]byte(needle))-1:]
}
