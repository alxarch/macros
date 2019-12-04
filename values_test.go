package macros

import (
	"strconv"
	"testing"
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

type Float64Value float64

func (v Float64Value) AppendValue(buf []byte) ([]byte, error) {
	return strconv.AppendFloat(buf, float64(v), 'f', -1, 64), nil
}

func BenchmarkInterfaceAlloc(b *testing.B) {
	b.ReportAllocs()
	var buf []byte
	var err error
	for i := 0; i < b.N; i++ {
		v := Any("foo", Float64Value(42.0))
		buf, err = v.AppendValue(buf[:0])
		if err != nil {
			b.Error(err)
			return
		}
	}

}
