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
func Base64(buf, value []byte) ([]byte, error) {
	size := base64.StdEncoding.EncodedLen(len(value))
	offset := len(buf)
	buf = append(buf, make([]byte, size)...)
	base64.StdEncoding.Encode(buf[offset:], value)
	return buf, nil
}
func growBuffer(buf []byte, size int) []byte {
	return append(buf, make([]byte, size)...)
}

// Base64URL is a filter converting a value to base64 string for URLs
func Base64URL(buf, value []byte) ([]byte, error) {
	size := base64.URLEncoding.EncodedLen(len(value))
	offset := len(buf)
	buf = growBuffer(buf, size)
	base64.URLEncoding.Encode(buf[offset:], value)
	return buf, nil
}

// Hex is a filter converting a value to hex string
func Hex(buf, value []byte) ([]byte, error) {
	offset := len(buf)
	buf = growBuffer(buf, 2*len(value))
	hex.Encode(buf[offset:], value)
	return buf, nil
}

// Filters is a series of filters to apply to replacements
type Filters map[string]Filter

var _ Option = (Filters)(nil)

// option implements `Option` interface
func (m Filters) option(p *Parser) {
	if p.filters == nil {
		p.filters = Filters{}
	}
	for name, filter := range m {
		p.filters[name] = filter
	}
}

// MissingFilterError is an error for missing macro filter
type MissingFilterError struct {
	filter string
}

func (m *MissingFilterError) Error() string {
	return "Missing macro filter " + m.filter
}
