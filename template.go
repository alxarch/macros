package macros

import (
	"strings"
)

// Template is a compiled template
type Template struct {
	chunks []chunk
	tail   string
}

// String renders a template
func (t *Template) String() string {
	var r Replacer
	var w strings.Builder
	r.render(&w, t)
	return w.String()
}

type chunk struct {
	prefix string
	token  Token
}

// Must creates a new templates or panics if there were any errors
func Must(tpl string, options ...Option) Template {
	t, err := Parse(tpl, options...)
	if err != nil {
		panic(err)
	}
	return t
}

// Parse creates a new template aplying options
func Parse(tpl string, options ...Option) (Template, error) {
	p, err := New(options...)
	if err != nil {
		return Template{}, err
	}
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

// Execute executes a template using a blank replacer appending it to a buffer
func (t *Template) Execute(buf []byte, values ...Value) ([]byte, error) {
	var r Replacer
	return r.execute(buf, t, values)
}

const defaultStartDelimiter = "${"
const defaultEndDelimiter = "}"

// DefaultDelimiters returns the default delimiters
func DefaultDelimiters() (string, string) {
	return defaultStartDelimiter, defaultEndDelimiter
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
