package macros

import (
	"errors"
	"fmt"
	"strings"
)

// Parser is a macro template parser
type Parser struct {
	delims  Delimiters
	filters Filters
	none    Value
	skip    map[Token]struct{}
	alias   map[Token]Token
	expand  map[Token]parsed
}

// Option is a parser option
type Option interface {
	option(p *Parser) error
}

// NewParser creates a template parser
func NewParser(options ...Option) (*Parser, error) {
	p := Parser{
		delims: DefaultDelimiters(),
	}
	for _, opt := range options {
		if err := opt.option(&p); err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// Parse compiles a new template
func (p *Parser) Parse(s string) (*Template, error) {
	var (
		token Token
		t     parsed
		err   error
	)
	for len(s) > 0 {
		start := strings.Index(s, p.delims.Start)
		if 0 <= start && start < len(s) {
			prefix, tail := s[:start], s[start:]
			token, s, err = p.parseToken(tail, start)
			if err != nil {
				return nil, err
			}
			macro, filters := token.split()
			macro = p.macroAlias(macro)
			if exp, ok := p.expand[macro]; ok {
				var w strings.Builder
				w.WriteString(prefix)
				p.render(&w, exp)
				w.WriteString(s)
				s = w.String()
				continue
			}
			if filters == "" {
				token = macro
			} else {
				token = NewToken(string(macro), filters.Filters()...)
			}
			t.chunks = append(t.chunks, chunk{
				prefix: prefix,
				token:  token,
			})
		} else {
			t.tail = s
			return &Template{t, p}, nil
		}
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
	for i := range t.chunks {
		chunk := &t.chunks[i]
		w.WriteString(chunk.prefix)
		macro, filters := chunk.token.split()
		macro = p.macroAlias(macro)
		if exp, ok := p.expand[macro]; ok {
			p.render(w, exp)
		} else {
			w.WriteString(p.delims.Start)
			w.WriteString(string(macro))
			if filters != "" {
				w.WriteByte(TokenDelimiter)
				w.WriteString(string(filters))
			}
			w.WriteString(p.delims.End)
		}
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
		token    Token
	)
	for len(tpl) > 0 {
		if i := strings.Index(tpl, p.delims.Start); 0 <= i && i < len(tpl) {
			prefix, tail := tpl[:i], tpl[i:]
			buf = append(buf, prefix...)
			if token, tpl, err = p.parseToken(tail, i); err != nil {
				return original, err
			}
			if buf, err = p.replaceToken(buf, token, values); err != nil {
				return original, err
			}
		} else {
			return append(buf, tpl...), nil
		}
	}
	return buf, nil
}

func unmatchedDelimiterError(d string, pos int) error {
	return fmt.Errorf("Unmatched delimiter %q at position %d", d, pos)
}

func (p *Parser) parseToken(src string, pos int) (Token, string, error) {
	if n := len(p.delims.Start); 0 <= n && n < len(src) {
		src := src[n:]
		if i := strings.Index(src, p.delims.End); 0 <= i && i <= len(src) {
			token := src[:i]
			if strings.Index(token, p.delims.Start) != -1 {
				return "", src, unmatchedDelimiterError(p.delims.Start, pos)
			}
			i += len(p.delims.End)
			return Token(strings.TrimSpace(token)), src[i:], nil
		}
	}
	return "", src, unmatchedDelimiterError(p.delims.Start, pos)
}

func (p *Parser) replaceToken(buf []byte, token Token, values []Value) ([]byte, error) {
	macro, filters := token.split()
	if _, skip := p.skip[macro]; skip {
		return p.delims.AppendToken(buf, token), nil
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
			return v.AppendValue(buf)
		}
	}
	return p.none.AppendValue(buf)
}

type optionFunc func(p *Parser) error

func (opt optionFunc) option(p *Parser) error {
	return opt(p)
}

// Expand assigns a template to be expanded by a macro
func Expand(macro Token, s string) Option {
	return optionFunc(func(p *Parser) error {
		if p.expand == nil {
			p.expand = make(map[Token]parsed)
		}
		macro, _ = macro.split()
		tpl, err := p.Parse(s)
		if err != nil {
			return err
		}
		p.expand[macro] = tpl.parsed
		return nil
	})
}

// Alias defines aliases for a macro
func Alias(macro Token, aliases ...Token) Option {
	return optionFunc(func(p *Parser) error {
		if p.alias == nil {
			p.alias = make(map[Token]Token)
		}
		macro, _ = macro.split()
		for _, alias := range aliases {
			alias, _ = alias.split()
			p.alias[macro] = alias
		}
		return nil
	})
}

// DefaultValue sets a value to be used when a macro has no value
func DefaultValue(value string) Option {
	return optionFunc(func(p *Parser) error {
		p.none = String("", value)
		return nil
	})

}

// Skip defines tokens that will not be replaced
func Skip(tokens ...Token) Option {
	return optionFunc(func(p *Parser) error {
		if p.skip == nil {
			p.skip = make(map[Token]struct{}, len(tokens))
		}
		for _, token := range tokens {
			p.skip[token] = struct{}{}
		}
		return nil
	})

}
