package macros

import (
	"strings"
)

// Option is a Replacer option
type Option interface {
	apply(r *Replacer)
}

type optionFunc func(*Replacer)

func (option optionFunc) apply(r *Replacer) {
	option(r)
}

// Delimiters sets Replacer delimiters
func Delimiters(start, end string) Option {
	return optionFunc(func(p *Replacer) {
		start = strings.TrimSpace(start)
		if start == "" {
			start = defaultStartDelimiter
		}
		end = strings.TrimSpace(end)
		if end == "" {
			end = defaultEndDelimiter
		}
		p.start, p.end = start, end
	})
}

// Alias defines aliases for a macro
func Alias(macro Token, aliases ...Token) Option {
	return optionFunc(func(p *Replacer) {
		if p.alias == nil {
			p.alias = make(map[Token]Token)
		}
		macro, _ = macro.split()
		for _, alias := range aliases {
			alias, _ = alias.split()
			p.alias[alias] = macro
		}
	})
}

// DefaultValue sets a value to be used when no macro replacement is found
func DefaultValue(value string) Option {
	return optionFunc(func(p *Replacer) {
		p.none = String("", value)
	})
}

// Skip defines macros that will not be replaced
func Skip(macros ...Token) Option {
	return optionFunc(func(p *Replacer) {
		if p.skip == nil {
			p.skip = make(map[Token]struct{}, len(macros))
		}
		for _, token := range macros {
			p.skip[token] = struct{}{}
		}
	})

}

type expand struct {
	macro Token
	tpl   string
}

func (e expand) apply(r *Replacer) {
	if r.expand == nil {
		r.expand = make(map[Token]string)
	}
	r.expand[e.macro] = e.tpl
}

// Expand creates a new value that replaces `macro` by expanding `tpl`
func Expand(macro Token, tpl string) Option {
	macro, _ = macro.split()
	return expand{macro, tpl}
}
