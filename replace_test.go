package macros

import "testing"

func TestFields(t *testing.T) {
	values := MacroValues{
		Bind("foo", String("bar")),
		Bind("bar", Float(4.2)),
		Bind("answer", Int(-42)),
		Bind("answer+", Uint(42)),
		Bind("ok", Bool(true)),
		Bind("not", Bool(false)),
	}
	{
		tpl := Must("${foo}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "bar" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${bar}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "4.2" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${answer}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "-42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${answer+}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${ok}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "true" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${not}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "false" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		tpl := Must("${bax}")
		buf, err := tpl.AppendTo(nil, values...)
		if err != ErrMacroNotFound {
			t.Errorf("Unexpected error %s", err)
		} else if buf != nil {
			t.Errorf("Invalid buf %v", buf)
		}
	}

}
