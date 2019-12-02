package macros

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MacroReplacer provides a value to replace a macro token.
//
// If no replacement value exists for this token a `None()`
// value should be returned
type MacroReplacer interface {
	ReplaceMacro(macro string) Value
}

// MacroValues is a list of values bound to a macro
type MacroValues []MacroValue

// ErrMacroNotFound is the error to return when a macro is not found
var ErrMacroNotFound = errors.New("Macro not found")

// ReplaceMacro implements `MacroReplacer` interface
func (values MacroValues) ReplaceMacro(macro string) Value {
	for i := range values {
		if values[i].macro == macro {
			return values[i].Value
		}
	}
	return None()
}

// MacroReplacerFunc is a closure implementing `MacroReplacer` interface
type MacroReplacerFunc func(macro string) Value

// ReplaceMacro implements `MacroReplacer` interface
func (f MacroReplacerFunc) ReplaceMacro(macro string) Value {
	return f(macro)
}

type fieldType uint

const (
	fieldTypeNone fieldType = iota
	fieldTypeString
	fieldTypeFloat
	fieldTypeInt
	fieldTypeUint
	fieldTypeAny
)

// Value is a value replacement for a macro
type Value struct {
	str string
	num uint64
	typ fieldType
	any interface{}
}

// MacroValue is a value replacement bound to a macro token
type MacroValue struct {
	macro string
	Value
}

// IsNone checks if a value is None value
func (v *Value) IsNone() bool {
	return v.typ == fieldTypeNone
}

func (v *Value) appendTo(buf []byte) []byte {
	switch v.typ {
	case fieldTypeString:
		return append(buf, v.str...)
	case fieldTypeFloat:
		f := math.Float64frombits(v.num)
		return strconv.AppendFloat(buf, f, 'f', -1, 64)
	case fieldTypeUint:
		return strconv.AppendUint(buf, v.num, 10)
	case fieldTypeInt:
		return strconv.AppendInt(buf, int64(v.num), 10)
	case fieldTypeAny:
		var w strings.Builder
		fmt.Fprintf(&w, "%s", v.any)
		return append(buf, w.String()...)
	default:
		return buf
	}
}

// String creates a new MacroValue field for a string value
func String(s string) Value {
	return Value{s, 0, fieldTypeString, nil}
}

// Bool creates a new MacroValue field for a bool value
func Bool(v bool) Value {
	if v {
		return Value{"true", 0, fieldTypeString, nil}
	}
	return Value{"false", 0, fieldTypeString, nil}
}

// Float creates a new MacroValue field for a float value
func Float(f float64) Value {
	return Value{"", math.Float64bits(f), fieldTypeFloat, nil}
}

// Uint creates a new MacroValue field for a uint value
func Uint(u uint) Value {
	return Value{"", uint64(u), fieldTypeUint, nil}
}

// Uint64 creates a new MacroValue field for a uint64 value
func Uint64(u uint64) Value {
	return Value{"", u, fieldTypeUint, nil}
}

// Uint32 creates a new MacroValue field for a uint64 value
func Uint32(u uint32) Value {
	return Value{"", uint64(u), fieldTypeUint, nil}
}

// Uint16 creates a new MacroValue field for a uint64 value
func Uint16(u uint16) Value {
	return Value{"", uint64(u), fieldTypeUint, nil}
}

// Uint8 creates a new MacroValue field for a uint64 value
func Uint8(u uint8) Value {
	return Value{"", uint64(u), fieldTypeUint, nil}
}

// Int creates a new MacroValue field for an int value
func Int(i int) Value {
	return Value{"", uint64(int64(i)), fieldTypeInt, nil}
}

// Int64 creates a new MacroValue field for a uint64 value
func Int64(i int64) Value {
	return Value{"", uint64(i), fieldTypeInt, nil}
}

// Int32 creates a new MacroValue field for a uint64 value
func Int32(i int32) Value {
	return Value{"", uint64(int64(i)), fieldTypeInt, nil}
}

// Int16 creates a new MacroValue field for a uint64 value
func Int16(i int16) Value {
	return Value{"", uint64(int64(i)), fieldTypeInt, nil}
}

// Int8 creates a new MacroValue field for a uint64 value
func Int8(i int8) Value {
	return Value{"", uint64(int64(i)), fieldTypeInt, nil}
}

// Any creates a new MacroValue field for any value
func Any(x interface{}) Value {
	return Value{"", 0, fieldTypeAny, x}
}

// None creates a new value for no value found
func None() Value {
	return Value{"", 0, fieldTypeNone, nil}
}

// Values creates a `MacroValues` from various MacroValues
func Values(values ...MacroValue) MacroValues {
	return MacroValues(values)
}

// Bind adds a value binding to a `MacroValues`
func (values MacroValues) Bind(macro string, value Value) MacroValues {
	return append(values, MacroValue{macro: macro, Value: value})
}

// Bind binds a value to a macro token
func Bind(macro string, value Value) MacroValue {
	return MacroValue{
		macro: macro,
		Value: value,
	}
}

// ValueMap is a `MacroReplacer` that maps macros to values
type ValueMap map[string]Value

// ReplaceMacro implements `MacroReplacer` interface
func (m ValueMap) ReplaceMacro(macro string) Value {
	return m[macro]
}
