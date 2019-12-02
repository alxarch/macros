package macros

import "testing"

func TestFields(t *testing.T) {
	fields := MacroValues{
		Bind("foo", String("bar")),
		Bind("bar", Float(4.2)),
		Bind("answer", Int(-42)),
		Bind("answer+", Uint(42)),
		Bind("ok", Bool(true)),
		Bind("not", Bool(false)),
	}
	{
		v := fields.ReplaceMacro("foo")
		if buf := v.appendTo(nil); string(buf) != "bar" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("bar")
		if buf := v.appendTo(nil); string(buf) != "4.2" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("answer")
		if buf := v.appendTo(nil); string(buf) != "-42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("answer+")
		if buf := v.appendTo(nil); string(buf) != "42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("ok")
		if buf := v.appendTo(nil); string(buf) != "true" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("not")
		if buf := v.appendTo(nil); string(buf) != "false" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		v := fields.ReplaceMacro("baz")
		if !v.IsNone() {
			t.Errorf("Invalid value %v", v)
		}
	}

}
