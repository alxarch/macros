package macros

import (
	"net/url"
	"sort"
	"strings"
)

// Template is a compiled template
type Template struct {
	chunks     []chunk
	tail       string
	delimiters Delimiters
	filters    Filters
}

// Option is a Template option
type Option interface {
	option(t *Template)
}

type chunk struct {
	prefix string
	macro  string
}

// Must creates a new templates or panics if there were any errors
func Must(tpl string, options ...Option) *Template {
	t, err := New(tpl, options...)
	if err != nil {
		panic(err)
	}
	return t
}

// New creates a new template aplying options
func New(tpl string, options ...Option) (*Template, error) {
	t := Template{
		delimiters: DefaultDelimiters(),
	}
	for _, opt := range options {
		opt.option(&t)
	}
	if err := t.parse(tpl); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *Template) checkMacro(macro string) error {
	_, macro = splitMacro(macro)
	for len(macro) > 0 {
		filter, tail := splitMacro(macro)
		if _, ok := t.filters[filter]; !ok {
			return &MissingFilterError{filter}
		}
		macro = tail
	}
	return nil
}

// Compile compiles a new template
func (t *Template) parse(tpl string) error {
	i := 0
	for 0 <= i && i < len(tpl) {
		tail := tpl[i:]
		i = strings.Index(tail, t.delimiters.Start)
		if 0 <= i && i < len(tail) {
			prefix := tail[:i]
			i += len(t.delimiters.Start)
			if 0 <= i && i < len(tail) {
				tail = tail[i:]
				j := strings.Index(tail, t.delimiters.End)
				if 0 <= j && j < len(tail) {
					macro := tail[:j]
					if err := t.checkMacro(macro); err != nil {
						return err
					}
					t.chunks = append(t.chunks, chunk{
						prefix: prefix,
						macro:  macro,
					})
					tpl = tail[j:]
					i = len(t.delimiters.End)
					continue
				} else {
					// End not found
					continue
				}
			}
		} else {
			t.tail = tail
			break
		}
	}
	return nil
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

func (t *Template) replaceMacro(buf []byte, macro string, values MacroValues) ([]byte, error) {
	macro, filters := splitMacro(macro)
	offset := len(buf)
	v := values.ReplaceMacro(macro)
	if v.typ == fieldTypeNone {
		return nil, ErrMacroNotFound
	}
	buf = v.appendTo(buf)
	var err error
	if len(filters) == 0 {
		return buf, nil
	}
	value := buf[offset:]
	var name string
	for len(filters) > 0 {
		name, filters = splitMacro(filters)
		n := len(buf)
		buf, err = t.filters.Apply(buf, value, name)
		if err != nil {
			return buf[:offset], err
		}
		value = buf[n:]
	}
	return append(buf[:offset], value...), nil

}

// AppendTo appends the template to a buffer replacing values from fields
func (t *Template) AppendTo(buf []byte, values ...MacroValue) ([]byte, error) {
	var (
		err      error
		r        = MacroValues(values)
		original = buf[:]
	)
	for i := range t.chunks {
		chunk := &t.chunks[i]
		buf = append(buf, chunk.prefix...)
		if buf, err = t.replaceMacro(buf, chunk.macro, r); err != nil {
			return original, err
		}
	}
	buf = append(buf, t.tail...)
	return buf, nil
}

// Delimiters returns the template's delimiters
func (t *Template) Delimiters() Delimiters {
	return t.delimiters
}

func splitMacro(macro string) (string, string) {
	if pos := strings.IndexByte(macro, ':'); 0 <= pos && pos < len(macro) {
		return macro[:pos], macro[pos+1:]
	}
	return macro, ""
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
func (d *Delimiters) AppendToken(dst []byte, macro string) []byte {
	dst = append(dst, d.Start...)
	dst = append(dst, macro...)
	dst = append(dst, d.End...)
	return dst
}

var _ Option = Delimiters{}

func (d Delimiters) option(t *Template) {
	t.delimiters = d
}

// Token renders a macro token
func (d *Delimiters) Token(macro string) string {
	return d.Start + macro + d.End
}

// URL converts a URL string to a template string with a query
func (d *Delimiters) URL(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, err
	}
	q := u.Query()
	for key, macro := range params {
		q.Set(key, d.Token(macro))
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
