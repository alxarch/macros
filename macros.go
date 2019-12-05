package macros

import "io"

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
