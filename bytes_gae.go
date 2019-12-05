// +build appengine

package macros

import "io"

func b2s(b []byte) string {
	return string(b)
}

func s2b(s string) []byte {
	return []byte(s)
}

func writeString(w io.Writer, s string) (int, error) {
	if w, ok := w.(io.StringWriter); ok {
		return w.WriteString(s)
	}
	return w.Write([]byte(s))
}
