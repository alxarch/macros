package macros

import "testing"

func TestNew(t *testing.T) {
	p, _ := New(Delimiters("{{", "}}"))
	buf, _ := p.Replace(nil, "{{FOO}}", String("FOO", "foo"))
	if string(buf) != "foo" {
		t.Errorf("Invalid replacement %q", buf)
	}
	tpl, _ := p.Parse("{{ FOO }} {{BAR }} ")
	if s := tpl.String(); s != "${FOO} ${BAR} " {
		t.Errorf("Invalid parse %q", s)
	}
}
