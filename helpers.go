package normalizer

import (
	"regexp"
	"strings"
)

var (
	dots            = regexp.MustCompile(`\.+`)
	spaces          = regexp.MustCompile(`\s+`)
	allSpaces       = regexp.MustCompile(`[\s\p{Z}]+`)
	dash            = regexp.MustCompile(`-+`)
	invalidSymbols  = regexp.MustCompile(`[^\p{L}\p{N}\-.]+`)
	spacesAroundDot = regexp.MustCompile(`\s*\.\s*`)
	sepRuns         = regexp.MustCompile(`[-.]+`)
	dotLikes        = regexp.MustCompile(`[\x{3002}\x{FF0E}\x{FF61}]`)
)

func collapseDots(s string) string {
	return dots.ReplaceAllString(s, ".")
}

func collapseSpaces(s string) string {
	return spaces.ReplaceAllString(s, " ")
}

func collapseDashes(s string) string {
	return dash.ReplaceAllString(s, "-")
}

func stripInvalidChars(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "\\", " ")
	return invalidSymbols.ReplaceAllString(s, "")
}

func trimDashes(s string) string {
	s = strings.TrimLeft(s, "-")
	s = strings.TrimRight(s, "-")
	return s
}

func trimDots(s string) string {
	s = strings.TrimLeft(s, ".")
	s = strings.TrimRight(s, ".")
	return s
}

func trimSpacesAroundDot(s string) string {
	return spacesAroundDot.ReplaceAllString(s, ".")
}

func normalizeSpaces(s string) string {
	return allSpaces.ReplaceAllString(s, " ")
}

func normalizeSepRuns(s string) string {
	return sepRuns.ReplaceAllStringFunc(s, func(run string) string {
		hasDot := strings.Contains(run, ".")
		hasDash := strings.Contains(run, "-")
		if hasDot && hasDash {
			if run == "-." || run == ".-" {
				return "."
			}
			return "-"
		}
		if hasDot {
			return "."
		}
		return "-"
	})
}

func normalizeDotLikes(s string) string {
	return dotLikes.ReplaceAllString(s, ".")
}

func normalizeLabelsKeepPunycode(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, ".")
	for i, p := range parts {
		if p == "" {
			continue
		}

		pTrim := strings.Trim(p, "-")

		if strings.HasPrefix(pTrim, "xn") {
			rest := strings.TrimPrefix(pTrim, "xn")
			rest = strings.TrimLeft(rest, "-")
			p = "xn--" + rest
		} else {
			p = pTrim
		}

		parts[i] = p
	}

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}

	return strings.Join(out, ".")
}
