package macros

import (
	"io"
	"strings"
)

// Replace appends the template `tpl` to `buf` replacing tokens with values
func Replace(buf []byte, tpl string, values ...Value) ([]byte, error) {
	var r Replacer
	return r.Replace(buf, tpl, values...)
}

// FPrintf writes the template `tpl` to an `io.Writer` replacing tokens with values using a blank Replacer
func FPrintf(w io.Writer, tpl string, values ...Value) (int, error) {
	var p Printer
	return p.FPrintf(w, tpl, values...)
}

// SPrintf writes the template `tpl` to an `io.Writer` replacing tokens with values using a blank Replacer
func SPrintf(tpl string, values ...Value) (string, error) {
	var p Printer
	return p.SPrintf(tpl, values...)
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

func (token Token) alias(alias Token) Token {
	if pos := strings.IndexByte(string(token), TokenDelimiter); 0 <= pos && pos < len(token) {
		return alias + token[pos:]
	}
	return alias
}

func (token Token) split() (Token, Token) {
	if pos := strings.IndexByte(string(token), TokenDelimiter); 0 <= pos && pos < len(token) {
		return token[:pos], token[pos+1:]
	}
	return token, ""
}
