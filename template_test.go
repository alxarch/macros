package macros

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl, err := New("${FOO}")
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

	d := tpl.Delimiters()
	if token := d.AppendToken(nil, "foo"); string(token) != "${foo}" {
		t.Errorf("Invalid token %q", string(token))
	}
	buf, err := tpl.AppendTo(nil)
	if err != ErrMacroNotFound {
		t.Errorf("Invalid error %s", err)
	}
	buf, err = tpl.AppendTo(nil, Bind("FOO", String("bar")))

	if err != nil {
		t.Error(err)
		return
	}
	if string(buf) != "bar" {
		t.Errorf("Invalid buf: %q", buf)
	}
}

func TestURLTemplate(t *testing.T) {
	params := map[string]string{
		"foo": "FOO",
		"bar": "BAR",
	}
	delim := DefaultDelimiters()
	tpl, err := delim.URL("http://example.org/foo/bar?foo=bar&bar=baz&baz=foo", params)
	if err != nil {
		t.Fatal(err)
	}
	expect := "http://example.org/foo/bar?bar=${BAR}&baz=foo&foo=${FOO}"
	if tpl != expect {
		t.Errorf("Invalid url %s != %s", tpl, expect)
	}
}
