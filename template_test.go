package macros

import (
	"testing"
)

func TestReplacer(t *testing.T) {
	p := New()
	for _, src := range []string{
		"",
		"${FOO}",
		"${FOO} ${BAR}",
		"${} ${BAR}",
		"  ${FOO} ${BAR} ",
	} {
		tpl, err := p.Parse(src)
		if err != nil {
			t.Errorf("[%s] Failed to parse: %s", src, err)
			continue
		}
		if s := tpl.String(); s != src {
			t.Errorf("[%s] Invalid template: %q", src, s)

		}

	}
}
func TestTemplate(t *testing.T) {
	s := "${FOO}"
	tpl, err := Parse(s)
	if err != nil {
		t.Errorf("Unexpected error")
	}
	if len(tpl.tail) != 0 {
		t.Errorf("Invalid tail: %s", tpl.tail)
		return

	}
	if size := tpl.EstimateSize(42); size != 42 {
		t.Errorf("Invalid size estimation: %d", size)
	}

	buf, err := tpl.Replace(nil)
	if err != ErrMacroNotFound {
		t.Errorf("Invalid error %s", err)
	}
	buf, err = tpl.Replace(nil, String("FOO", "bar"))

	if err != nil {
		t.Error(err)
		return
	}
	if string(buf) != "bar" {
		t.Errorf("Invalid buf: %q", buf)
	}
}

func TestURLTemplate(t *testing.T) {
	params := map[string]Token{
		"foo": "FOO",
		"bar": "BAR",
	}
	var p Replacer
	tpl, err := p.URL("http://example.org/foo/bar?foo=bar&bar=baz&baz=foo", params)
	if err != nil {
		t.Fatal(err)
	}
	expect := "http://example.org/foo/bar?bar=${BAR}&baz=foo&foo=${FOO}"
	if tpl != expect {
		t.Errorf("Invalid url %s != %s", tpl, expect)
	}
}
