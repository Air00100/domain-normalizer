package normalizer

import (
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"single space", " ", ""},
		{"many spaces", "           ", ""},
		{"trim spaces", "   example.com   ", "example.com"},
		{"uppercase", "ExAmPlE.CoM", "example.com"},
		{"multiple dots", "foo..bar...com", "foo.bar.com"},
		{"leading dot", ".example.com", "example.com"},
		{"trailing dot", "example.com.", "example.com"},
		{"spaces inside", "exa mple com", "exa-mple-com"},
		{"multiple spaces", "foo   bar com", "foo-bar-com"},
		{"mixed dots and spaces", "foo .. bar . com", "foo.bar.com"},
		{"multiple dashes", "foo---bar--com", "foo-bar-com"},
		{"leading dashes", "---example.com", "example.com"},
		{"trailing dashes", "example.com---", "example.com"},
		{"symbols filtered", "exa$mple!com?", "examplecom"},
		{"unicode letters", "пример.рф", "пример.рф"},
		{"unicode with spaces", "  при мер   . рф ", "при-мер.рф"},
		{"mix of all", " SomE...DOMa in..com!!! ", "some.doma-in.com"},

		{"dash around dot", "foo-.bar.-baz", "foo.bar.baz"},
		{"numbers only", "123.456", "123.456"},
		{"unicode chinese", "测试.公司", "测试.公司"},
		{"unicode arabic", "مثال.مصر", "مثال.مصر"},
		{"unicode hebrew", "דוגמה.ישראל", "דוגמה.ישראל"},
		{"unicode chinese2", "\u6D4B\u8BD5.\u516C\u53F8", "\u6D4B\u8BD5.\u516C\u53F8"},
		{
			"unicode arabic2",
			"\u0645\u062B\u0627\u0644.\u0645\u0635\u0631",
			"\u0645\u062B\u0627\u0644.\u0645\u0635\u0631",
		},
		{
			"unicode hebrew2",
			"\u05D3\u05D5\u05D2\u05DE\u05D4.\u05D9\u05E9\u05E8\u05D0\u05DC",
			"\u05D3\u05D5\u05D2\u05DE\u05D4.\u05D9\u05E9\u05E8\u05D0\u05DC",
		},
		{
			"unicode russian",
			"\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
			"\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
		},
		{"unicode with dashes", "пример---сайт.рф", "пример-сайт.рф"},
		{"unicode with mixed symbols", "тест$$$.рф", "тест.рф"},
		{"emoji filtered", "😀example.com", "example.com"},
		{"mixed digits and unicode", "дом123.рф", "дом123.рф"},
		{
			"long label with dashes",
			"foo--------------------------------------------bar.com",
			"foo-bar.com",
		},
		{"dot and dash mess", "..foo--.-.-bar..", "foo-bar"},
		{"only symbols", "!!!@@@###", ""},
		{"only dots", "...", ""},

		{"tabs and newlines", "\tfoo \n bar \r\n com\t", "foo-bar-com"},
		{"nbsp spaces", "foo\u00A0\u00A0bar\u00A0com", "foo-bar-com"},
		{"underscores removed", "exa_mple..com", "example.com"},
		{"brackets stripped", "[example].com", "example.com"},
		{"quotes stripped", `"exa'mple".com`, "example.com"},
		{"slash/backslash stripped", `exa/mpl\e.com`, "example.com"},
		{"non-ascii dashes removed", "ex—am–ple.com", "example.com"},
		{"nbhyphen removed", "go\u2011lang.org", "golang.org"},
		{"zero-width removed", "ex\u200Dample.\u200Ccom", "example.com"},
		{"spaces-around-dot", "foo .  bar .   baz", "foo.bar.baz"},
		{"mixed run around dot", "foo-.-.-bar", "foo-bar"},
		{"dash touching dots", "-.foo.-", "foo"},
		{"ip-like stays", "192.168.0.1", "192.168.0.1"},
		{"subdomains deep", "a...b..c.d....e.com", "a.b.c.d.e.com"},
		{"leading/trailing unicode spaces", " \u2003пример.\u2007сайт.\u202Frf ", "пример.сайт.rf"},
		{"ending dot with junk", "example.com.!!!", "example.com"},
		{"label-only dashes", "---", ""},
		{"mixed separators cluster", "foo.-.--..-.-.bar", "foo-bar"},
		{"punctuation storm", "ex,am;ple:.co!m", "example.com"},
		{"mixed scripts safe", "пример-example.ком", "пример-example.ком"},
		{"punycode preserved", "xn--d1acufc.xn--p1ai", "xn--d1acufc.xn--p1ai"},
		{"punycode with extra dashes", "xn----d1acufc.xn---p1ai", "xn--d1acufc.xn--p1ai"},
		{"punycode", "xn--d1acufc.xn--p1ai", "xn--d1acufc.xn--p1ai"},
		{"dot likes to dot", "example。com", "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.s); got != tt.want {
				t.Errorf("Normalize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type want struct {
		Normalized  string
		ASCII       string
		TLD         string
		Registrable string
		SLD         string
		Subdomain   string
		IsIDN       bool
		HasPuny     bool
	}

	tests := []struct {
		name    string
		in      string
		want    want
		wantErr bool
	}{
		{
			name: "simple domain",
			in:   "example.com",
			want: want{
				Normalized:  "example.com",
				ASCII:       "example.com",
				TLD:         "com",
				Registrable: "example.com",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "with subdomains",
			in:   "a.b.example.com",
			want: want{
				Normalized:  "a.b.example.com",
				ASCII:       "a.b.example.com",
				TLD:         "com",
				Registrable: "example.com",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "publicsuffix co.uk",
			in:   "a.b.example.co.uk",
			want: want{
				Normalized:  "a.b.example.co.uk",
				ASCII:       "a.b.example.co.uk",
				TLD:         "co.uk",
				Registrable: "example.co.uk",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "unicode IDN (пример.рф)",
			in:   "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
			want: want{
				Normalized:  "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				ASCII:       "xn--e1afmkfd.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       true,
				HasPuny:     true,
			},
		},
		{
			name: "punycode input (пример.рф)",
			in:   "xn--e1afmkfd.xn--p1ai",
			want: want{
				Normalized:  "xn--e1afmkfd.xn--p1ai",
				ASCII:       "xn--e1afmkfd.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "puny in TLD",
			in:   "example.xn--p1ai",
			want: want{
				Normalized:  "example.xn--p1ai",
				ASCII:       "example.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "example.\u0440\u0444",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "messy input normalized",
			in:   "  SomE...DOMa in..com!!!  ",
			want: want{
				Normalized:  "some.doma-in.com",
				ASCII:       "some.doma-in.com",
				TLD:         "com",
				Registrable: "doma-in.com",
				SLD:         "doma-in",
				Subdomain:   "some",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name:    "empty after normalize",
			in:      "   ",
			wantErr: true,
		},
		{
			name:    "label too long (>63)",
			in:      strings.Repeat("a", 64) + ".com",
			wantErr: true,
		},
		{
			name: "domain length = 253 ok",
			in: func() string {
				l63a := strings.Repeat("a", 63)
				l63b := strings.Repeat("b", 63)
				l63c := strings.Repeat("c", 63)
				l61d := strings.Repeat("d", 61)
				return strings.Join([]string{l63a, l63b, l63c, l61d}, ".")
			}(),
			want: want{
				TLD:         strings.Repeat("d", 61),
				Registrable: strings.Repeat("c", 63) + "." + strings.Repeat("d", 61),
				SLD:         strings.Repeat("c", 63),
			},
		},
		{
			name: "domain length > 253",
			in: func() string {
				l63a := strings.Repeat("a", 63)
				l63b := strings.Repeat("b", 63)
				l63c := strings.Repeat("c", 63)
				l62d := strings.Repeat("d", 62)
				return strings.Join([]string{l63a, l63b, l63c, l62d}, ".")
			}(),
			wantErr: true,
		},
		{
			name:    "leading hyphen in label",
			in:      "-abc.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) expected error, got nil", tt.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.in, err)
			}

			if tt.want.Normalized != "" && got.Normalized != tt.want.Normalized {
				t.Errorf("Normalized = %q, want %q", got.Normalized, tt.want.Normalized)
			}
			if tt.want.ASCII != "" && got.ASCII != tt.want.ASCII {
				t.Errorf("ASCII = %q, want %q", got.ASCII, tt.want.ASCII)
			}
			if tt.want.TLD != "" && got.Tld != tt.want.TLD {
				t.Errorf("TLD = %q, want %q", got.Tld, tt.want.TLD)
			}
			if tt.want.Registrable != "" && got.Registerable != tt.want.Registrable {
				t.Errorf("Registerable = %q, want %q", got.Registerable, tt.want.Registrable)
			}
			if tt.want.SLD != "" && got.Sld != tt.want.SLD {
				t.Errorf("SLD = %q, want %q", got.Sld, tt.want.SLD)
			}
			if tt.want.Subdomain != "" && got.SubDomain != tt.want.Subdomain {
				t.Errorf("Subdomain = %q, want %q", got.SubDomain, tt.want.Subdomain)
			}
			if got.IsIDN != tt.want.IsIDN {
				t.Errorf("IsIDN = %v, want %v", got.IsIDN, tt.want.IsIDN)
			}
			if got.HasPunycode != tt.want.HasPuny {
				t.Errorf("HasPunycode = %v, want %v", got.HasPunycode, tt.want.HasPuny)
			}
		})
	}
}

func TestParse_ValidDomains(t *testing.T) {
	type want struct {
		Normalized  string
		ASCII       string
		TLD         string
		Registrable string
		SLD         string
		Subdomain   string
		IsIDN       bool
		HasPuny     bool
	}

	tests := []struct {
		name string
		in   string
		want want
	}{
		{
			name: "simple com",
			in:   "example.com",
			want: want{
				Normalized:  "example.com",
				ASCII:       "example.com",
				TLD:         "com",
				Registrable: "example.com",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "www subdomain",
			in:   "www.example.com",
			want: want{
				Normalized:  "www.example.com",
				ASCII:       "www.example.com",
				TLD:         "com",
				Registrable: "example.com",
				SLD:         "example",
				Subdomain:   "www",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "deep subdomains",
			in:   "a.b.c.example.com",
			want: want{
				Normalized:  "a.b.c.example.com",
				ASCII:       "a.b.c.example.com",
				TLD:         "com",
				Registrable: "example.com",
				SLD:         "example",
				Subdomain:   "a.b.c",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "publicsuffix co.uk",
			in:   "a.b.example.co.uk",
			want: want{
				Normalized:  "a.b.example.co.uk",
				ASCII:       "a.b.example.co.uk",
				TLD:         "co.uk",
				Registrable: "example.co.uk",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "publicsuffix gov.au",
			in:   "example.gov.au",
			want: want{
				Normalized:  "example.gov.au",
				ASCII:       "example.gov.au",
				TLD:         "gov.au",
				Registrable: "example.gov.au",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "hyphen inside",
			in:   "my-domain.io",
			want: want{
				Normalized:  "my-domain.io",
				ASCII:       "my-domain.io",
				TLD:         "io",
				Registrable: "my-domain.io",
				SLD:         "my-domain",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "numeric and letters",
			in:   "123-456.net",
			want: want{
				Normalized:  "123-456.net",
				ASCII:       "123-456.net",
				TLD:         "net",
				Registrable: "123-456.net",
				SLD:         "123-456",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "unicode IDN (пример.рф)",
			in:   "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
			want: want{
				Normalized:  "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				ASCII:       "xn--e1afmkfd.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       true,
				HasPuny:     true,
			},
		},
		{
			name: "punycode input (пример.рф)",
			in:   "xn--e1afmkfd.xn--p1ai",
			want: want{
				Normalized:  "xn--e1afmkfd.xn--p1ai",
				ASCII:       "xn--e1afmkfd.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "puny in TLD",
			in:   "example.xn--p1ai",
			want: want{
				Normalized:  "example.xn--p1ai",
				ASCII:       "example.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "example.\u0440\u0444",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "unicode SLD with subdomains (пример.рф)",
			in:   "a.b.\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
			want: want{
				Normalized:  "a.b.\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				ASCII:       "a.b.xn--e1afmkfd.xn--p1ai",
				TLD:         "\u0440\u0444",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.\u0440\u0444",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "a.b",
				IsIDN:       true,
				HasPuny:     true,
			},
		},
		{
			name: "co.uk simple",
			in:   "example.co.uk",
			want: want{
				Normalized:  "example.co.uk",
				ASCII:       "example.co.uk",
				TLD:         "co.uk",
				Registrable: "example.co.uk",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.uk with subdomains",
			in:   "a.b.example.co.uk",
			want: want{
				Normalized:  "a.b.example.co.uk",
				ASCII:       "a.b.example.co.uk",
				TLD:         "co.uk",
				Registrable: "example.co.uk",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.uk IDN SLD (unicode)",
			in:   "\u043F\u0440\u0438\u043C\u0435\u0440.co.uk",
			want: want{
				Normalized:  "\u043F\u0440\u0438\u043C\u0435\u0440.co.uk",
				ASCII:       "xn--e1afmkfd.co.uk",
				TLD:         "co.uk",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.co.uk",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       true,
				HasPuny:     true,
			},
		},
		{
			name: "co.uk IDN SLD (puny input)",
			in:   "xn--e1afmkfd.co.uk",
			want: want{
				Normalized:  "xn--e1afmkfd.co.uk",
				ASCII:       "xn--e1afmkfd.co.uk",
				TLD:         "co.uk",
				Registrable: "\u043F\u0440\u0438\u043C\u0435\u0440.co.uk",
				SLD:         "\u043F\u0440\u0438\u043C\u0435\u0440",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "gov.uk simple",
			in:   "example.gov.uk",
			want: want{
				Normalized:  "example.gov.uk",
				ASCII:       "example.gov.uk",
				TLD:         "gov.uk",
				Registrable: "example.gov.uk",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "gov.uk with subdomain",
			in:   "sub.example.gov.uk",
			want: want{
				Normalized:  "sub.example.gov.uk",
				ASCII:       "sub.example.gov.uk",
				TLD:         "gov.uk",
				Registrable: "example.gov.uk",
				SLD:         "example",
				Subdomain:   "sub",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "ac.uk simple",
			in:   "example.ac.uk",
			want: want{
				Normalized:  "example.ac.uk",
				ASCII:       "example.ac.uk",
				TLD:         "ac.uk",
				Registrable: "example.ac.uk",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.uk deep chain",
			in:   "x.y.z.example.co.uk",
			want: want{
				Normalized:  "x.y.z.example.co.uk",
				ASCII:       "x.y.z.example.co.uk",
				TLD:         "co.uk",
				Registrable: "example.co.uk",
				SLD:         "example",
				Subdomain:   "x.y.z",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.jp simple",
			in:   "example.co.jp",
			want: want{
				Normalized:  "example.co.jp",
				ASCII:       "example.co.jp",
				TLD:         "co.jp",
				Registrable: "example.co.jp",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.jp with subdomains",
			in:   "a.b.example.co.jp",
			want: want{
				Normalized:  "a.b.example.co.jp",
				ASCII:       "a.b.example.co.jp",
				TLD:         "co.jp",
				Registrable: "example.co.jp",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.jp IDN SLD (unicode)",
			in:   "\u4f8b\u3048.co.jp",
			want: want{
				Normalized:  "\u4f8b\u3048.co.jp",
				ASCII:       "xn--r8jz45g.co.jp",
				TLD:         "co.jp",
				Registrable: "\u4f8b\u3048.co.jp",
				SLD:         "\u4f8b\u3048",
				Subdomain:   "",
				IsIDN:       true,
				HasPuny:     true,
			},
		},
		{
			name: "co.jp IDN SLD (puny input)",
			in:   "xn--r8jz45g.co.jp",
			want: want{
				Normalized:  "xn--r8jz45g.co.jp",
				ASCII:       "xn--r8jz45g.co.jp",
				TLD:         "co.jp",
				Registrable: "\u4f8b\u3048.co.jp",
				SLD:         "\u4f8b\u3048",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
		{
			name: "ne.jp",
			in:   "example.ne.jp",
			want: want{
				Normalized:  "example.ne.jp",
				ASCII:       "example.ne.jp",
				TLD:         "ne.jp",
				Registrable: "example.ne.jp",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "or.jp with sub",
			in:   "sub.example.or.jp",
			want: want{
				Normalized:  "sub.example.or.jp",
				ASCII:       "sub.example.or.jp",
				TLD:         "or.jp",
				Registrable: "example.or.jp",
				SLD:         "example",
				Subdomain:   "sub",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.au",
			in:   "example.com.au",
			want: want{
				Normalized:  "example.com.au",
				ASCII:       "example.com.au",
				TLD:         "com.au",
				Registrable: "example.com.au",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.au with subdomains",
			in:   "a.b.example.com.au",
			want: want{
				Normalized:  "a.b.example.com.au",
				ASCII:       "a.b.example.com.au",
				TLD:         "com.au",
				Registrable: "example.com.au",
				SLD:         "example",
				Subdomain:   "a.b",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.nz",
			in:   "example.co.nz",
			want: want{
				Normalized:  "example.co.nz",
				ASCII:       "example.co.nz",
				TLD:         "co.nz",
				Registrable: "example.co.nz",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.mx",
			in:   "example.com.mx",
			want: want{
				Normalized:  "example.com.mx",
				ASCII:       "example.com.mx",
				TLD:         "com.mx",
				Registrable: "example.com.mx",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "gob.mx (government MX)",
			in:   "example.gob.mx",
			want: want{
				Normalized:  "example.gob.mx",
				ASCII:       "example.gob.mx",
				TLD:         "gob.mx",
				Registrable: "example.gob.mx",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.br",
			in:   "example.com.br",
			want: want{
				Normalized:  "example.com.br",
				ASCII:       "example.com.br",
				TLD:         "com.br",
				Registrable: "example.com.br",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.za",
			in:   "example.co.za",
			want: want{
				Normalized:  "example.co.za",
				ASCII:       "example.co.za",
				TLD:         "co.za",
				Registrable: "example.co.za",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.kr",
			in:   "example.co.kr",
			want: want{
				Normalized:  "example.co.kr",
				ASCII:       "example.co.kr",
				TLD:         "co.kr",
				Registrable: "example.co.kr",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.tr",
			in:   "example.com.tr",
			want: want{
				Normalized:  "example.com.tr",
				ASCII:       "example.com.tr",
				TLD:         "com.tr",
				Registrable: "example.com.tr",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.id",
			in:   "example.co.id",
			want: want{
				Normalized:  "example.co.id",
				ASCII:       "example.co.id",
				TLD:         "co.id",
				Registrable: "example.co.id",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.cn",
			in:   "example.com.cn",
			want: want{
				Normalized:  "example.com.cn",
				ASCII:       "example.com.cn",
				TLD:         "com.cn",
				Registrable: "example.com.cn",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.tw",
			in:   "example.com.tw",
			want: want{
				Normalized:  "example.com.tw",
				ASCII:       "example.com.tw",
				TLD:         "com.tw",
				Registrable: "example.com.tw",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "com.hk",
			in:   "example.com.hk",
			want: want{
				Normalized:  "example.com.hk",
				ASCII:       "example.com.hk",
				TLD:         "com.hk",
				Registrable: "example.com.hk",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.in",
			in:   "example.co.in",
			want: want{
				Normalized:  "example.co.in",
				ASCII:       "example.co.in",
				TLD:         "co.in",
				Registrable: "example.co.in",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "gov.in",
			in:   "example.gov.in",
			want: want{
				Normalized:  "example.gov.in",
				ASCII:       "example.gov.in",
				TLD:         "gov.in",
				Registrable: "example.gov.in",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "co.il",
			in:   "example.co.il",
			want: want{
				Normalized:  "example.co.il",
				ASCII:       "example.co.il",
				TLD:         "co.il",
				Registrable: "example.co.il",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "it.com",
			in:   "example.it.com",
			want: want{
				Normalized:  "example.it.com",
				ASCII:       "example.it.com",
				TLD:         "it.com",
				Registrable: "example.it.com",
				SLD:         "example",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     false,
			},
		},
		{
			name: "it.com with punycode",
			in:   "xn--b1agh1afp.it.com",
			want: want{
				Normalized:  "xn--b1agh1afp.it.com",
				ASCII:       "xn--b1agh1afp.it.com",
				TLD:         "it.com",
				Registrable: "привет.it.com",
				SLD:         "привет",
				Subdomain:   "",
				IsIDN:       false,
				HasPuny:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.in, err)
			}

			if got.Normalized != tt.want.Normalized {
				t.Errorf("Normalized = %q, want %q", got.Normalized, tt.want.Normalized)
			}
			if got.ASCII != tt.want.ASCII {
				t.Errorf("ASCII = %q, want %q", got.ASCII, tt.want.ASCII)
			}
			if got.Tld != tt.want.TLD {
				t.Errorf("TLD = %q, want %q", got.Tld, tt.want.TLD)
			}
			if got.Registerable != tt.want.Registrable {
				t.Errorf("Registerable = %q, want %q", got.Registerable, tt.want.Registrable)
			}
			if got.Sld != tt.want.SLD {
				t.Errorf("SLD = %q, want %q", got.Sld, tt.want.SLD)
			}
			if got.SubDomain != tt.want.Subdomain {
				t.Errorf("Subdomain = %q, want %q", got.SubDomain, tt.want.Subdomain)
			}
			if got.IsIDN != tt.want.IsIDN {
				t.Errorf("IsIDN = %v, want %v", got.IsIDN, tt.want.IsIDN)
			}
			if got.HasPunycode != tt.want.HasPuny {
				t.Errorf("HasPunycode = %v, want %v", got.HasPunycode, tt.want.HasPuny)
			}
		})
	}
}
