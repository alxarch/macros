package macros

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// Filter is a converter for values
type Filter func(dst, value []byte) ([]byte, error)

// QueryEscape is a filter escaping a value for URL query strings
func QueryEscape(dst, value []byte) ([]byte, error) {
	q := url.QueryEscape(string(value))
	return append(dst, q...), nil
}

// Base64 is a filter converting a value to base64 string
func Base64(dst, value []byte) ([]byte, error) {
	size := base64.StdEncoding.EncodedLen(len(value))
	buf := make([]byte, size)
	base64.StdEncoding.Encode(buf, value)
	return append(dst, buf...), nil
}

// Base64URL is a filter converting a value to base64 string for URLs
func Base64URL(dst, value []byte) ([]byte, error) {
	size := base64.URLEncoding.EncodedLen(len(value))
	buf := make([]byte, size)
	base64.URLEncoding.Encode(buf, value)
	return append(dst, buf...), nil
}

// Hex is a filter converting a value to hex string
func Hex(dst, value []byte) ([]byte, error) {
	buf := make([]byte, 2*len(value))
	hex.Encode(buf, value)
	return append(dst, buf...), nil
}

// Filters is a series of filters to apply to replacements
type Filters map[string]Filter

var _ Option = (Filters)(nil)

// option implements `Option` interface
func (m Filters) option(t *Template) {
	if t.filters == nil {
		t.filters = Filters{}
	}
	for name, filter := range m {
		t.filters[name] = filter
	}
}

// Apply applies the filter `name` to `value` appending to `dst`
func (m Filters) Apply(dst []byte, value []byte, name string) ([]byte, error) {
	filter, ok := m[name]
	if !ok {
		return nil, &MissingFilterError{name}
	}
	dst, err := filter(dst, value)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// MissingFilterError is an error for missing macro filter
type MissingFilterError struct {
	filter string
}

func (m *MissingFilterError) Error() string {
	return "Missing macro filter " + m.filter
}
