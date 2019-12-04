package macros

import "testing"

func TestNewParser(t *testing.T) {
	{
		p := NewParser()
		if p.Delimiters() != DefaultDelimiters() {
			t.Errorf("Invalid delimiters for new parser")
		}
	}
	{
		delims := Delimiters{"{{", "}}"}
		p := NewParser(delims)
		if p.Delimiters() != delims {
			t.Errorf("Invalid custom delimiters for new parser")
		}
		buf, _ := p.AppendReplace(nil, "{{FOO}}", String("FOO", "foo"))
		if string(buf) != "foo" {
			t.Errorf("Invalid replacement %q", buf)
		}
		tpl, _ := p.Parse("{{ FOO }} {{BAR }} ")
		if s := p.Render(tpl); s != "{{FOO}} {{BAR}} " {
			t.Errorf("Invalid parse %q", s)

		}

	}

}
