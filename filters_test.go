package macros

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestQueryEscape(t *testing.T) {
	data, err := QueryEscape(nil, []byte("Hello world!"))
	if err != nil {
		t.Error(err)
	}
	if string(data) != "Hello+world%21" {
		t.Errorf("Invalid filter %s", data)
	}

}

func TestHex(t *testing.T) {
	v := []byte("Hello world!")
	expect := hex.EncodeToString(v)
	data, err := Hex(nil, v)
	if err != nil {
		t.Error(err)
	}
	if string(data) != expect {
		t.Errorf("Invalid filter result %q", data)
	}

}
func TestBase64URL(t *testing.T) {
	v := []byte("Hello world!")
	expect := base64.URLEncoding.EncodeToString(v)
	data, err := Base64URL(nil, v)
	if err != nil {
		t.Error(err)
	}
	if string(data) != expect {
		t.Errorf("Invalid filter result %q", data)
	}

}
func TestBase64(t *testing.T) {
	v := []byte("Hello world!")
	expect := base64.StdEncoding.EncodeToString(v)
	data, err := Base64(nil, v)
	if err != nil {
		t.Error(err)
	}
	if string(data) != expect {
		t.Errorf("Invalid filter result %q", data)
	}

}

func TestFilters(t *testing.T) {
	p, _ := New(Filters{
		"hex": Hex,
	})
	buf, err := p.Replace(nil, "${foo:hex}", String("foo", "\x00\xff"))
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "00ff" {
		t.Errorf("Invalid filter replacement %q", buf)
	}
}
