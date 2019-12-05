package macros

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Replacer is a macro template Replacer
type Replacer struct {
	start   string
	end     string
	filters Filters
	none    Value
	skip    map[Token]struct{}
	alias   map[Token]Token
	expand  map[Token]string
}

// New creates a new `Replacer` applying options
func New(options ...Option) *Replacer {
	var r Replacer
	r.applyOptions(options)
	return &r
}

// Delimiters returns the Replacer's delimiters
func (r *Replacer) Delimiters() (start string, end string) {
	start, end = r.start, r.end
	if start == "" {
		start = defaultStartDelimiter
	}
	if end == "" {
		end = defaultEndDelimiter
	}
	return
}

var errEOF = errors.New("EOF")

// Parse compiles a new template using `r` options
func (r *Replacer) Parse(s string) (*Template, error) {
	t := Template{
		config: *r,
	}
	if err := t.parse(s); err != nil {
		return nil, err
	}
	return &t, nil
}

// Alias returns an alias for a token
func (r *Replacer) Alias(token Token) Token {
	if pos := strings.IndexByte(string(token), TokenDelimiter); 0 <= pos && pos < len(token) {
		macro, filters := token[:pos], token[pos:]
		if alias, ok := r.alias[macro]; ok {
			return alias + filters
		}
		return token
	}
	if alias, ok := r.alias[token]; ok {
		return alias
	}
	return token
}

const minBufferSize = 64

func (r *Replacer) applyOptions(options []Option) {
	for _, opt := range options {
		opt.apply(r)
	}
}

// Replace appends the template to a buffer replacing tokens with values
func (r *Replacer) Replace(buf []byte, tpl string, values ...Value) ([]byte, error) {
	var (
		err      error
		original = buf[:]
		chunk    chunk
	)
	for len(tpl) > 0 {
		if tpl, err = r.parseToken(tpl, &chunk); err != nil {
			if err == errEOF {
				return append(buf, tpl...), nil
			}
			return original, err
		}
		buf = append(buf, chunk.prefix...)
		if buf, err = r.replaceToken(buf, chunk.token, values); err != nil {
			return original, err
		}
	}
	return buf, nil
}

func unmatchedDelimiterError(d string, pos int) error {
	return fmt.Errorf("Unmatched delimiter %q at position %d", d, pos)
}

func (r *Replacer) parseToken(s string, chunk *chunk) (string, error) {
	start, end := r.Delimiters()
	if i := strings.Index(s, start); 0 <= i && i < len(s) {
		var src string
		chunk.prefix, src = s[:i], s[i:]
		if n := len(start); 0 <= n && n < len(src) {
			src := src[n:]
			if i := strings.Index(src, end); 0 <= i && i <= len(src) {
				token := src[:i]
				if strings.Index(token, end) != -1 {
					return s, unmatchedDelimiterError(start, i)
				}
				i += len(end)
				chunk.token = Token(strings.TrimSpace(token))
				return src[i:], nil
			}
		}
		return s, unmatchedDelimiterError(start, i)
	}
	return s, errEOF
}

func (r *Replacer) appendToken(buf []byte, macro, filters Token) []byte {
	start, end := r.Delimiters()
	buf = append(buf, start...)
	buf = append(buf, string(macro)...)
	buf = append(buf, string(filters)...)
	return append(buf, end...)
}

func (r *Replacer) replaceToken(buf []byte, token Token, values []Value) ([]byte, error) {
	var (
		err            error
		value          []byte
		offset         = len(buf)
		macro, filters = token.split()
	)
	if alias, ok := r.alias[macro]; ok {
		macro = alias
	}
	if _, skip := r.skip[macro]; skip {
		return r.appendToken(buf, macro, filters), nil
	}
	if exp, ok := r.expand[macro]; ok {
		buf, err = r.Replace(buf, exp, values...)
		if err != nil {
			return buf[:offset], fmt.Errorf("Expand %q failed: %s", macro, err)
		}
	} else {
		var v *Value
		for i := range values {
			v = &values[i]
			if v.macro == macro {
				goto done
			}
		}
		v = &r.none
	done:
		buf, err = v.AppendValue(buf)
		if err != nil {
			return buf[:offset], err
		}
	}
	if filters == "" {
		return buf, nil
	}
	value = buf[offset:]
	for len(filters) > 1 {
		filters = filters[1:]
		macro, filters = filters.split()
		filter := r.filters[string(macro)]
		if filter == nil {
			return nil, &MissingFilterError{string(macro)}
		}
		n := len(buf)
		if buf, err = filter(buf, value); err != nil {
			return buf[:offset], err
		}
		value = buf[n:]
	}

	return append(buf[:offset], value...), nil
}

// ErrMacroNotFound is the error to return when a macro is not found
var ErrMacroNotFound = errors.New("Macro not found")

// URL converts a URL string to a template string with a query
func (r *Replacer) URL(rawurl string, params map[string]Token) (string, error) {
	start, end := r.Delimiters()
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, err
	}
	q := u.Query()
	for key, macro := range params {
		q.Set(key, start+macro.String()+end)
	}
	bs := strings.Builder{}
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i > 0 {
			bs.WriteByte('&')
		}
		for j, v := range q[k] {
			if j > 0 {
				bs.WriteByte('&')
			}
			bs.WriteString(url.QueryEscape(k))
			bs.WriteByte('=')
			if _, isMacro := params[k]; isMacro {
				bs.WriteString(v)
			} else {
				bs.WriteString(url.QueryEscape(v))
			}
		}
	}
	u.RawQuery = bs.String()
	return u.String(), nil
}
