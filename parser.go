package macros

import (
	"errors"
	"fmt"
	"strings"
)

// Parser is a macro template parser
type Parser struct {
	start   string
	end     string
	filters FilterMap
	none    Value
	skip    map[Token]struct{}
	alias   map[Token]Token
}

// Option is a parser option
type Option func(p *Parser) error

// NewParser creates a new parser applying options
func NewParser(options ...Option) (p Parser, err error) {
	for _, option := range options {
		if err = option(&p); err != nil {
			return
		}
	}
	return
}

// Delimiters returns the parser's delimiters
func (p *Parser) Delimiters() (start string, end string) {
	start, end = p.start, p.end
	if start == "" {
		start = defaultStartDelimiter
	}
	if end == "" {
		end = defaultEndDelimiter
	}
	return
}

var errEOF = errors.New("EOF")

// Parse compiles a new template
func (p *Parser) Parse(s string) (*Template, error) {
	var (
		t     parsed
		err   error
		chunk chunk
	)
	for len(s) > 0 {
		if s, err = p.parseToken(s, &chunk); err != nil {
			if err == errEOF {
				t.tail = s
				return &Template{t, p}, nil
			}
			return nil, err
		}
		macro, filters := chunk.token.split()
		macro = p.macroAlias(macro)
		if filters == "" {
			chunk.token = macro
		} else {
			chunk.token = NewToken(string(macro), filters.Filters()...)
		}
		t.chunks = append(t.chunks, chunk)
	}
	return &Template{t, p}, nil
}

// Alias returns an alias for a token
func (p *Parser) Alias(token Token) Token {
	macro, _ := token.split()
	if alias, ok := p.alias[macro]; ok {
		filters := token.Filters()
		return NewToken(string(alias), filters...)
	}
	return token
}

func (p *Parser) macroAlias(macro Token) Token {
	if alias, ok := p.alias[macro]; ok {
		return alias
	}
	return macro
}

func (p *Parser) render(w *strings.Builder, t parsed) {
	start, end := p.Delimiters()
	for i := range t.chunks {
		chunk := &t.chunks[i]
		w.WriteString(chunk.prefix)
		macro, filters := chunk.token.split()
		macro = p.macroAlias(macro)
		w.WriteString(start)
		w.WriteString(string(macro))
		if filters != "" {
			w.WriteByte(TokenDelimiter)
			w.WriteString(string(filters))
		}
		w.WriteString(end)
	}
	w.WriteString(t.tail)
}

func (p *Parser) appendTemplate(buf []byte, t parsed, values []Value) ([]byte, error) {
	var (
		err      error
		original = buf[:]
	)
	for i := range t.chunks {
		chunk := &t.chunks[i]
		buf = append(buf, chunk.prefix...)
		if buf, err = p.replaceToken(buf, chunk.token, values); err != nil {
			return original, err
		}
	}
	return append(buf, t.tail...), nil
}

// AppendReplace appends the template to a buffer replacing tokens with values
func (p *Parser) AppendReplace(buf []byte, tpl string, values ...Value) ([]byte, error) {
	var (
		err      error
		original = buf[:]
		chunk    chunk
	)
	for len(tpl) > 0 {
		if tpl, err = p.parseToken(tpl, &chunk); err != nil {
			if err == errEOF {
				return append(buf, tpl...), nil
			}
			return original, err
		}
		buf = append(buf, chunk.prefix...)
		if buf, err = p.replaceToken(buf, chunk.token, values); err != nil {
			return original, err
		}
	}
	return buf, nil
}

func unmatchedDelimiterError(d string, pos int) error {
	return fmt.Errorf("Unmatched delimiter %q at position %d", d, pos)
}

func (p *Parser) parseToken(s string, chunk *chunk) (string, error) {
	start, end := p.Delimiters()
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

func (p *Parser) replaceToken(buf []byte, token Token, values []Value) ([]byte, error) {
	macro, filters := token.split()
	if _, skip := p.skip[macro]; skip {
		start, end := p.Delimiters()
		buf = append(buf, start...)
		buf = append(buf, string(token)...)
		buf = append(buf, end...)
		return buf, nil
	}
	offset := len(buf)
	buf, err := p.replace(buf, macro, values)
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
		filter := p.filters[string(name)]
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

func (p *Parser) replace(buf []byte, macro Token, values []Value) ([]byte, error) {
	var v *Value
	for i := range values {
		v = &values[i]
		if v.macro == macro {
			if v.typ == typeExpand {
				return p.AppendReplace(buf, v.str, values...)
			}
			return v.AppendValue(buf)
		}
	}
	return p.none.AppendValue(buf)
}

// FilterMap is maps names to filters
type FilterMap map[string]Filter

// Filters sets parser filters
func Filters(filters FilterMap) Option {
	return func(p *Parser) error {
		if len(filters) == 0 {
			return nil
		}

		if p.filters == nil {
			p.filters = FilterMap{}
		}
		for name, filter := range filters {
			p.filters[name] = filter
		}
		return nil
	}
}

// Delimiters sets parser delimiters
func Delimiters(start, end string) Option {
	return func(p *Parser) error {
		start = strings.TrimSpace(start)
		if start == "" {
			return fmt.Errorf("Invalid start delimiter")
		}
		end = strings.TrimSpace(end)
		if end == "" {
			return fmt.Errorf("Invalid end delimiter")
		}
		p.start, p.end = start, end
		return nil
	}
}

// Alias defines aliases for a macro
func Alias(macro Token, aliases ...Token) Option {
	return func(p *Parser) error {
		if p.alias == nil {
			p.alias = make(map[Token]Token)
		}
		macro, _ = macro.split()
		for _, alias := range aliases {
			alias, _ = alias.split()
			p.alias[macro] = alias
		}
		return nil
	}
}

// DefaultValue sets a value to be used when no macro replacement is found
func DefaultValue(value string) Option {
	return func(p *Parser) error {
		p.none = String("", value)
		return nil
	}
}

// Skip defines tokens that will not be replaced
func Skip(tokens ...Token) Option {
	return func(p *Parser) error {
		if p.skip == nil {
			p.skip = make(map[Token]struct{}, len(tokens))
		}
		for _, token := range tokens {
			p.skip[token] = struct{}{}
		}
		return nil
	}

}
