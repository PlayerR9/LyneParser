package ConflictSolver

import "fmt"

type ConflictSolver struct {
	Elements []*Helper
}

func NewConflictSolver() *ConflictSolver {
	return &ConflictSolver{
		Elements: make([]*Helper, 0),
	}
}

func solveMinimum(elements []*Helper) error {
	// 1. Bucket sort the items by their position
	buckets := make(map[int][]*Item)

	for _, elem := range elements {
		buckets[elem.Item.Pos] = append(buckets[elem.Item.Pos], elem.Item)
	}

	for limit, bucket := range buckets {
		err := minimum(bucket, solutions, limit)
		if err != nil {
			return err
		}
	}

	// Now, find conflicts between buckets

	return nil
}

// Find the minimum number of elements required to uniquely identify
// any item.
func minimum(items []*Item, solutions []Actioner, limit int) error {
	doneMap := make(map[*Item]bool)

	for _, item := range items {
		doneMap[item] = false
	}

	for i := limit; i >= 0; i-- {
		rhsPerLevel := make(map[string][]int)

		for j, item := range items {
			done, ok := doneMap[item]
			if !ok {
				return fmt.Errorf("item %v not found in doneMap", item)
			} else if done {
				continue
			}

			rhs, err := item.Rule.GetRhsAt(i)
			if err != nil {
				return err
			}

			rhsPerLevel[rhs] = append(rhsPerLevel[rhs], j)
		}

		for rhs, indices := range rhsPerLevel {
			if len(indices) == 1 {
				// No conflict
				err := solutions[indices[0]].AppendRhs(rhs)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func solveSRConflict(symbol string, items []*Item) ([]Actioner, error) {
	itemMap := make(map[*Item]bool)

	for _, item := range items {
		itemMap[item] = false
	}

	results := make([]Actioner, 0, len(items))

	// 1. Assign a simple action for each item
	if symbol == gr.EOFTokenID {
		for _, item := range items {
			if item.IsReduce {
				results = append(results, NewAcceptAction(item.ruleIndex))
			} else {
				results = append(results, NewActShift())
			}
		}
	} else {
		for _, item := range items {
			if item.IsReduce {
				results = append(results, NewActReduce(item.ruleIndex))
			} else {
				results = append(results, NewActShift())
			}
		}
	}

	// 2. For each shift action, check if it has a look-ahead.
	for i, result := range results {
		act, ok := result.(*ActShift)
		if !ok {
			continue
		}

		lookahead, err := items[i].Rule.GetRhsAt(items[i].Pos + 1)
		if err == nil && gr.IsTerminal(lookahead) {
			act.SetLookahead(&lookahead)
		}
	}

	// Now, those that have the lookahead are no longer in conflict with
	// their reduce counterparts. But there still might be conflicts between
	// shift actions that has the same lookahead.

	for i := 0; i < len(results); i++ {
		act, ok := results[i].(*ActShift)
		if !ok || act.Lookahead == nil {
			continue
		}

		conflicts := make([]int, 0)

		for j := i + 1; j < len(results); j++ {
			act, ok := results[i].(*ActShift)
			if !ok || act.Lookahead == nil {
				continue
			}

			if *act.Lookahead == *results[j].(*ActShift).Lookahead {
				conflicts = append(conflicts, j)
			}
		}

		if len(conflicts) == 0 {
			itemMap[items[i]] = true
		}

		// Solve conflicts
	}

	for _, item := range items {
		if item.Pos == 0 {
			results = append(results, &hp.HResult[*Item]{
				Result: item,
				Reason: nil,
			})

			continue
		}

		var reason error = nil

		for i := item.Pos - 1; i >= 0; i-- {
			rhs, err := item.Rule.GetRhsAt(i)
			if err != nil {
				reason = fmt.Errorf("could not get RHS at index %d", i)
				break
			} else if stack.IsEmpty() {
				reason = ers.NewErrUnexpected(nil, rhs)
				break
			}

			if top := stack.Pop(); top.GetID() != rhs {
				reason = ers.NewErrUnexpected(top, rhs)
				break
			}
		}

		results = append(results, &Helper{Elem: item, Reason: reason})

		stack.Refuse()
	}

	success := make([]*Helper, 0)
	fail := make([]*Helper, 0)

	for _, r := range results {
		if r.Reason == nil {
			success = append(success, r)
		} else {
			fail = append(fail, r)
		}
	}

	if len(success) == 0 {
		if len(shifts) == 0 {
			// Return the most likely error
			// As of now, we will return the first error
			return NewErrorAction(fail[0].Reason)
		}

		// We can only shift
		return NewShiftAction()
	}

	// Find the actual reduce item

	weights := slext.ApplyWeightFunc(success, func(h *Helper) (float64, bool) {
		return float64(h.Elem.Pos), true
	})

	final := slext.FilterByPositiveWeight(weights)

	if len(final) == 1 {
		return NewReduceAction(final[0].Elem.ruleIndex)
	}

	// AMBIGUOUS GRAMMAR

	// SHIFT-REDUCE CONFLICT
}

func (dt *DecisionTable) FixConflicts() error {
	dt.actions = make(map[string][]*MatcherAction)

	possibleSRConflict := make(map[string][]*Item)

	for symbol, items := range dt.table {
		// Split the items into shifts and reduces.
		shifts, reduces := SplitShiftReduce(items)

		// If there are no reduce items, then we can only shift.
		if len(reduces) == 0 {
			if len(shifts) == 0 {
				return fmt.Errorf("no actions found for symbol %s", symbol)
			}

			dt.actions[symbol] = []*MatcherAction{NewMatcherAction(NewShiftAction())}
		} else {
			possibleSRConflict[symbol] = items
		}
	}

	if len(possibleSRConflict) == 0 {
		// No shift-reduce conflicts
		return nil
	}

	// For each rule, find the least amount of elements required to uniquely identify
	// any item.

	for symbol, items := range possibleSRConflict {
		err := solveSRConflict(symbol, items)
		if err != nil {
			return fmt.Errorf("could not solve shift-reduce conflict for symbol %s: %s",
				symbol, err.Error(),
			)
		}
	}
}
