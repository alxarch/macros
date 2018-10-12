package macros

import (
	"bytes"
	"io"
	"strings"
)

type Template struct {
	chunks [][]byte
	macros []string
	tail   []byte
}

func NewTemplate(tpl, start, end string) *Template {
	blob := []byte(tpl)
	s := []byte(start)
	e := []byte(end)
	t := Template{}

	i := 0
	for 0 <= i && i < len(blob) {
		tail := blob[i:]
		i = bytes.Index(tail, s)
		if 0 <= i && i < len(tail) {
			chunk := tail[:i]
			i += len(start)
			if 0 <= i && i < len(tail) {
				tail = tail[i:]
				j := bytes.Index(tail, e)
				if 0 <= j && j < len(tail) {
					t.macros = append(t.macros, string(tail[:j]))
					t.chunks = append(t.chunks, []byte(chunk))
					blob = tail[j:]
					i = len(end)
					continue
				} else {
					// End not found
					continue
				}
			}
		} else {
			t.tail = tail
			break
		}
	}
	return &t
}

func (t *Template) ExecuteString(w *strings.Builder, a ReplacerFunc, buf []byte) (string, error) {
	n := w.Len()
	var err error
	for i, chunk := range t.chunks {
		if len(chunk) > 0 {
			w.Write(chunk)
		}
		macro := t.macros[i]
		if buf, err = a(buf[:0], macro); err != nil {
			goto done
		}
		w.Write(buf)
	}
	w.Write(t.tail)
done:
	s := w.String()
	if 0 <= n && n < len(s) {
		s = s[n:]
	}
	return s, err
}

func (t *Template) EstimateSize(size int) int {
	size *= len(t.macros)
	for _, chunk := range t.chunks {
		size += len(chunk)
	}
	size += len(t.tail)
	return size
}

func (t *Template) ExecuteBuffer(w io.Writer, a ReplacerFunc, buf []byte) (n int64, err error) {
	var nn int
	for i, chunk := range t.chunks {
		if len(chunk) > 0 {
			nn, err = w.Write(chunk)
			if n += int64(nn); err != nil {
				return
			}
		}
		macro := t.macros[i]
		buf, err = a(buf[:0], macro)
		nn, err = w.Write(buf)
		if n += int64(nn); err != nil {
			return
		}
	}
	nn, err = w.Write(t.tail)
	n += int64(nn)
	return
}

func (t *Template) Append(dst []byte, a Replacer) ([]byte, error) {
	var err error
	for i, chunk := range t.chunks {
		dst = append(dst, chunk...)
		dst, err = a.Replace(dst, t.macros[i])
		if err != nil {
			return dst, err
		}
	}
	dst = append(dst, t.tail...)
	return dst, nil
}

func (t *Template) AppendFunc(dst []byte, a ReplacerFunc) ([]byte, error) {
	var err error
	for i, chunk := range t.chunks {
		dst = append(dst, chunk...)
		dst, err = a(dst, t.macros[i])
		if err != nil {
			return dst, err
		}
	}
	dst = append(dst, t.tail...)
	return dst, nil
}
