package macros

import "testing"

func TestDelimiters(t *testing.T) {
	{
		p, err := New(Delimiters("%", "%"))
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		b, err := p.Replace(nil, "foo %bar%", String("bar", "baz"))
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		if string(b) != "foo baz" {
			t.Errorf("Invalid replacement %q", b)
		}

	}
	{
		p, err := New(Delimiters("", "%"))
		if err == nil {
			t.Errorf("Unexpected error %s", err)
		}
		if p != nil {
			t.Errorf("Non nil Replacer")
		}
	}
	{
		p, err := New(Delimiters("%", ""))
		if err == nil {
			t.Errorf("Unexpected error %s", err)
		}
		if p != nil {
			t.Errorf("Non nil Replacer")
		}
	}
	{
		p, err := New(Delimiters("", ""))
		if err == nil {
			t.Errorf("Unexpected error %s", err)
		}
		if p != nil {
			t.Errorf("Non nil Replacer")
		}
	}

}

func TestAlias(t *testing.T) {
	r, err := New(Alias("foo", "FOO", "Foo"))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if alias := r.Alias("Foo"); alias != "foo" {
		t.Errorf("Invalid alias %q", alias)
	}
	if alias := r.Alias("FOO"); alias != "foo" {
		t.Errorf("Invalid alias %q", alias)
	}
	if alias := r.Alias("foo"); alias != "foo" {
		t.Errorf("Invalid alias %q", alias)
	}
	b, err := r.Replace(nil, "${Foo}", String("foo", "bar"))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if string(b) != "bar" {
		t.Errorf("Invalid replacement %q", b)
	}
}

func TestSkip(t *testing.T) {
	r, err := New(Skip("foo", "bar"))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	b, err := r.Replace(nil, "${foo} ${baz}", String("baz", "bar"))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if string(b) != "${foo} bar" {
		t.Errorf("Invalid replacement %q", b)
	}
}
