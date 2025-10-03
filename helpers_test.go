package normalizer

import "testing"

func Test_collapseDots(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "a.b.c", "a.b.c"},
		{"...", "...", "."},
		{"test.com", "test....com", "test.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collapseDots(tt.s); got != tt.want {
				t.Errorf("collapseDots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collapseSpaces(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "a    b  c", "a b c"},
		{"spaces", "          ", " "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collapseSpaces(tt.s); got != tt.want {
				t.Errorf("collapseSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collapseDashes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "a---b--c", "a-b-c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collapseDashes(tt.s); got != tt.want {
				t.Errorf("collapseDashes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripInvalidChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "a.b.c", "a.b.c"},
		{"has invalid chars", "^&some_-domain.com!,", "some-domain.com"},
		{"numbers", "123.456", "123.456"},
		{"unicode letters", "пример.рф", "пример.рф"},
		{"unicode with invalid", "тест$$$.рф", "тест.рф"},
		{"emoji filtered", "😀example.com", "example.com"},
		{"spaces removed", "foo bar com", "foobarcom"},
		{"non-breaking hyphen", "go\u2011lang.org", "golang.org"},
		{"en/em dash removed", "go\u2013long\u2014domain.com", "golongdomain.com"},
		{"punctuation storm", "ex,am;ple:.co!m", "example.com"},
		{"multiple symbols", "@@@exa!!!mple###.com$$$", "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripInvalidChars(tt.s); got != tt.want {
				t.Errorf("stripInvalidChars() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_trimDashes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"dash only", "---", ""},
		{"simple", "--a.b.c-", "a.b.c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimDashes(tt.s); got != tt.want {
				t.Errorf("trimDashes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimDots(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"dots only", "...", ""},
		{"simple", "..a.b.c.", "a.b.c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimDots(tt.s); got != tt.want {
				t.Errorf("trimDots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimSpaceArounDot(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"dot only", " . ", "."},
		{"simple", " . a . b . c . ", ".a.b.c."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimSpacesAroundDot(tt.s); got != tt.want {
				t.Errorf("trimSpacesAroundDot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeSpaces(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"ascii space", "a b  c", "a b c"},
		{"tabs", "a\t\tb\tc", "a b c"},
		{"newlines", "a\nb\r\nc", "a b c"},
		{"mixed ws", " a \t b \n c  ", " a b c "},
		{"nbsp", "a\u00A0\u00A0b\u00A0c", "a b c"},
		{"thin/figure spaces", "a\u2009b\u2007c", "a b c"},
		{"em space", "a\u2003\u2003b", "a b"},
		{"nnbsp", "a\u202Fb\u202Fc", "a b c"},
		{"zero width (kept, not Z)", "a\u200Db", "a\u200Db"},
		{"leading/trailing", "\u2003  a  \u2003", " a "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSpaces(tt.in); got != tt.want {
				t.Errorf("normalizeSpaces() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_normalizeSepRuns(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single dot", ".", "."},
		{"single dash", "-", "-"},
		{"dots only", "....", "."},
		{"dashes only", "----", "-"},
		{"dash then dot (pair)", "-.", "."},
		{"dot then dash (pair)", ".-", "."},
		{"mixed short", "-.-", "-"},
		{"mixed long", "--.-.-", "-"},
		{"dot-run between words", "foo...bar", "foo.bar"},
		{"dash-run between words", "foo---bar", "foo-bar"},
		{"mixed cluster between words", "foo--.-.-bar", "foo-bar"},
		{"leading mixed", "-.-.-foo", "-foo"},
		{"trailing mixed", "bar-.-.-", "bar-"},
		{"around dot kept as dot", "foo-.bar.-baz", "foo.bar.baz"},
		{"only mixed becomes dash", "-.-.-", "-"},
		{"dot and dash mess full", "..foo--.-.-bar..", ".foo-bar."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSepRuns(tt.in); got != tt.want {
				t.Errorf("normalizeSepRuns() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_normalizeDotLikes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain dot", ".", "."},
		{"ideographic full stop", "\u3002", "."},
		{"fullwidth full stop", "\uFF0E", "."},
		{"halfwidth middle dot", "\uFF61", "."},
		{"ascii unchanged", "example.com", "example.com"},
		{"ideographic dot in domain", "example\u3002com", "example.com"},
		{"fullwidth dot in domain", "example\uFF0Ecom", "example.com"},
		{"halfwidth dot in domain", "example\uFF61com", "example.com"},
		{"mixed dotlikes", "a\u3002b\uFF0Ec\uFF61d", "a.b.c.d"},
		{"no dotlikes", "foo-bar", "foo-bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeDotLikes(tt.in); got != tt.want {
				t.Errorf("normalizeDotLikes() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_normalizeLabelsKeepPunycode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"simple domain", "example.com", "example.com"},
		{"preserve punycode", "xn--d1acufc.xn--p1ai", "xn--d1acufc.xn--p1ai"},
		{"fix single-dash punycode", "xn-d1acufc.xn-p1ai", "xn--d1acufc.xn--p1ai"},
		{"collapse extra dashes after xn", "xn----d1acufc", "xn--d1acufc"},
		{"trim label edges", "-foo.-bar-", "foo.bar"},
		{"remove empty labels", "a..b...c", "a.b.c"},
		{"mix: puny + trims", "-xn---test-.", "xn--test"},

		{"puny with edges", "-xn-d1acufc-", "xn--d1acufc"},
		{"puny collapsed head", "xn-----abc", "xn--abc"},
		{"puny just prefix", "xn--", "xn--"},

		{"only dashes label", "----", ""},
		{"inner double dashes kept", "foo--bar.com", "foo--bar.com"},
		{"keep multiple labels", "sub.-mid-.tail.", "sub.mid.tail"},
		{"leading/trailing dots", ".foo.bar.", "foo.bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeLabelsKeepPunycode(tt.in); got != tt.want {
				t.Errorf("normalizeLabelsKeepPunycode() = %q, want %q", got, tt.want)
			}
		})
	}
}
