package ConflictSolver

import (
	"fmt"

	gr "github.com/PlayerR9/LyneParser/Grammar"
)

type ConflictSolver struct {
	Elements []*Helper
}

func NewConflictSolver(leadingSymbol string, items []*Item) *ConflictSolver {
	helpers := make([]*Helper, 0, len(items))

	var act Actioner

	if leadingSymbol == gr.EOFTokenID {
		for _, item := range items {
			if item.IsReduce {
				act = NewAcceptAction(item.ruleIndex)
			} else {
				act = NewActShift()
			}

			helpers = append(helpers, NewHelper(item, act))
		}
	} else {
		for _, item := range items {
			if item.IsReduce {
				act = NewActReduce(item.ruleIndex)
			} else {
				act = NewActShift()
			}

			helpers = append(helpers, NewHelper(item, act))
		}
	}

	return &ConflictSolver{
		Elements: helpers,
	}
}

func solveMinimum(elems []*Helper) error {
	// 1. Bucket sort the items by their position.
	buckets := make(map[int][]*Helper)

	for _, elem := range elems {
		buckets[elem.Item.Pos] = append(buckets[elem.Item.Pos], elem)
	}

	for limit, bucket := range buckets {
		err := minimum(bucket, limit)
		if err != nil {
			return err
		}
	}

	// Now, find conflicts between buckets.
	// FIXME: Solve conflicts between buckets.

	return nil
}

// Find the minimum number of elements required to uniquely identify
// any item.
func minimum(helpers []*Helper, limit int) error {
	// Set all helpers to not done.
	elemsEval := make(map[*Helper]bool)

	for _, h := range helpers {
		elemsEval[h] = false
	}

	for i := limit; i >= 0; i-- {
		rhsPerLevel := make(map[string][]int)

		for j, h := range helpers {
			isDone, ok := elemsEval[h]
			if !ok {
				panic(fmt.Sprintf("item %v not found in doneMap", h))
			} else if isDone {
				continue
			}

			rhs, err := h.Item.Rule.GetRhsAt(i)
			if err != nil {
				panic(err)
			}

			rhsPerLevel[rhs] = append(rhsPerLevel[rhs], j)
		}

		for rhs, indices := range rhsPerLevel {
			currentH := helpers[indices[0]]

			// Add the RHS.
			err := currentH.Action.AppendRhs(rhs)
			if err != nil {
				panic(err)
			}

			if len(indices) == 1 {
				// No conflict. Mark it as done.
				elemsEval[currentH] = true
			}
		}
	}

	// FIXME: Check if it was successful

	return nil
}

func (cs *ConflictSolver) SolveConflicts() error {
	// Set all helpers to not done.
	elemsEval := make(map[*Helper]bool)

	for _, elem := range cs.Elements {
		elemsEval[elem] = false
	}

	// 1. FOR EACH shift action, check if it has a look-ahead.
	for _, elem := range cs.Elements {
		act, ok := elem.Action.(*ActShift)
		if !ok {
			continue
		}

		item := elem.Item

		lookahead, err := item.Rule.GetRhsAt(item.Pos + 1)
		if err == nil && gr.IsTerminal(lookahead) {
			act.SetLookahead(&lookahead)
		}
	}

	// Now, those shift actions that have the look-ahead are no longer
	// in conflict with their reduce counterparts.
	// However, there still might be conflicts between shift actions
	// with the same look-ahead.

	laConflicts := make(map[string][]*Helper)

	for _, elem := range cs.Elements {
		act, ok := elem.Action.(*ActShift)
		if !ok || act.Lookahead == nil {
			continue
		}

		lookahead := *act.Lookahead

		laConflicts[lookahead] = append(laConflicts[lookahead], elem)
	}

	for _, elems := range laConflicts {
		if len(elems) == 1 {
			// Mark as done.
			elemsEval[elems[0]] = true
		} else {
			// Solve conflicts.
			err := solveMinimum(elems)
			if err != nil {
				return fmt.Errorf("ambiguous grammar")
			}

			for _, elem := range elems {
				elemsEval[elem] = true
			}
		}
	}

	// FIXME:

	// AMBIGUOUS GRAMMAR

	// SHIFT-REDUCE CONFLICT

	return nil
}
