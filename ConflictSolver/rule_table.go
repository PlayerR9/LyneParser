package ConflictSolver

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
)

// RuleTable represents a table of items.
type RuleTable struct {
	// items is the items of the rule table.
	items []*Item

	// buckets is the buckets of the rule table.
	buckets map[string]*uts.Bucket[*Helper]
}

// NewRuleTable is a constructor of RuleTable.
//
// Parameters:
//   - symbols: The symbols to use.
//   - rules: The rules to use.
//
// Returns:
//   - *RuleTable: The new rule table.
func NewRuleTable(symbols []string, rules []*gr.Production) *RuleTable {
	rt := &RuleTable{
		items: make([]*Item, 0),
	}

	for _, s := range symbols {
		for _, r := range rules {
			indices := r.IndicesOfRhs(s)

			for _, i := range indices {
				item, err := NewItem(r, i, len(rt.items))
				if err != nil {
					panic(err)
				}

				rt.items = append(rt.items, item)
			}
		}
	}

	rt.buckets = rt.getItemBuckets()

	return rt
}

// getItemBuckets gets the item buckets of the rule table.
//
// Returns:
//   - map[string]*uts.Bucket[*Helper]: The item buckets.
func (rt *RuleTable) getItemBuckets() map[string]*uts.Bucket[*Helper] {
	buckets := make(map[string]*uts.Bucket[*Helper])

	for _, item := range rt.items {
		symbol, err := item.Rule.GetRhsAt(item.Pos)
		if err != nil {
			panic(err)
		}

		lastIndex := item.Rule.Size() - 1

		var act HelperElem

		if item.Pos == lastIndex {
			act = NewActReduce(item.Rule, symbol == gr.EOFTokenID)
		} else {
			act = NewActShift()
		}

		h := NewHelper(item, act)

		prev, ok := buckets[symbol]
		if !ok {
			buckets[symbol] = uts.NewBucket([]*Helper{h})
		} else {
			prev.Add(h)
		}
	}

	return buckets
}

// GetBucketsCopy gets a copy of the buckets of the rule table.
//
// Returns:
//   - map[string]*uts.Bucket[*Helper]: The copy of the buckets.
func (rt *RuleTable) GetBucketsCopy() map[string]*uts.Bucket[*Helper] {
	buckets := make(map[string]*uts.Bucket[*Helper])

	for k, v := range rt.buckets {
		buckets[k] = v.Copy().(*uts.Bucket[*Helper])
	}

	return buckets
}
