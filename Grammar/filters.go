package Grammar

import "slices"

// FilterProductionsWithoutLHS filters out productions without the specified
// left-hand side.
//
// Parameters:
//   - lhs: The left-hand side to filter.
//
// Returns:
//   - bool: True if the production has the specified left-hand side, false
//     otherwise.
func (b *GrammarBuilder) FilterProductionsWithoutLHS(lhs string) bool {
	filterProductionWithLHS := func(p Productioner) bool {
		return p != nil && p.GetLhs() == lhs
	}

	return slices.ContainsFunc(b.productions, filterProductionWithLHS)
}
