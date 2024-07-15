package ConflictSolver

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type HelperElem[T gr.TokenTyper] interface {
	// SetLookahead sets the lookahead of the action.
	//
	// Parameters:
	//   - lookahead: The lookahead to set.
	SetLookahead(lookahead *T)

	// AppendRhs appends a symbol to the right-hand side of the action.
	//
	// Parameters:
	//   - symbol: The symbol to append.
	AppendRhs(symbol T)

	Actioner[T]

	fmt.Stringer
	uc.Copier
}
