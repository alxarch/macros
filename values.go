package macros

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

// Value is a value replacement for a macro
type Value struct {
	macro Token
	str   string
	num   uint64
	typ   valueType
	any   interface{}
}

type valueType uint

const (
	typeNone valueType = iota
	typeString
	typeFloat
	typeInt
	typeUint
	typeAny
	typeTime
	typeConcat
)

// String creates a new value replacing `macro` with a string
func String(macro Token, s string) Value {
	return Value{macro, s, 0, typeString, nil}
}

// Concat creates a new value replacing `macro` with a string
func Concat(macro Token, sep string, values []string) Value {
	return Value{macro, sep, 0, typeConcat, values}
}

// Bool creates a new value replacing `macro` with "true" or "false"
func Bool(macro Token, v bool) Value {
	if v {
		return Value{macro, "true", 0, typeString, nil}
	}
	return Value{macro, "false", 0, typeString, nil}
}

// Float64 creates a new value replacing `macro` with a float64
func Float64(macro Token, f float64) Value {
	return Value{macro, "", math.Float64bits(f), typeFloat, nil}
}

// Float32 creates a new value replacing `macro` with a float32 value
func Float32(macro Token, f float32) Value {
	return Value{macro, "", math.Float64bits(float64(f)), typeFloat, nil}
}

// Uint creates a new value that replaces `macro` with a `uint`
func Uint(macro Token, u uint) Value {
	return Value{macro, "", uint64(u), typeUint, nil}
}

// Uint64 creates a new value that replaces `macro` with a `uint64`
func Uint64(macro Token, u uint64) Value {
	return Value{macro, "", u, typeUint, nil}
}

// Uint32 creates a new value that replaces `macro` with a `uint32`
func Uint32(macro Token, u uint32) Value {
	return Value{macro, "", uint64(u), typeUint, nil}
}

// Uint16 creates a new value that replaces `macro` with a `uint16`
func Uint16(macro Token, u uint16) Value {
	return Value{macro, "", uint64(u), typeUint, nil}
}

// Uint8 creates a new value that replaces `macro` with a `uint8`
func Uint8(macro Token, u uint8) Value {
	return Value{macro, "", uint64(u), typeUint, nil}
}

// Int creates a new value that replaces `macro` with an `int`
func Int(macro Token, i int) Value {
	return Value{macro, "", uint64(int64(i)), typeInt, nil}
}

// Int64 creates a new value that replaces `macro` with an `int64`
func Int64(macro Token, i int64) Value {
	return Value{macro, "", uint64(i), typeInt, nil}
}

// Int32 creates a new value that replaces `macro` with an `int32`
func Int32(macro Token, i int32) Value {
	return Value{macro, "", uint64(int64(i)), typeInt, nil}
}

// Int16 creates a new value that replaces `macro` with an `int16`
func Int16(macro Token, i int16) Value {
	return Value{macro, "", uint64(int64(i)), typeInt, nil}
}

// Int8 creates a new value that replaces `macro` with an `int8`
func Int8(macro Token, i int8) Value {
	return Value{macro, "", uint64(int64(i)), typeInt, nil}
}

// Unix crates a new value that replaces `macro` with the unix timestamp of `tm`
func Unix(macro Token, tm time.Time) Value {
	return Value{macro, "", uint64(tm.Unix()), typeInt, nil}
}

// Time creates a new value that replaces `macro` with `tm` formatted according to `layout`
func Time(macro Token, tm time.Time, layout string) Value {
	return Value{macro, layout, 0, typeTime, tm}
}

// Any creates a new value that replaces `macro` with any value
func Any(macro Token, x interface{}) Value {
	if any, ok := x.(ValueAppender); ok {
		return Value{macro, "", 0, typeAny, any}
	}
	return Value{macro, "", 0, typeAny, any{x}}
}

// Bind creates a new value that replaces `macro` with any value
func Bind(macro Token, v ValueAppender) Value {
	return Value{macro, "", 0, typeAny, v}
}

// AppendValue implements `ValueAppender` interface
func (v *Value) AppendValue(buf []byte) ([]byte, error) {
	switch v.typ {
	case typeString:
		return append(buf, v.str...), nil
	case typeFloat:
		f := math.Float64frombits(v.num)
		return strconv.AppendFloat(buf, f, 'f', -1, 64), nil
	case typeUint:
		return strconv.AppendUint(buf, v.num, 10), nil
	case typeInt:
		return strconv.AppendInt(buf, int64(v.num), 10), nil
	case typeTime:
		return v.any.(time.Time).AppendFormat(buf, v.str), nil
	case typeConcat:
		values := v.any.([]string)
		sep := v.str
		for i, v := range values {
			if i > 0 {
				buf = append(buf, sep...)
			}
			buf = append(buf, v...)
		}
		return buf, nil
	case typeAny:
		if v, ok := v.any.(ValueAppender); ok {
			return v.AppendValue(buf)
		}
		return any{v.any}.AppendValue(buf)
	case typeNone:
		return buf, ErrMacroNotFound
	default:
		return nil, errors.New("Invalid value type")
	}
}

// ValueAppender appends a value to a buffer
type ValueAppender interface {
	AppendValue([]byte) ([]byte, error)
}

type any struct {
	value interface{}
}

func (v any) AppendValue(buf []byte) ([]byte, error) {
	switch v := v.value.(type) {
	case string:
		return append(buf, v...), nil
	case []byte:
		return append(buf, v...), nil
	case int:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case int8:
		return strconv.AppendInt(buf, int64(v), 10), nil
	case uint:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case uint8:
		return strconv.AppendUint(buf, uint64(v), 10), nil
	case float64:
		return strconv.AppendFloat(buf, float64(v), 'f', -1, 64), nil
	case float32:
		return strconv.AppendFloat(buf, float64(v), 'f', -1, 32), nil
	case bool:
		if v {
			return append(buf, "true"...), nil
		}
		return append(buf, "false"...), nil
	case fmt.Stringer:
		return append(buf, v.String()...), nil
	default:
		s := fmt.Sprintf("%s", v)
		return append(buf, s...), nil
	}
}
