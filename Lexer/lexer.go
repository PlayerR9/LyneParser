package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/StatusTree"
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Lexer is a lexer that uses a grammar to tokenize a string.
type Lexer struct {
	// productions are the production rules to use.
	productions []*gr.RegProduction

	// toSkip are the tokens to skip.
	toSkip []string
}

// NewLexer creates a new lexer.
//
// Parameters:
//   - grammar: The grammar to use.
//
// Returns:
//   - Lexer: The new lexer.
//
// Example:
//
//	lexer, err := NewLexer(grammar)
//	if err != nil {
//	    // Handle error.
//	}
//
//	branches, err := lexer.Lex(lexer, []byte("1 + 2"))
//	if err != nil {
//	    // Handle error.
//	}
//
// // Continue with parsing.
func NewLexer(grammar *Grammar) *Lexer {
	lex := new(Lexer)

	if grammar == nil {
		return lex
	}

	lex.productions = grammar.GetRegexProds()
	lex.toSkip = grammar.GetToSkip()

	return lex
}

// Lex is the main function of the lexer. This can be parallelized.
//
// Parameters:
//   - lexer: The lexer to use.
//   - source: The source to lex.
//
// Returns:
//   - Lexer: The active lexer. Nil if there are no tokens to lex or
//     grammar is invalid.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//   - *gr.ErrNoProductionRulesFound: No production rules are found in the grammar.
func (l *Lexer) Lex(input []byte) *LexerIterator {
	if len(input) == 0 || len(l.productions) == 0 {
		return nil
	}

	prodCopy := make([]*gr.RegProduction, len(l.productions))
	copy(prodCopy, l.productions)
	toSkip := make([]string, len(l.toSkip))
	copy(toSkip, l.toSkip)

	stream := cds.NewStream(input)

	si := newSourceIterator(stream, prodCopy)

	li := newLexerIterator(toSkip, si)

	return li
}

type SourceIterator struct {
	// source is the source to lex.
	source *cds.Stream[byte]

	// tree is the tree to use.
	tree *ud.History[*tr.Tree[EvalStatus, *gr.LeafToken]]

	// productions are the production rules to use.
	productions []*gr.RegProduction

	// isFirst is a flag that indicates if the lexer is the first.
	isFirst bool

	// canContinue is a flag that indicates if the lexer can continue.
	canContinue bool
}

// Size implements the Iterater interface.
func (si *SourceIterator) Size() (count int) {
	f := func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		count = data.Size()
	}

	si.tree.ReadData(f)

	return
}

// Consume implements the Iterater interface.
func (si *SourceIterator) Consume() ([]*tr.TreeNode[EvalStatus, *gr.LeafToken], error) {
	if !si.canContinue {
		return nil, ui.NewErrExhaustedIter()
	}

	var leaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	err := si.lexOne()
	if err != nil {
		si.canContinue = false

		f := func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
			leaves = data.GetLeaves()
		}

		si.tree.ReadData(f)
	} else {
		leaves, si.canContinue = si.getCompletedBranch()
	}

	return leaves, nil
}

// Restart implements the Iterater interface.
func (si *SourceIterator) Restart() {
	si.isFirst = true
	si.canContinue = true

	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)
	si.tree = tree
}

func newSourceIterator(source *cds.Stream[byte], productions []*gr.RegProduction) *SourceIterator {
	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)

	si := &SourceIterator{
		source:      source,
		productions: productions,
		tree:        tree,
		isFirst:     true,
		canContinue: true,
	}

	return si
}

func (si *SourceIterator) lexOne() error {
	if si.isFirst {
		matches, err := matchFrom(si.source, 0, si.productions)
		if err != nil {
			return ers.NewErrAt(0, "position", err)
		}

		children := generateEvalTrees(matches)

		cmd := tr.NewSetChildrenCmd(children)
		si.tree.ExecuteCommand(cmd)

		si.tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
			root := data.Root()
			root.ChangeStatus(EvalComplete)
		})

		si.isFirst = false
	} else {
		p := filterLeaves(si.source, si.productions)

		cmd1 := tr.NewProcessLeavesCmd(p)
		err := si.tree.ExecuteCommand(cmd1)
		if err != nil {
			si.tree.UndoLastCommand()
			return err
		}

		si.tree.Accept()

		cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
		err = si.tree.ExecuteCommand(cmd2)
		if err != nil {
			si.tree.UndoLastCommand()
			return err
		}

		si.tree.Accept()
	}

	return nil
}

func (si *SourceIterator) getCompletedBranch() ([]*tr.TreeNode[EvalStatus, *gr.LeafToken], bool) {
	var leaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	f := func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		leaves = data.GetLeaves()
	}

	si.tree.ReadData(f)

	canContinue := false

	var completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	for _, leaf := range leaves {
		status := leaf.GetStatus()

		if status == EvalComplete {
			completedLeaves = append(completedLeaves, leaf)
		} else if status == EvalIncomplete {
			canContinue = true
		}
	}

	return completedLeaves, canContinue
}

func (si *SourceIterator) deleteBranch(leaf *tr.TreeNode[EvalStatus, *gr.LeafToken]) {
	cmd := tr.NewDeleteBranchContainingCmd(leaf)
	si.tree.ExecuteCommand(cmd)
	si.tree.Accept()
}

// LexerIterator is an iterator that uses a grammar to tokenize a string.
type LexerIterator struct {
	// toSkip are the tokens to skip.
	toSkip []string

	// completedLeaves are the leaves that have been completed.
	completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

	// sourceIter is the source iterator.
	sourceIter *SourceIterator
}

// Size implements the Iterater interface.
//
// Returns the current size of the lexing tree.
func (li *LexerIterator) Size() (count int) {
	count = li.sourceIter.Size()

	return
}

// Consume implements the Iterater interface.
func (li *LexerIterator) Consume() (*cds.Stream[*gr.LeafToken], error) {
	if len(li.completedLeaves) == 0 {
		leaves, err := li.sourceIter.Consume()
		if err != nil {
			return nil, err
		}

		li.completedLeaves = leaves
	}

	// Extract the entire branch
	leaf := li.completedLeaves[0]

	anch := leaf.GetAncestors()
	anch = append(anch, leaf)

	li.completedLeaves = li.completedLeaves[1:]

	li.sourceIter.deleteBranch(leaf)

	branch := convertBranchToTokenStream(anch, li.toSkip)

	return branch, nil
}

// Restart implements the Iterater interface.
func (li *LexerIterator) Restart() {
	li.completedLeaves = nil
	li.sourceIter.Restart()
}

func newLexerIterator(toSkip []string, sourceIter *SourceIterator) *LexerIterator {
	li := &LexerIterator{
		toSkip:     toSkip,
		sourceIter: sourceIter,
	}

	return li
}

// FIXME: Remove this function once MyGoLib is updated.
func extractData(tree *ud.History[*tr.Tree[EvalStatus, *gr.LeafToken]]) *tr.Tree[EvalStatus, *gr.LeafToken] {
	var tmp *tr.Tree[EvalStatus, *gr.LeafToken]

	tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		tmp = data
	})

	return tmp
}

// generateEvalTrees adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func generateEvalTrees(matches []*gr.MatchedResult[*gr.LeafToken]) []*tr.Tree[EvalStatus, *gr.LeafToken] {
	// Get the longest match.
	matches = selectBestMatches(matches)

	children := make([]*tr.Tree[EvalStatus, *gr.LeafToken], 0, len(matches))

	for _, match := range matches {
		currMatch := match.GetMatch()

		tree := tr.NewTree(EvalIncomplete, currMatch)
		children = append(children, tree)
	}

	return children
}

// filterLeaves processes the leaves in the tree evaluator.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func filterLeaves(source *cds.Stream[byte], productions []*gr.RegProduction) uc.EvalManyFunc[*tr.TreeNode[EvalStatus, *gr.LeafToken], uc.Pair[EvalStatus, *gr.LeafToken]] {
	filterFunc := func(leaf *tr.TreeNode[EvalStatus, *gr.LeafToken]) ([]uc.Pair[EvalStatus, *gr.LeafToken], error) {
		nextAt := leaf.Data.GetPos() + len(leaf.Data.Data)

		if nextAt >= source.Size() {
			leaf.ChangeStatus(EvalComplete)
			return nil, nil
		}

		matches, err := matchFrom(source, nextAt, productions)
		if err != nil {
			leaf.ChangeStatus(EvalError)
			return nil, nil
		}

		// Get the longest match.
		matches = selectBestMatches(matches)

		children := make([]uc.Pair[EvalStatus, *gr.LeafToken], 0, len(matches))

		for _, match := range matches {
			curr := match.GetMatch()
			p := uc.NewPair(EvalIncomplete, curr)

			children = append(children, p)
		}

		leaf.ChangeStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

/*

// executeLexing is the main function of the tree evaluator.
//
// Parameters:
//   - source: The source to executeLexing.
//   - root: The root of the tree evaluator.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrEmptyInput: The source is empty.
//   - *ers.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func executeLexing(source *cds.Stream[byte], productions []*gr.RegProduction) (*tr.Tree[EvalStatus, *gr.LeafToken], error) {
	rootNode := gr.NewRootToken()

	tree := tr.NewTreeWithHistory(EvalIncomplete, rootNode)

	matches, err := matchFrom(source, 0, productions)
	if err != nil {
		data := extractData(tree)

		return data, ers.NewErrAt(0, "position", err)
	}

	children := generateEvalTrees(matches)

	cmd := tr.NewSetChildrenCmd(children)
	tree.ExecuteCommand(cmd)

	tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
		root := data.Root()
		root.ChangeStatus(EvalComplete)
	})

	shouldContinue := true

	for shouldContinue {
		tree.Accept() // Accept the changes made by the last command.

		p := filterLeaves(source, productions)

		cmd1 := tr.NewProcessLeavesCmd(p)
		err := tree.ExecuteCommand(cmd1)
		if err != nil {
			tree.UndoLastCommand()

			data := extractData(tree)
			return data, err
		}

		tree.Accept()

		cmd2 := tr.NewPruneTreeCmd(filterErrorLeaves)
		err = tree.ExecuteCommand(cmd2)
		if err != nil {
			tree.UndoLastCommand()

			data := extractData(tree)
			return data, err
		}

		tree.Accept()

		var completedLeaves []*tr.TreeNode[EvalStatus, *gr.LeafToken]

		tree.ReadData(func(data *tr.Tree[EvalStatus, *gr.LeafToken]) {
			completedLeaves = data.GetLeaves()
		})

		shouldContinue = false
		for _, leaf := range completedLeaves {
			status := leaf.GetStatus()

			if status == EvalIncomplete {
				shouldContinue = true
				break
			}
		}
	}

	cmd3 := tr.NewPruneTreeCmd(filterIncompleteLeaves)
	err = tree.ExecuteCommand(cmd3)
	if err != nil {
		tree.UndoLastCommand()

		data := extractData(tree)
		return data, NewErrAllMatchesFailed()
	}

	tree.Accept()

	data := extractData(tree)

	return data, nil
}
*/
