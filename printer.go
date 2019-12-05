package macros

import (
	"io"
	"strings"
)

// Printer prints templates to writers
// It has a buffer to replace tokens and it's *not* thread-safe
type Printer struct {
	config Replacer
	buffer []byte
}

// NewPrinter creates a new `Printer` using options
func NewPrinter(size int, options ...Option) (*Printer, error) {
	if size < minBufferSize {
		size = minBufferSize
	}
	p := Printer{
		buffer: make([]byte, size),
	}
	if err := p.config.applyOptions(options); err != nil {
		return nil, err
	}
	return &p, nil
}

// Execute executes a parsed template appending to a buffer
func (p *Printer) Execute(w io.Writer, t *Template, values ...Value) (n int, err error) {
	var nn int
	for i := range t.chunks {
		chunk := &t.chunks[i]
		nn, err = w.Write(s2b(chunk.prefix))
		n += nn
		if err != nil {
			return
		}

		if p.buffer, err = p.config.replaceToken(p.buffer[:0], chunk.token, values); err != nil {
			return
		}
		nn, err = w.Write(p.buffer)
		n += nn
		if err != nil {
			return
		}
	}
	nn, err = w.Write(s2b(t.tail))
	n += nn
	return
}

// SPrintf renders the template `tpl` replacing tokens with values
func (p *Printer) SPrintf(tpl string, values ...Value) (string, error) {
	var (
		err   error
		w     strings.Builder
		chunk chunk
	)
	for len(tpl) > 0 {
		if tpl, err = p.config.parseToken(tpl, &chunk); err != nil {
			if err == errEOF {
				w.WriteString(tpl)
				return w.String(), nil
			}
			return "", err
		}
		w.WriteString(chunk.prefix)
		if p.buffer, err = p.config.replaceToken(p.buffer[:0], chunk.token, values); err != nil {
			return "", err
		}
		w.Write(p.buffer)
	}
	return w.String(), nil
}

// FPrintf writes the template `tpl` to an `io.Writer` replacing tokens with values
func (p *Printer) FPrintf(w io.Writer, tpl string, values ...Value) (n int, err error) {
	var (
		nn    int
		chunk chunk
	)

	for len(tpl) > 0 {
		if tpl, err = p.config.parseToken(tpl, &chunk); err != nil {
			if err == errEOF {
				nn, err = writeString(w, tpl)
				n += nn
			}
			return
		}
		nn, err = writeString(w, tpl)
		n += nn
		if err != nil {
			return
		}
		if p.buffer, err = p.config.replaceToken(p.buffer[:0], chunk.token, values); err != nil {
			return n, err
		}
		nn, err = w.Write(p.buffer)
		n += nn
		if err != nil {
			return
		}
	}
	return n, nil
}
