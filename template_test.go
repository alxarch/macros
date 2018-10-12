package macros

import (
	"bytes"
	"errors"
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl := NewTemplate("${FOO}", "${", "}")
	if tpl == nil {
		t.Errorf("Nil template")
	}
	if len(tpl.chunks) != len(tpl.macros) {
		t.Errorf("Invalid chunks: %d", len(tpl.chunks))
		return
	}
	if len(tpl.tail) != 0 {
		t.Errorf("Invalid tail: %s", tpl.tail)
		return

	}
	buf := bytes.NewBuffer(nil)
	n, err := tpl.ExecuteBuffer(buf, ReplacerFunc(func(dst []byte, macro string) ([]byte, error) {
		if macro == "FOO" {
			return append(dst, "foo"...), nil
		}
		return dst, errors.New("Invalid macro")
	}), nil)
	if err != nil {
		t.Error(err)
		return
	}
	if int64(buf.Len()) != n {
		t.Errorf("Invalid size: %d", n)
	}
	if buf.String() != "foo" {
		t.Errorf("Invalid result: %s", buf.String())
	}
}
