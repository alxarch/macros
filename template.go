package macros

import (
	"io"
	"strings"
)

// Template is a compiled template
type Template struct {
	chunks []chunk
	tail   string
	config Replacer
}

// String renders a template
func (t *Template) String() string {
	var w strings.Builder
	start, end := t.config.Delimiters()
	for i := range t.chunks {
		chunk := &t.chunks[i]
		w.WriteString(chunk.prefix)
		macro, filters := chunk.token.split()
		macro = t.config.macroAlias(macro)
		w.WriteString(start)
		w.WriteString(string(macro))
		if filters != "" {
			w.WriteByte(TokenDelimiter)
			w.WriteString(string(filters))
		}
		w.WriteString(end)
	}
	w.WriteString(t.tail)
	return w.String()
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
	t := Template{}
	if err := t.config.applyOptions(options); err != nil {
		return nil, err
	}
	if err := t.parse(tpl); err != nil {
		return nil, err
	}
	return &t, nil
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

// Replace executes a template appending it to a buffer
func (t *Template) Replace(b []byte, values ...Value) (buf []byte, err error) {
	buf = b
	for i := range t.chunks {
		chunk := &t.chunks[i]
		buf = append(buf, chunk.prefix...)
		if buf, err = t.config.replaceToken(buf, chunk.token, values); err != nil {
			return b, err
		}
	}
	return append(buf, t.tail...), nil
}

// Execute writes a template to `w` replacing macros with `values` using `buffer` as scratch space.
func (t *Template) Execute(w io.Writer, buffer []byte, values ...Value) (n int64, err error) {
	if buffer == nil {
		buffer = make([]byte, minBufferSize)
	}
	var nn int
	for i := range t.chunks {
		chunk := &t.chunks[i]
		nn, err = w.Write(s2b(chunk.prefix))
		n += int64(nn)
		if err != nil {
			return
		}

		if buffer, err = t.config.replaceToken(buffer[:0], chunk.token, values); err != nil {
			return
		}
		nn, err = w.Write(buffer)
		n += int64(nn)
		if err != nil {
			return
		}
	}
	nn, err = w.Write(s2b(t.tail))
	n += int64(nn)
	return
}

func (t *Template) parse(s string) (err error) {
	var chunk chunk
	for len(s) > 0 {
		if s, err = t.config.parseToken(s, &chunk); err != nil {
			if err == errEOF {
				t.tail = s
				err = nil
			}
			return
		}
		macro, filters := chunk.token.split()
		macro = t.config.macroAlias(macro)
		if filters == "" {
			chunk.token = macro
		} else {
			chunk.token = macro + Token(string(TokenDelimiter)) + filters
		}
		t.chunks = append(t.chunks, chunk)
	}
	return
}
