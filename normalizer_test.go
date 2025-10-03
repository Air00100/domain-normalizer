package normalizer

import "testing"

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
		//nolint:gosmopolitan
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
