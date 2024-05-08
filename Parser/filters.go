package Parser

import (
	cs "github.com/PlayerR9/LyneParser/ConflictSolver"
	hlp "github.com/PlayerR9/MyGoLib/CustomData/Helpers"
)

/////////////////////////////////////////////////////////////

// HResultWeightFunc is a function that, given a HResult, returns the weight of the result.
//
// Parameters:
//   - h: The HResult to calculate the weight of.
//
// Returns:
//   - float64: The weight of the result.
//   - bool: True if the weight is valid, false otherwise.
func HResultWeightFunc(h hlp.HResult[cs.Actioner]) (float64, bool) {
	return float64(h.First.Size()), true
}
