package macros

import (
	"testing"
	"time"
)

func TestFields(t *testing.T) {
	p, _ := NewParser()
	values := []Value{
		String("foo", "bar"),
		Float64("bar", 4.2),
		Int("answer", -42),
		Uint("answer+", 42),
		Bool("ok", true),
		Bool("not", false),
	}
	{
		buf, err := p.AppendReplace(nil, "${foo} ${bar}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "bar 4.2" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${foo}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "bar" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${bar}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "4.2" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${answer}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "-42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${answer+}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "42" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${ok}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "true" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${not}", values...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		} else if string(buf) != "false" {
			t.Errorf("Invalid replacement %q", buf)
		}
	}
	{
		buf, err := p.AppendReplace(nil, "${bax}", values...)
		if err != ErrMacroNotFound {
			t.Errorf("Unexpected error %s", err)
		} else if buf != nil {
			t.Errorf("Invalid buf %v", buf)
		}
	}

}

func BenchmarkInterfaceAlloc(b *testing.B) {
	p, _ := NewParser()
	tpl := "${foo} ${bar}"
	var buf []byte
	var err error
	b.ReportAllocs()
	now := time.Now()
	for i := 0; i < b.N; i++ {
		buf, err = p.AppendReplace(buf[:0], tpl,
			String("foo", "bar"),
			Bool("bar", true),
			// Bind("zap", Float64Value(42.0)),
			// Bind("zap", TimeValue{now, time.RFC3339Nano}),
			Time("ts", now, time.RFC3339Nano),
		)
		if err != nil {
			b.Error(err)
			return
		}
		if string(buf) != "bar true" {
			b.Errorf("Invalid replace %q", buf)
			return

		}
	}

}
