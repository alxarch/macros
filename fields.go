package macros

import (
	"math"
	"strconv"
)

type fieldType uint

const (
	_ fieldType = iota
	fieldTypeString
	fieldTypeFloat
	fieldTypeInt
	fieldTypeUint
)

type Field struct {
	macro string
	str   string
	num   uint64
	typ   fieldType
}

type Fields []Field

func (f *Field) Replace(dst []byte) []byte {
	switch f.typ {
	case fieldTypeString:
		return append(dst, f.str...)
	case fieldTypeFloat:
		f := math.Float64frombits(f.num)
		return strconv.AppendFloat(dst, f, 'f', -1, 64)
	case fieldTypeUint:
		return strconv.AppendUint(dst, f.num, 10)
	case fieldTypeInt:
		return strconv.AppendInt(dst, int64(f.num), 10)
	default:
		return dst
	}
}
func replaceFields(fields []Field, dst []byte, macro string) ([]byte, error) {
	for i := range fields {
		f := &fields[i]
		if f.macro == macro {
			return f.Replace(dst), nil
		}
	}
	return dst, nil
}
func Bytes(macro string, v []byte) Field {
	return Field{macro, b2s(v), 0, fieldTypeString}
}
func String(macro, v string) Field {
	return Field{macro, v, 0, fieldTypeString}
}

func Bool(macro string, v bool) Field {
	if v {
		return Field{macro, "true", 0, fieldTypeString}
	}
	return Field{macro, "false", 0, fieldTypeString}
}

func Float(macro string, f float64) Field {
	return Field{macro, "", math.Float64bits(f), fieldTypeFloat}
}
func Uint(macro string, u uint64) Field {
	return Field{macro, "", u, fieldTypeUint}
}
func Int(macro string, i int64) Field {
	return Field{macro, "", uint64(i), fieldTypeInt}
}
