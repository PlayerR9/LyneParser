package LyneLexer

import (
	"fmt"
	"regexp"
	"slices"
	"unicode"

	Stack "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Stack"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

type RuleOption func(*Rule) error

func WithSkip() RuleOption {
	return func(r *Rule) error {
		r.Skip = true

		return nil
	}
}

type Rule struct {
	Name  string
	Regex string
	*regexp.Regexp
	Skip bool
}

func (r *Rule) String() string {
	if r == nil {
		return "Rule[nil]"
	}

	if r.Regexp == nil {
		if r.Skip {
			return fmt.Sprintf("Rule[Name=%s, Regex=%v, Compiled=N/A, Skip=Y]", r.Name, r.Regex)
		} else {
			return fmt.Sprintf("Rule[Name=%s, Regex=%v, Compiled=N/A, Skip=N]", r.Name, r.Regex)
		}
	}

	if r.Skip {
		return fmt.Sprintf("Rule[Name=%s, Regex=%v, Compiled=%v, Skip=Y]", r.Name, r.Regex, r.Regexp)
	} else {
		return fmt.Sprintf("Rule[Name=%s, Regex=%v, Compiled=%v, Skip=N]", r.Name, r.Regex, r.Regexp)
	}
}

func NewRule(name string, regex string, options ...RuleOption) (*Rule, error) {
	if !IsLexToken(name) {
		return nil, ers.NewErrInvalidParameter("name").
			Wrap(NewErrInvalidTokenName(name))
	}

	regexp, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	r := &Rule{
		Name:   name,
		Regex:  regex,
		Regexp: regexp,
	}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Rule) Equals(other *Rule) bool {
	if r == nil && other == nil {
		return true
	}

	if r == nil || other == nil {
		return false
	}

	return r.Name == other.Name && r.Regex == other.Regex
}

type RuleSetBuilder struct {
	rules []*Rule
}

func (b *RuleSetBuilder) Insert(rule *Rule) {
	if rule == nil {
		return
	}

	if b.rules == nil {
		b.rules = make([]*Rule, 0)
	}

	b.rules = append(b.rules, rule)
}

func (b *RuleSetBuilder) AddRule(name string, regex string, options ...RuleOption) error {
	if b.rules == nil {
		b.rules = make([]*Rule, 0)
	}

	r, err := NewRule(name, regex)
	if err != nil {
		return err
	}

	for _, option := range options {
		if err := option(r); err != nil {
			return err
		}
	}

	b.rules = append(b.rules, r)

	return nil
}

func (b *RuleSetBuilder) Build() []*Rule {
	if len(b.rules) == 0 {
		return nil
	}

	// Remove duplicates
	for i := 0; i < len(b.rules); i++ {
		for j := i + 1; j < len(b.rules); {
			if b.rules[i].Equals(b.rules[j]) {
				b.rules = append(b.rules[:j], b.rules[j+1:]...)
			} else {
				j++
			}
		}
	}

	result := make([]*Rule, len(b.rules))
	copy(result, b.rules)

	// Clear the builder
	for i := range b.rules {
		b.rules[i] = nil
	}

	b.rules = nil

	return result
}

func (b *RuleSetBuilder) Clear() {
	for i := range b.rules {
		b.rules[i] = nil
	}

	b.rules = nil
}

type Lexer struct {
	rules  []*Rule
	toSkip []string

	roots  []*matchNode
	leaves []*matchNode
}

func NewLexer(rules []*Rule) *Lexer {
	ruleSet := make([]*Rule, len(rules))
	copy(ruleSet, rules)

	// Sort the rules by length of the regex (longest first)
	slices.SortFunc(ruleSet, func(lr, rr *Rule) int {
		return len(lr.Regex) - len(rr.Regex)
	})

	l := &Lexer{
		rules:  ruleSet,
		toSkip: make([]string, 0),
		roots:  make([]*matchNode, 0),
		leaves: make([]*matchNode, 0),
	}

	for _, rule := range l.rules {
		if rule.Skip {
			l.toSkip = append(l.toSkip, rule.Name)
		}
	}

	return l
}

func (l *Lexer) deleteBranch(node *matchNode) error {
	parent, branchNode := node.findBranchingPoint()
	if parent == nil || branchNode == nil {
		return fmt.Errorf("no branching point found")
	}

	// Dereference all the children
	branchNode.removeBranch()

	var prev *matchNode = nil

	for node := parent.firstChild; node != nil && node != branchNode; node = node.nextSibling {
		prev = node
	}

	if prev == nil {
		parent.firstChild = branchNode.nextSibling
	} else {
		prev.nextSibling = branchNode.nextSibling
	}

	return nil
}

func (l *Lexer) findMatches(input *[]byte, at int) []*Token {
	subInput := (*input)[at:]

	tokens := make([]*Token, 0, len(l.rules))

	for i, r := range l.rules {
		matched := r.Find(subInput)

		if len(matched) == 0 {
			continue
		}

		tk, _ := NewToken(r.Name, matched, i)
		tokens = append(tokens, tk)
	}

	if len(tokens) == 0 {
		return nil
	}

	return tokens
}

func (l *Lexer) Tokenize(input []byte) error {
	l.roots = make([]*matchNode, 0)
	l.leaves = make([]*matchNode, 0)

	matched := l.findMatches(&input, 0)
	if len(matched) == 0 {
		return fmt.Errorf("no match found at position 0")
	}

	for _, match := range matched {
		node := newMatchNode(match, 0)

		l.roots = append(l.roots, node)
		l.leaves = append(l.leaves, node)
	}

	for {
		leavesToProcess := make([]*matchNode, 0, len(l.leaves))

		for _, leaf := range l.leaves {
			if !leaf.isDone {
				leavesToProcess = append(leavesToProcess, leaf)
			}
		}

		if len(leavesToProcess) == 0 {
			break
		}

		for _, leaf := range leavesToProcess {
			nextPos := leaf.pos + len(leaf.token.GetData())

			if nextPos >= len(input) {
				leaf.isDone = true
				continue
			}

			matched := l.findMatches(&input, nextPos)
			if len(matched) == 0 {
				// Delete branch
				err := l.deleteBranch(leaf)
				if err != nil {
					return err
				}

				continue
			}

			nodes := make([]*matchNode, 0)

			for _, match := range matched {
				if !slices.Contains(l.toSkip, match.GetID()) {
					nodes = append(nodes, newMatchNode(match, nextPos))
				}
			}

			leaf.addChildren(nodes...)
		}

		// Update the leaves
		newLeaves := make([]*matchNode, 0)

		for _, leaf := range l.leaves {
			newLeaves = append(newLeaves, leaf.getLeaves()...)
		}

		l.leaves = newLeaves
	}

	return nil
}

func LexString(input string, rules []*Rule) ([]Stack.Stacker[*Token], error) {
	lexer := NewLexer(rules)

	err := lexer.Tokenize([]byte(input))
	if err != nil {
		return nil, err
	}

	return lexer.GetTokens(), nil
}

func (l *Lexer) GetTokens() []Stack.Stacker[*Token] {
	results := make([][]*Token, 0)

	for _, root := range l.roots {
		for _, branches := range root.snakeTraversal() {
			tmp := make([]*Token, 0, len(branches))

			for _, branch := range branches {
				tmp = append(tmp, branch.token)
			}

			slices.Reverse(tmp)

			results = append(results, tmp)
		}
	}

	inputStreams := make([]Stack.Stacker[*Token], len(results))

	for i, branches := range results {
		inputStreams[i] = Stack.NewLinkedStack(branches...)
	}

	return inputStreams
}

func IsLexToken(tokenType string) bool {
	if tokenType == "" {
		return false
	}

	firstChar := []rune(tokenType)[0]

	return unicode.IsLetter(firstChar) && unicode.IsLower(firstChar)
}
