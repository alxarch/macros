package macros

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

type Filter func(dst, value []byte) ([]byte, error)

func QueryEscape(dst, value []byte) ([]byte, error) {
	q := url.QueryEscape(b2s(value))
	return append(dst, q...), nil
}

func Base64(dst, value []byte) ([]byte, error) {
	size := base64.StdEncoding.EncodedLen(len(value))
	buf := make([]byte, size)
	base64.StdEncoding.Encode(buf, value)
	return append(dst, buf...), nil
}

func Base64URL(dst, value []byte) ([]byte, error) {
	size := base64.URLEncoding.EncodedLen(len(value))
	buf := make([]byte, size)
	base64.URLEncoding.Encode(buf, value)
	return append(dst, buf...), nil
}

func Hex(dst, value []byte) ([]byte, error) {
	buf := make([]byte, 2*len(value))
	hex.Encode(buf, value)
	return append(dst, buf...), nil
}

type FilterMap map[string]Filter

type Filters struct {
	R       ReplacerFunc
	Filters FilterMap
}

func (r *Filters) Replace(dst []byte, macro string) ([]byte, error) {
	for i := 0; 0 <= i && i < len(macro); i++ {
		if macro[i] == ':' {
			var err error
			offset := len(dst)
			dst, err = r.R(dst, macro[:i])
			if err != nil {
				return dst[:offset], err
			}
			value := dst[offset:]

			for i++; 0 <= i && i < len(macro); i++ {
				name := macro[i:]
				if j := strings.IndexByte(name, ':'); 0 <= j && j < len(name) {
					name = name[:j]
				}
				i += len(name)
				filter := r.Filters[name]
				if filter == nil {
					return dst, fmt.Errorf("Missing macro filter %s", name)
				}
				j := len(dst)
				dst, err = filter(dst, value)
				if err != nil {
					return dst[:offset], err
				}
				value = dst[j:]
			}
			return append(dst[:offset], value...), nil
		}
	}
	return r.R(dst, macro)
}
