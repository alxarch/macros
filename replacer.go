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
}

// New creates a new `Replacer` applying options
func New(options ...Option) (*Replacer, error) {
	var r Replacer
	if err := r.applyOptions(options); err != nil {
		return nil, err
	}
	return &r, nil
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
	macro, _ := token.split()
	if alias, ok := r.alias[macro]; ok {
		return token.alias(alias)
	}
	return token
}

func (r *Replacer) macroAlias(macro Token) Token {
	if alias, ok := r.alias[macro]; ok {
		return alias
	}
	return macro
}

const minBufferSize = 64

func (r *Replacer) applyOptions(options []Option) error {
	for _, opt := range options {
		if err := opt.apply(r); err != nil {
			return err
		}
	}
	return nil
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

func (r *Replacer) appendToken(buf []byte, token Token) []byte {
	start, end := r.Delimiters()
	buf = append(buf, start...)
	buf = append(buf, string(token)...)
	return append(buf, end...)
}

func (r *Replacer) replaceToken(buf []byte, token Token, values []Value) ([]byte, error) {
	macro, filters := token.split()
	macro = r.macroAlias(macro)
	if _, skip := r.skip[macro]; skip {
		return r.appendToken(buf, token), nil
	}
	offset := len(buf)
	buf, err := r.replace(buf, macro, values)
	if err != nil {
		return buf[:offset], err
	}
	if len(filters) == 0 {
		return buf, nil
	}
	value := buf[offset:]
	var name Token
	for len(filters) > 0 {
		name, filters = filters.split()
		filter := r.filters[string(name)]
		if filter == nil {
			return buf[:offset], &MissingFilterError{string(name)}
		}
		n := len(buf)
		buf, err = filter(buf, value)
		if err != nil {
			return buf[:offset], err
		}
		value = buf[n:]
	}
	return append(buf[:offset], value...), nil
}

// ErrMacroNotFound is the error to return when a macro is not found
var ErrMacroNotFound = errors.New("Macro not found")

func (r *Replacer) replace(buf []byte, macro Token, values []Value) ([]byte, error) {
	var v *Value
	for i := range values {
		v = &values[i]
		if v.macro == macro {
			if v.typ == typeExpand {
				return r.Replace(buf, v.str, values...)
			}
			return v.AppendValue(buf)
		}
	}
	return r.none.AppendValue(buf)
}

// Filters is maps names to filters
type Filters map[string]Filter

func (filters Filters) apply(r *Replacer) error {
	if len(filters) == 0 {
		return nil
	}

	if r.filters == nil {
		r.filters = Filters{}
	}
	for name, filter := range filters {
		r.filters[name] = filter
	}
	return nil
}

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
