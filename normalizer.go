package normalizer

import (
	"strings"
)

func Normalize(s string) string {
	if len(s) == 0 || s == " " {
		return ""
	}

	s = strings.TrimSpace(s)
	s = normalizeSpaces(s)

	if len(s) == 0 {
		return ""
	}

	s = strings.ToLower(s)

	s = normalizeDotLikes(s)
	s = trimSpacesAroundDot(s)
	s = collapseDots(s)
	s = collapseSpaces(s)
	s = trimSpacesAroundDot(s)

	s = strings.ReplaceAll(s, " ", "-")

	s = collapseDashes(s)
	s = normalizeSepRuns(s)
	s = stripInvalidChars(s)
	s = normalizeLabelsKeepPunycode(s)
	s = trimDashes(s)
	s = trimDots(s)

	return s
}
