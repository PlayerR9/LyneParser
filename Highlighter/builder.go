package Highlighter

import (
	"github.com/gdamore/tcell"
)

// Builder is a builder for Highlighter.
type Builder struct {
	// rules is a map of rules to apply.
	rules map[string]tcell.Style
}

// AddRule is a method of Builder that adds a rule to the builder.
// It overrides any existing rules with the same ID.
//
// Parameters:
//   - style: The style to apply.
//   - ids: The IDs to apply the style to.
func (b *Builder) AddRule(style tcell.Style, ids ...string) {
	if b.rules == nil {
		b.rules = make(map[string]tcell.Style)
	}

	for _, id := range ids {
		b.rules[id] = style
	}
}

// Build is a method of Builder that builds a Highlighter from the rules.
// It resets the builder after building the Highlighter.
//
// Returns:
//   - Highlighter: The new Highlighter.
func (b *Builder) Build() Highlighter {
	h := Highlighter{
		rules: make(map[string]tcell.Style),
	}

	for id, style := range b.rules {
		h.rules[id] = style
	}

	b.rules = nil

	return h
}

// Reset is a method of Builder that resets the builder.
func (b *Builder) Reset() {
	b.rules = nil
}
