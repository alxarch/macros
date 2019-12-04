package macros

import "testing"

func TestNewParser(t *testing.T) {
	p, _ := NewParser(Delimiters("{{", "}}"))
	buf, _ := p.AppendReplace(nil, "{{FOO}}", String("FOO", "foo"))
	if string(buf) != "foo" {
		t.Errorf("Invalid replacement %q", buf)
	}
	tpl, _ := p.Parse("{{ FOO }} {{BAR }} ")
	if s := tpl.String(); s != "{{FOO}} {{BAR}} " {
		t.Errorf("Invalid parse %q", s)

	}

}
