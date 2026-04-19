package helper

import "strings"

func UpperAndTrim(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

var preservedWords = map[string]bool{
	"MAX": true, "MIN": true,
	"XS": true, "S": true, "M": true, "L": true, "XL": true, "XXL": true, "XXXL": true,
	"LED": true, "USB": true, "RAM": true, "ROM": true,
}

func TitleCase(s string) string {
	words := strings.Fields(strings.TrimSpace(s))
	for i, word := range words {
		upper := strings.ToUpper(word)
		if preservedWords[upper] {
			words[i] = upper
			continue
		}
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func LowerAndTrim(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
