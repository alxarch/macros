package macros

import (
	"errors"
	"math"
)

// MacroReplacer provides a value to replace a macro token.
//
// If no replacement value exists for this token a `None()`
// value should be returned
type MacroReplacer interface {
	ReplaceMacro(macro string) Value
}

// ErrMacroNotFound is the error to return when a macro is not found
var ErrMacroNotFound = errors.New("Macro not found")

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
	fieldTypeTemplate
)

// Value is a value replacement for a macro
type Value struct {
	macro Token
	str   string
	num   uint64
	typ   fieldType
	any   interface{}
}

// String creates a new MacroValue field for a string value
func String(macro Token, s string) Value {
	return Value{macro, s, 0, fieldTypeString, nil}
}

// Bool creates a new MacroValue field for a bool value
func Bool(macro Token, v bool) Value {
	if v {
		return Value{macro, "true", 0, fieldTypeString, nil}
	}
	return Value{macro, "false", 0, fieldTypeString, nil}
}

// Float creates a new MacroValue field for a float value
func Float(macro Token, f float64) Value {
	return Value{macro, "", math.Float64bits(f), fieldTypeFloat, nil}
}

// Uint creates a new MacroValue field for a uint value
func Uint(macro Token, u uint) Value {
	return Value{macro, "", uint64(u), fieldTypeUint, nil}
}

// Uint64 creates a new MacroValue field for a uint64 value
func Uint64(macro Token, u uint64) Value {
	return Value{macro, "", u, fieldTypeUint, nil}
}

// Uint32 creates a new MacroValue field for a uint64 value
func Uint32(macro Token, u uint32) Value {
	return Value{macro, "", uint64(u), fieldTypeUint, nil}
}

// Uint16 creates a new MacroValue field for a uint64 value
func Uint16(macro Token, u uint16) Value {
	return Value{macro, "", uint64(u), fieldTypeUint, nil}
}

// Uint8 creates a new MacroValue field for a uint64 value
func Uint8(macro Token, u uint8) Value {
	return Value{macro, "", uint64(u), fieldTypeUint, nil}
}

// Int creates a new MacroValue field for an int value
func Int(macro Token, i int) Value {
	return Value{macro, "", uint64(int64(i)), fieldTypeInt, nil}
}

// Int64 creates a new MacroValue field for a uint64 value
func Int64(macro Token, i int64) Value {
	return Value{macro, "", uint64(i), fieldTypeInt, nil}
}

// Int32 creates a new MacroValue field for a uint64 value
func Int32(macro Token, i int32) Value {
	return Value{macro, "", uint64(int64(i)), fieldTypeInt, nil}
}

// Int16 creates a new MacroValue field for a uint64 value
func Int16(macro Token, i int16) Value {
	return Value{macro, "", uint64(int64(i)), fieldTypeInt, nil}
}

// Int8 creates a new MacroValue field for a uint64 value
func Int8(macro Token, i int8) Value {
	return Value{macro, "", uint64(int64(i)), fieldTypeInt, nil}
}

// Any creates a new MacroValue field for any value
func Any(macro Token, x interface{}) Value {
	return Value{macro, "", 0, fieldTypeAny, x}
}

// Nested creates a new value that renders a template
func Nested(macro Token, tpl *Template) Value {
	return Value{macro, "", 0, fieldTypeTemplate, tpl}
}
