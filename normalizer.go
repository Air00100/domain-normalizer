package normalizer

import (
	"errors"
	"strings"

	"golang.org/x/net/idna"
)

var (
	ErrEmptyInput   = errors.New("empty input")
	ErrTooLong      = errors.New("domain too long")
	ErrInvalidLabel = errors.New("invalid label")
	ErrIDNA         = errors.New("idna conversion failed")
)

type Domain struct {
	Raw          string
	Normalized   string
	ASCII        string
	SubDomain    string
	Sld          string
	Tld          string
	Registerable string
	Labels       []string
	ASCIILabels  []string
	IsIDN        bool
	HasPunycode  bool
	Icann        bool
}

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

func Parse(s string) (Domain, error) {
	domain := Domain{Raw: s}
	normalized := Normalize(s)

	if len(normalized) == 0 {
		return domain, ErrEmptyInput
	}

	domain.Normalized = normalized
	domain.Labels = strings.Split(normalized, ".")

	ascii, err := idna.Lookup.ToASCII(normalized)
	if err != nil {
		return domain, errors.Join(ErrIDNA, err)
	}

	domain.ASCII = ascii
	domain.ASCIILabels = strings.Split(ascii, ".")
	domain.IsIDN = ascii != normalized
	domain.HasPunycode = false

	for _, lbl := range domain.ASCIILabels {
		if strings.HasPrefix(lbl, "xn--") {
			domain.HasPunycode = true
			break
		}
	}

	if err = validateASCII(ascii, domain.ASCIILabels); err != nil {
		return domain, err
	}

	handleParts(ascii, &domain)

	return domain, nil
}
