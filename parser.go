package macros

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Parser is a macro template parser
type Parser struct {
	delims  Delimiters
	filters Filters
	none    Value
	skip    map[Token]struct{}
	alias   map[Token]Token
	expand  map[Token]*Template
}

// Option is a parser option
type Option interface {
	option(p *Parser)
}

// NewParser creates a template parser
func NewParser(options ...Option) *Parser {
	p := Parser{
		delims: DefaultDelimiters(),
	}
	for _, opt := range options {
		opt.option(&p)
	}
	return &p
}

// Parse compiles a new template
func (p *Parser) Parse(s string) (*Template, error) {
	var (
		token Token
		t     Template
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
			if exp := p.expand[macro]; exp != nil {
				s = prefix + p.Render(exp) + s
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
			return &t, nil
		}
	}
	return &t, nil
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

func (p *Parser) render(w *strings.Builder, t *Template) {
	for i := range t.chunks {
		chunk := &t.chunks[i]
		w.WriteString(chunk.prefix)
		macro, filters := chunk.token.split()
		macro = p.macroAlias(macro)
		if exp := p.expand[macro]; exp != nil {
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

// Render renders a template
func (p *Parser) Render(t *Template) string {
	var w strings.Builder
	p.render(&w, t)
	return w.String()
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

// AppendTemplate appends the template to a buffer replacing tokens with values
func (p *Parser) AppendTemplate(buf []byte, tpl *Template, values ...Value) ([]byte, error) {
	var (
		err      error
		original = buf[:]
	)
	for i := range tpl.chunks {
		chunk := &tpl.chunks[i]
		buf = append(buf, chunk.prefix...)
		if buf, err = p.replaceToken(buf, chunk.token, values); err != nil {
			return original, err
		}
	}
	return append(buf, tpl.tail...), nil
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
			goto appendValue
		}
	}
	v = &p.none
appendValue:
	switch v.typ {
	case fieldTypeTemplate:
		return p.AppendTemplate(buf, v.any.(*Template), values...)
	case fieldTypeString:
		return append(buf, v.str...), nil
	case fieldTypeFloat:
		f := math.Float64frombits(v.num)
		return strconv.AppendFloat(buf, f, 'f', -1, 64), nil
	case fieldTypeUint:
		return strconv.AppendUint(buf, v.num, 10), nil
	case fieldTypeInt:
		return strconv.AppendInt(buf, int64(v.num), 10), nil
	case fieldTypeAny:
		return appendAny(buf, v.any)
	case fieldTypeNone:
		return buf, ErrMacroNotFound
	default:
		return nil, errors.New("Invalid value type")
	}
}

type optionFunc func(p *Parser)

func (opt optionFunc) option(p *Parser) {
	opt(p)

}

// Expand assigns a template to be expanded by a macro
func Expand(macro Token, tpl *Template) Option {
	return optionFunc(func(p *Parser) {
		if p.expand == nil {
			p.expand = make(map[Token]*Template)
		}
		macro, _ = macro.split()
		p.expand[macro] = tpl
	})
}

// Alias defines aliases for a macro
func Alias(macro Token, aliases ...Token) Option {
	return optionFunc(func(p *Parser) {
		if p.alias == nil {
			p.alias = make(map[Token]Token)
		}
		macro, _ = macro.split()
		for _, alias := range aliases {
			alias, _ = alias.split()
			p.alias[macro] = alias
		}
	})
}

// DefaultValue sets a value to be used when a macro has no value
func DefaultValue(value string) Option {
	return optionFunc(func(p *Parser) {
		p.none = String("", value)
	})

}

// Skip defines tokens that will not be replaced
func Skip(tokens ...Token) Option {
	return optionFunc(func(p *Parser) {
		if p.skip == nil {
			p.skip = make(map[Token]struct{}, len(tokens))
		}
		for _, token := range tokens {
			p.skip[token] = struct{}{}
		}
	})

}

func appendAny(buf []byte, v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return append(buf, v...), nil
	case []byte:
		return append(buf, v...), nil
	case int:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int8:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case uint:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint8:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case float64:
		return strconv.AppendFloat(buf, float64(v), 'f', -1, 64), nil
	case float32:
		return strconv.AppendFloat(buf, float64(v), 'f', -1, 32), nil
	case bool:
		if v {
			return append(buf, "true"...), nil
		}
		return append(buf, "false"...), nil
	case fmt.Stringer:
		return append(buf, v.String()...), nil
	default:
		s := fmt.Sprintf("%s", v)
		return append(buf, s...), nil
	}
}
