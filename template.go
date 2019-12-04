package macros

import (
	"net/url"
	"sort"
	"strings"
)

// Template is a compiled template
type Template struct {
	parsed
	p *Parser
}

// Parser returns the parser that parsed the template
func (t *Template) Parser() *Parser {
	return t.p
}

// String renders a template
func (t *Template) String() string {
	var w strings.Builder
	t.p.render(&w, t.parsed)
	return w.String()
}

type parsed struct {
	chunks []chunk
	tail   string
}

type chunk struct {
	prefix string
	token  Token
}

// Must creates a new templates or panics if there were any errors
func Must(tpl string, options ...Option) *Template {
	t, err := Parse(tpl, options...)
	if err != nil {
		panic(err)
	}
	return t
}

// Parse creates a new template aplying options
func Parse(tpl string, options ...Option) (*Template, error) {
	p, err := NewParser(options...)
	if err != nil {
		return nil, err
	}
	return p.Parse(tpl)
}

// EstimateSize estimates the rendered buffer size
func (t *parsed) EstimateSize(size int) int {
	size *= len(t.chunks)
	for i := range t.chunks {
		chunk := &t.chunks[i]
		size += len(chunk.prefix)
	}
	size += len(t.tail)
	return size
}

// AppendReplace appends the template to a buffer replacing tokens with values
func (t *Template) AppendReplace(buf []byte, values ...Value) ([]byte, error) {
	return t.p.appendTemplate(buf, t.parsed, values)
}

const defaultStartDelimiter = "${"
const defaultEndDelimiter = "}"

func DefaultDelimiters() (string, string) {
	return defaultStartDelimiter, defaultEndDelimiter
}

// var _ Option = Delimiters{}

// URL converts a URL string to a template string with a query
func (p *Parser) URL(rawurl string, params map[string]Token) (string, error) {
	start, end := p.Delimiters()
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
