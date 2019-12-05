package macros

import (
	"fmt"
	"strings"
)

// Option is a Replacer option
type Option interface {
	apply(r *Replacer) error
}

type optionFunc func(*Replacer) error

func (option optionFunc) apply(r *Replacer) error {
	return option(r)
}

// Delimiters sets Replacer delimiters
func Delimiters(start, end string) Option {
	return optionFunc(func(p *Replacer) error {
		start = strings.TrimSpace(start)
		if start == "" {
			return fmt.Errorf("Invalid start delimiter")
		}
		end = strings.TrimSpace(end)
		if end == "" {
			return fmt.Errorf("Invalid end delimiter")
		}
		p.start, p.end = start, end
		return nil
	})
}

// Alias defines aliases for a macro
func Alias(macro Token, aliases ...Token) Option {
	return optionFunc(func(p *Replacer) error {
		if p.alias == nil {
			p.alias = make(map[Token]Token)
		}
		macro, _ = macro.split()
		for _, alias := range aliases {
			alias, _ = alias.split()
			p.alias[macro] = alias
		}
		return nil
	})
}

// DefaultValue sets a value to be used when no macro replacement is found
func DefaultValue(value string) Option {
	return optionFunc(func(p *Replacer) error {
		p.none = String("", value)
		return nil
	})
}

// Skip defines tokens that will not be replaced
func Skip(tokens ...Token) Option {
	return optionFunc(func(p *Replacer) error {
		if p.skip == nil {
			p.skip = make(map[Token]struct{}, len(tokens))
		}
		for _, token := range tokens {
			p.skip[token] = struct{}{}
		}
		return nil
	})

}
