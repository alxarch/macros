package macros

import (
	"bytes"
)

type QuickTemplate struct {
	start, end []byte
}

func Quick(start, end string) *QuickTemplate {
	return &QuickTemplate{s2b(start), s2b(end)}
}

// Replace replaces macros in a template and appends the output to dst.
func (t *QuickTemplate) Replace(dst []byte, tpl string, r Replacer) ([]byte, error) {
	var (
		b   = s2b(tpl)
		s   = t.start
		e   = t.end
		i   int
		err error
	)
	if r == nil {
		r = NopReplacer{}
	}
	for {
		i = bytes.Index(b, s)
		if 0 <= i && i < len(b) {
			dst = append(dst, b[:i]...)
			if i += len(s); 0 <= i && i < len(b) {
				b = b[i:]
				i = bytes.Index(b, e)
				if 0 <= i && i < len(b) {
					dst, err = r.Replace(dst, b2s(b[:i]))
					if err != nil {
						return dst, err
					}
					if i += len(e); 0 <= i && i < len(b) {
						b = b[i:]
					} else {
						// eof
						b = nil
					}

				} else {
					dst = append(dst, s...)
					break
				}
			}
		} else {
			break
		}
	}
	dst = append(dst, b...)
	return dst, nil
}

type Replacer interface {
	Replace(dst []byte, macro string) ([]byte, error)
}

type ReplacerFunc func(dst []byte, macro string) ([]byte, error)

func Replacements(pairs ...string) Replacer {
	return replacePairs(pairs)
}

type replacePairs []string

func (r replacePairs) Replace(dst []byte, macro string) ([]byte, error) {
	for len(r) > 1 {
		if r[0] == macro {
			return append(dst, r[1]...), nil
		}
		r = r[2:]
	}
	return dst, nil
}

func (f ReplacerFunc) Replace(dst []byte, macro string) ([]byte, error) {
	return f(dst, macro)
}

type NopReplacer struct{}

func (NopReplacer) Replace(dst []byte, _ string) ([]byte, error) {
	return dst, nil
}

func nopReplace(dst []byte, _ string) ([]byte, error) {
	return dst, nil
}

// // Execute replaces macros in a template and writes the output to w.
// func (t *QuickTemplate) Execute(w io.Writer, tpl string, a Replacer, buf []byte) (n int64, err error) {
// 	var (
// 		b  = s2b(tpl)
// 		i  int
// 		nn int
// 	)
// 	for {
// 		i = bytes.Index(b, t.start)
// 		if 0 <= i && i < len(b) {
// 			nn, err = w.Write(b[:i])
// 			if n += int64(nn); err != nil {
// 				return
// 			}
// 			if i += len(t.start); 0 <= i && i < len(b) {
// 				b = b[i:]
// 				i = bytes.Index(b, t.end)
// 				if 0 <= i && i < len(b) {
// 					if t.Filters == nil {
// 						buf, err = a(buf[:0], b2s(b[:i]))
// 					} else {
// 						buf, err = t.Filters.replace(buf[:0], b2s(b[:i]), a)
// 					}
// 					if err != nil {
// 						return
// 					}
// 					if len(buf) > 0 {
// 						nn, err = w.Write(buf)
// 						if n += int64(nn); err != nil {
// 							return
// 						}
// 					}
// 					if i += len(t.end); 0 <= i && i < len(b) {
// 						b = b[i:]
// 					} else {
// 						// eof
// 						b = nil
// 					}

// 				} else {
// 					nn, err = w.Write(t.end)
// 					if n += int64(nn); err != nil {
// 						return
// 					}
// 					break
// 				}
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	nn, err = w.Write(b)
// 	n += int64(nn)
// 	return
// }
