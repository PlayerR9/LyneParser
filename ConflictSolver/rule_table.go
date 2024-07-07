package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// RuleTable represents a table of items.
type RuleTable[T gr.TokenTyper] struct {
	// items is the items of the rule table.
	items []*Item[T]

	// buckets is the buckets of the rule table.
	buckets map[T][]*Helper[T]
}

// NewRuleTable is a constructor of RuleTable.
//
// Parameters:
//   - symbols: The symbols to use.
//   - rules: The rules to use.
//
// Returns:
//   - *RuleTable: The new rule table.
func NewRuleTable[T gr.TokenTyper](symbols []T, rules []*gr.Production[T]) *RuleTable[T] {
	rt := &RuleTable[T]{
		items: make([]*Item[T], 0),
	}

	for _, s := range symbols {
		for _, r := range rules {
			indices := r.IndicesOfRhs(s)

			for _, i := range indices {
				item, err := NewItem(r, i, len(rt.items))
				uc.AssertF(err == nil, "NewItem failed: %s", err)

				rt.items = append(rt.items, item)
			}
		}
	}

	rt.buckets = rt.get_item_buckets()

	return rt
}

// get_item_buckets gets the item buckets of the rule table.
//
// Returns:
//   - map[T]*uts.Bucket[*Helper]: The item buckets.
func (rt *RuleTable[T]) get_item_buckets() map[T][]*Helper[T] {
	buckets := make(map[T][]*Helper[T])

	for _, item := range rt.items {
		symbol, err := item.Rule.GetRhsAt(item.Pos)
		uc.AssertF(err == nil, "GetRhsAt failed: %s", err)

		last_index := item.Rule.Size() - 1

		var act HelperElem[T]

		if item.Pos == last_index {
			if symbol.String() == gr.EOFTokenID {
				act = NewActAccept(item.Rule)
			} else {
				act = NewActReduce(item.Rule)
			}
		} else {
			act = NewActShift[T]()
		}

		h := NewHelper(item, act)

		prev, ok := buckets[symbol]
		if !ok {
			prev = []*Helper[T]{h}
		} else {
			prev = append(prev, h)
		}

		buckets[symbol] = prev
	}

	return buckets
}

// GetBucketsCopy gets a copy of the buckets of the rule table.
//
// Returns:
//   - map[T]*uts.Bucket[*Helper]: The copy of the buckets.
func (rt *RuleTable[T]) GetBucketsCopy() map[T][]*Helper[T] {
	buckets := make(map[T][]*Helper[T])

	for k, v := range rt.buckets {
		v_copy := make([]*Helper[T], len(v))
		copy(v_copy, v)

		buckets[k] = v_copy
	}

	return buckets
}
