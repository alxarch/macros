// +build appengine

package macros

func b2s(b []byte) string {
	return string(b)
}

func s2b(s string) []byte {
	return []byte(s)
}
