package macros

import (
	"net/url"
	"sort"
	"strings"
)

// Template is a compiled template
type Template struct {
	chunks []chunk
	tail   string
}

type chunk struct {
	prefix string
	token  Token
}

// Must creates a new templates or panics if there were any errors
func Must(tpl string, options ...Option) *Template {
	p := NewParser(options...)
	t, err := p.Parse(tpl)
	if err != nil {
		panic(err)
	}
	return t
}

// Parse creates a new template aplying options
func Parse(tpl string, options ...Option) (*Template, error) {
	p := NewParser(options...)
	return p.Parse(tpl)
}

// EstimateSize estimates the rendered buffer size
func (t *Template) EstimateSize(size int) int {
	size *= len(t.chunks)
	for i := range t.chunks {
		chunk := &t.chunks[i]
		size += len(chunk.prefix)
	}
	size += len(t.tail)
	return size
}

// Delimiters returns the template's delimiters
func (p *Parser) Delimiters() Delimiters {
	return p.delims
}

// Delimiters are the macro delimiters for templates
type Delimiters struct {
	Start string
	End   string
}

// DefaultDelimiters returns the default delimiters for new templates
func DefaultDelimiters() Delimiters {
	return Delimiters{
		Start: "${",
		End:   "}",
	}
}

// AppendToken appends a macro token
func (d *Delimiters) AppendToken(dst []byte, token Token) []byte {
	dst = append(dst, d.Start...)
	dst = append(dst, string(token)...)
	dst = append(dst, d.End...)
	return dst
}

var _ Option = Delimiters{}

func (d Delimiters) option(p *Parser) {
	p.delims = d
}

// Token renders a macro token
func (d *Delimiters) Token(macro Token) string {
	return d.Start + string(macro) + d.End
}

// URL converts a URL string to a template string with a query
func (d *Delimiters) URL(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, err
	}
	q := u.Query()
	for key, macro := range params {
		q.Set(key, d.Token(Token(macro)))
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

// TokenDelimiter is the token delimiter for macro and filters
const TokenDelimiter = ':'

// Token is a macro token with optional filters
type Token string

// NewToken creates a new macro token
func NewToken(macro string, filters ...string) Token {
	return Token(strings.Join(append([]string{macro}, filters...), ":"))
}
func (token Token) String() string {
	return string(token)
}

// Macro returns the macro part of the token
func (token Token) Macro() string {
	m, _ := token.split()
	return string(m)
}

// Filters returns the filters of a token
func (token Token) Filters() []string {
	if _, filters := token.split(); len(filters) > 0 {
		return strings.Split(string(filters), string(TokenDelimiter))
	}
	return nil
}

func (token Token) split() (Token, Token) {
	if pos := strings.IndexByte(string(token), TokenDelimiter); 0 <= pos && pos < len(token) {
		return token[:pos], token[pos+1:]
	}
	return token, ""
}
