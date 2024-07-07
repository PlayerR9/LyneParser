package Grammar

// MatchedResult represents the result of a match operation.
type MatchedResult[T TokenTyper] struct {
	// Matched is the matched token.
	Matched *Token[T]

	// RuleIndex is the index of the production that matched.
	RuleIndex int
}

// NewMatchResult is a constructor of MatchedResult.
//
// Parameters:
//   - matched: The matched token.
//   - ruleIndex: The index of the production that matched.
//
// Returns:
//   - *MatchedResult: A new MatchedResult.
func NewMatchResult[T TokenTyper](matched *Token[T], ruleIndex int) *MatchedResult[T] {
	mr := &MatchedResult[T]{
		Matched:   matched,
		RuleIndex: ruleIndex,
	}

	return mr
}
