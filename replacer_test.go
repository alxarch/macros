package macros

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	p := New(Delimiters("{{", "}}"))
	buf, _ := p.Replace(nil, "{{FOO}}", String("FOO", "foo"))
	if string(buf) != "foo" {
		t.Errorf("Invalid replacement %q", buf)
	}
	tpl, _ := p.Parse("{{ FOO }} {{BAR }} ")
	if s := tpl.String(); s != "{{FOO}} {{BAR}} " {
		t.Errorf("Invalid parse %q", s)
	}
}

func TestParse(t *testing.T) {
	p := New(Filters{"hex": Hex})
	tpl, err := p.Parse("${foo:hex} bar")
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	expect := Template{
		config: *p,
		chunks: []chunk{{
			token: Token("foo:hex"),
		}},
		tail: " bar",
	}
	if !reflect.DeepEqual(*tpl, expect) {
		t.Errorf("Invalid parse %v", t)

	}

}
