package macros

import (
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
	expect := make([]byte, 2*len(v))
	hex.Encode(expect, v)
	data, err := Hex(nil, v)
	if err != nil {
		t.Error(err)
	}
	if string(data) != string(expect) {
		t.Errorf("Invalid filter %s", data)
	}

}
