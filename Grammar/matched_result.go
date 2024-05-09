package Grammar

// MatchedResult represents the result of a match operation.
type MatchedResult[T Tokener] struct {
	// Matched is the matched token.
	Matched T

	// RuleIndex is the index of the production that matched.
	RuleIndex int
}

// GetMatch returns the matched token.
//
// Returns:
//   - T: The matched token.
func (mr *MatchedResult[T]) GetMatch() T {
	return mr.Matched
}

// NewMatchResult is a constructor of MatchedResult.
//
// Parameters:
//   - matched: The matched token.
//   - ruleIndex: The index of the production that matched.
//
// Returns:
//   - *MatchedResult: A new MatchedResult.
func NewMatchResult[T Tokener](matched T, ruleIndex int) *MatchedResult[T] {
	return &MatchedResult[T]{Matched: matched, RuleIndex: ruleIndex}
}
