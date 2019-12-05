package macros

import "testing"

func TestDelimiters(t *testing.T) {
	{
		p := New(Delimiters("%", "%"))
		b, err := p.Replace(nil, "foo %bar%", String("bar", "baz"))
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		if string(b) != "foo baz" {
			t.Errorf("Invalid replacement %q", b)
		}

	}

}

func TestAlias(t *testing.T) {
	r := New(Alias("foo", "FOO", "Foo"))
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
	r := New(Skip("foo", "bar"))
	b, err := r.Replace(nil, "${foo} ${baz}", String("baz", "bar"))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if string(b) != "${foo} bar" {
		t.Errorf("Invalid replacement %q", b)
	}
}
