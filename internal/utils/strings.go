package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func NormalizeCityName(city string) string {
	t := norm.NFD.String(city)
	var sb strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	normalized := strings.ReplaceAll(sb.String(), " ", "")
	return normalized
}
