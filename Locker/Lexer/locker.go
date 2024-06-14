package Lexer

import (
	"strings"

	gw "github.com/PlayerR9/LyneParser/Util/CodeWriter"
)

var (
	HelperFunctionsFile *gw.GoFile
	LexerMatcherFile    *gw.GoFile
	ErrorsFile          *gw.GoFile
)

type Info struct {
	NumberOfToSkip int
}

func MakeGrammarFile(info *Info, productions string, toSkip string) *gw.GoFile {
	lines := strings.Split(productions, "\n")
	content1 := gw.WriteBacktickString(lines)

	lines = strings.Split(toSkip, "\n")
	content2 := gw.WriteBacktickString(lines)

	var content string

	switch info.NumberOfToSkip {
	case 0:
		content = `
		package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var (
	Productions []*gr.RegProduction
)

func init() {
	grammar, err := gr.NewLexerGrammar(
		` + content1 + `,
		"",
	)
	if err != nil {
		panic(err)
	}

	Productions = grammar.GetRegProductions()
	if len(Productions) == 0 {
		panic("No productions found")
	}
}
	`
	case 1:
		content = `
		package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var (
	Productions []*gr.RegProduction
	ToSkip      string
)

func init() {
	grammar, err := gr.NewLexerGrammar(
		` + content1 + `,
		` + content2 + `,
	)
	if err != nil {
		panic(err)
	}

	Productions = grammar.GetRegProductions()
	if len(Productions) == 0 {
		panic("No productions found")
	}

	ToSkip = grammar.GetToSkip()[0]
}
	`
	default:
		content = `
		package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
)

var (
	Productions []*gr.RegProduction
	ToSkip      []string
)

func init() {
	grammar, err := gr.NewLexerGrammar(
		` + content1 + `,
		` + content2 + `,
	)
	if err != nil {
		panic(err)
	}

	Productions = grammar.GetRegProductions()
	if len(Productions) == 0 {
		panic("No productions found")
	}

	ToSkip = grammar.GetToSkip()
}
	`
	}

	GrammarFile := &gw.GoFile{
		FileName: "grammar.go",
		Content:  content,
	}

	return GrammarFile
}

func init() {
	ErrorsFile = &gw.GoFile{
		FileName: "errors.go",
		Content: `
		
// ErrNoMatches is an error that is returned when there are no
// matches at a position.
type ErrNoMatches struct{}

// Error returns the error message: "no matches".
//
// Returns:
//   - string: The error message.
func (e *ErrNoMatches) Error() string {
	return "no matches"
}

// NewErrNoMatches creates a new error of type *ErrNoMatches.
//
// Returns:
//   - *ErrNoMatches: The new error.
func NewErrNoMatches() *ErrNoMatches {
	return &ErrNoMatches{}
}

// ErrAllMatchesFailed is an error that is returned when all matches
// fail.
type ErrAllMatchesFailed struct{}

// Error returns the error message: "all matches failed".
//
// Returns:
//   - string: The error message.
func (e *ErrAllMatchesFailed) Error() string {
	return "all matches failed"
}

// NewErrAllMatchesFailed creates a new error of type *ErrAllMatchesFailed.
//
// Returns:
//   - *ErrAllMatchesFailed: The new error.
func NewErrAllMatchesFailed() *ErrAllMatchesFailed {
	return &ErrAllMatchesFailed{}
}

// ErrInvalidElement is an error that is returned when an invalid element
// is found.
type ErrInvalidElement struct{}

// Error returns the error message: "invalid element".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidElement) Error() string {
	return "invalid element"
}

// NewErrInvalidElement creates a new error of type *ErrInvalidElement.
//
// Returns:
//   - *ErrInvalidElement: The new error.
func NewErrInvalidElement() *ErrInvalidElement {
	return &ErrInvalidElement{}
}

		`,
	}

	HelperFunctionsFile = &gw.GoFile{
		FileName: "helper_functions.go",
		Content: `
		package Lexer

import (
	com "github.com/PlayerR9/LyneParser/Common"
	gr "github.com/PlayerR9/LyneParser/Grammar"
	teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
)
	
// SetEOFToken sets the end-of-file token in the token stream.
//
// If the end-of-file token is already present, it will not be added again.
func setEOFToken(tokens []*gr.LeafToken) []*gr.LeafToken {
	if len(tokens) != 0 && tokens[len(tokens)-1].ID == gr.EOFTokenID {
		// EOF token is already present
		return tokens
	}

	return append(tokens, gr.NewEOFToken())
}

// SetLookahead sets the lookahead token for all the tokens in the stream.
func setLookahead(tokens []*gr.LeafToken) {
	for i, token := range tokens[:len(tokens)-1] {
		token.SetLookahead(tokens[i+1])
	}
}

// convertBranchToTokenStream converts a branch to a token stream.
//
// Parameters:
//   - branch: The branch to convert.
//
// Returns:
//   - *gr.TokenStream: The token stream.
func convertBranchToTokenStream(branch []*teval.CurrentEval[*gr.LeafToken]) *com.TokenStream {
	ts := make([]*gr.LeafToken, 0, len(branch))

	for _, leaf := range branch {
		ts = append(ts, leaf.GetElem())
	}

	ts = setEOFToken(ts)

	setLookahead(ts)

	return com.NewTokenStream(ts)
}
`,
	}

	LexerMatcherFile = &gw.GoFile{
		FileName: "lexer_matcher.go",
		Content: `
		package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// LexerMatcher is a struct that represents a lexer matcher.
type LexerMatcher struct {
	// source is the source to match.
	source *cds.Stream[byte]
}

// IsDone is a function that checks if the matcher is done.
//
// Parameters:
//   - from: The starting position of the match.
//
// Returns:
//   - bool: True if the matcher is done, false otherwise.
func (lm *LexerMatcher) IsDone(from int) bool {
	return lm.source.IsDone(from, 1)
}

// Match is a function that matches the element.
//
// Parameters:
//   - from: The starting position of the match.
//
// Returns:
//   - []Matcher: The list of matchers.
//   - error: An error if the matchers cannot be created.
func (lm *LexerMatcher) Match(from int) ([]*gr.MatchedResult[*gr.LeafToken], error) {
	matched, err := MatchFrom(lm.source, from, Productions)
	if err != nil {
		return nil, err
	}

	return matched, nil
}

// SelectBestMatches selects the best matches from the list of matches.
// Usually, the best matches' euristic is the longest match.
//
// Parameters:
//   - matches: The list of matches.
//
// Returns:
//   - []Matcher: The best matches.
func (lm *LexerMatcher) SelectBestMatches(matches []*gr.MatchedResult[*gr.LeafToken]) []*gr.MatchedResult[*gr.LeafToken] {
	weights := us.ApplyWeightFunc(matches, MatchWeightFunc)
	pairs := us.FilterByPositiveWeight(weights)

	return us.ExtractResults(pairs)
}

// GetNext is a function that returns the next position of an element.
//
// Parameters:
//   - elem: The element to get the next position of.
//
// Returns:
//   - int: The next position of the element.
func (lm *LexerMatcher) GetNext(elem *gr.LeafToken) int {
	return elem.GetPos() + len(elem.Data)
}
	`,
	}
}

func MakeCommonFile(info *Info) *gw.GoFile {
	var content string

	switch info.NumberOfToSkip {
	case 0:
		content = `package Lexer

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"
	cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
	teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// MatchWeightFunc is a weight function that returns the length of the match.
	//
	// Parameters:
	//   - match: The match to weigh.
	//
	// Returns:
	//   - float64: The weight of the match.
	//   - bool: True if the weight is valid, false otherwise.
	MatchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]

	// FilterEmptyBranch is a filter that filters out empty branches.
	//
	// Parameters:
	//   - branch: The branch to filter.
	//
	// Returns:
	//   - bool: True if the branch is not empty, false otherwise.
	FilterEmptyBranch us.PredicateFilter[[]*teval.CurrentEval[*gr.LeafToken]]

	// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
	//
	// Parameters:
	//   - branches: The branches to remove tokens from.
	//
	// Returns:
	//   - []gr.TokenStream: The branches with the tokens removed.
	RemoveToSkipTokens teval.FilterBranchesFunc[*gr.LeafToken]
)

func init() {
	MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
		return float64(len(elem.Matched.Data)), true
	}

	FilterEmptyBranch = func(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
		return len(branch) != 0
	}

	RemoveToSkipTokens = func(branches [][]*teval.CurrentEval[*gr.LeafToken]) ([][]*teval.CurrentEval[*gr.LeafToken], error) {
		var newBranches [][]*teval.CurrentEval[*gr.LeafToken]

		for _, branch := range branches {
			if len(branch) != 0 {
				newBranches = append(newBranches, branch[1:])
			}
		}

		return newBranches, nil
	}
}

// MatchFrom matches the source stream from a given index with a list of production rules.
//
// Parameters:
//   - s: The source stream to match.
//   - from: The index to start matching from.
//   - ps: The production rules to match.
//
// Returns:
//   - matches: A slice of MatchedResult that match the input token.
//   - reason: An error if no matches are found.
//
// Errors:
//   - *ue.ErrInvalidParameter: The from index is out of bounds.
//   - *ErrNoMatches: No matches are found.
func MatchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
	size := s.Size()

	if from < 0 || from >= size {
		reason = ue.NewErrInvalidParameter(
			"from",
			ue.NewErrOutOfBounds(from, 0, size),
		)

		return
	}

	subSet, err := s.Get(from, size)
	if err != nil {
		panic(err)
	}

	for i, p := range ps {
		matched := p.Match(from, subSet)
		if matched != nil {
			matches = append(matches, gr.NewMatchResult(matched, i))
		}
	}

	if len(matches) == 0 {
		reason = NewErrNoMatches()
	}

	return
}

// Lex is the main function of the lexer.
//
// Parameters:
//   - source: The source to lex.
//
// Returns:
//   - []*gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//
// Example:
//
//	err := Lex([]byte("1 + 2")))
//	if err != nil {
//	    // Handle error.
//	}
//
//	tokenBranches, err := lexer.GetTokens()
//	if err != nil {
//	    // Handle error.
//	} else if len(tokenBranches) == 0 {
//	    // No tokens found.
//	}
//
//	tokenBranches = lexer.RemoveToSkipTokens(tokenBranches) // prepare for parsing
//
// // DEBUG: Print tokens.
//
//	for _, branch := range tokenBranches {
//	    for _, token := range branch {
//	        fmt.Println(token)
//	    }
//	}
//
// // Continue with parsing.
func Lex(source []byte) ([]*com.TokenStream, error) {
	lexer := teval.NewTreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher](
		RemoveToSkipTokens,
	)

	matcher := &LexerMatcher{
		source: cds.NewStream(source),
	}

	err := lexer.Evaluate(matcher, gr.NewRootToken())

	branches, err2 := lexer.GetBranches()

	var result []*com.TokenStream
	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)

		result = append(result, conv)
	}

	if err != nil {
		return result, err
	}

	return result, err2
}

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc[O any] func(branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], error)

// MatchResult is an interface that represents a match result.
type MatchResulter[O any] interface {
	// GetMatch returns the match.
	//
	// Returns:
	//   - O: The match.
	GetMatch() O
}

// Matcher is an interface that represents a matcher.
type Matcher[R MatchResulter[O], O any] interface {
	// IsDone is a function that checks if the matcher is done.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - bool: True if the matcher is done, false otherwise.
	IsDone(from int) bool

	// Match is a function that matches the element.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - []R: The list of matched results.
	//   - error: An error if the matchers cannot be created.
	Match(from int) ([]R, error)

	// SelectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []T: The best matches.
	SelectBestMatches(matches []R) []R

	// GetNext is a function that returns the next position of an element.
	//
	// Parameters:
	//   - elem: The element to get the next position of.
	//
	// Returns:
	//   - int: The next position of the element.
	GetNext(elem O) int
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves[O any](h *CurrentEval[O]) bool {
	return h == nil || h.Status == EvalError
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches[O any](branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], int) {
	branches, ok := us.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	} else if len(branches) == 0 {
		return nil, -1
	}

	// Return the longest branch.
	weights := us.ApplyWeightFunc(branches, HelperWeightFunc)
	weights = us.FilterByPositiveWeight(weights)

	elems := weights[0].GetData().First

	return [][]*CurrentEval[O]{elems}, len(elems)
}

	`
	case 1:
		content = `package Lexer
	
	import (
		gr "github.com/PlayerR9/LyneParser/Grammar"
		cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
		teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
		ue "github.com/PlayerR9/MyGoLib/Units/errors"
		us "github.com/PlayerR9/MyGoLib/Units/slice"
	)
	
	var (
		// MatchWeightFunc is a weight function that returns the length of the match.
		//
		// Parameters:
		//   - match: The match to weigh.
		//
		// Returns:
		//   - float64: The weight of the match.
		//   - bool: True if the weight is valid, false otherwise.
		MatchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]
	
		// FilterEmptyBranch is a filter that filters out empty branches.
		//
		// Parameters:
		//   - branch: The branch to filter.
		//
		// Returns:
		//   - bool: True if the branch is not empty, false otherwise.
		FilterEmptyBranch us.PredicateFilter[[]*teval.CurrentEval[*gr.LeafToken]]
	
		// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
		//
		// Parameters:
		//   - branches: The branches to remove tokens from.
		//
		// Returns:
		//   - []gr.TokenStream: The branches with the tokens removed.
		RemoveToSkipTokens teval.FilterBranchesFunc[*gr.LeafToken]
	)
	
	func init() {
		MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
			return float64(len(elem.Matched.Data)), true
		}
	
		FilterEmptyBranch = func(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
			return len(branch) != 0
		}
	
		RemoveToSkipTokens = func(branches [][]*teval.CurrentEval[*gr.LeafToken]) ([][]*teval.CurrentEval[*gr.LeafToken], error) {
			var newBranches [][]*teval.CurrentEval[*gr.LeafToken]
			var reason error

			for _, branch := range branches {
				if len(branch) != 0 {
					newBranches = append(newBranches, branch[1:])
				}
			}

			newBranches = us.SliceFilter(newBranches, FilterEmptyBranch)
			if len(newBranches) == 0 {
				reason = teval.NewErrAllMatchesFailed()

				return newBranches, reason
			}

			filterTokenDifferentID := func(h *teval.CurrentEval[*gr.LeafToken]) bool {
				return h.GetElem().ID != ToSkip
			}

			for i := 0; i < len(newBranches); i++ {
				newBranches[i] = us.SliceFilter(newBranches[i], filterTokenDifferentID)
			}

			return newBranches, reason
		}
	}
	}
	
	// MatchFrom matches the source stream from a given index with a list of production rules.
	//
	// Parameters:
	//   - s: The source stream to match.
	//   - from: The index to start matching from.
	//   - ps: The production rules to match.
	//
	// Returns:
	//   - matches: A slice of MatchedResult that match the input token.
	//   - reason: An error if no matches are found.
	//
	// Errors:
	//   - *ue.ErrInvalidParameter: The from index is out of bounds.
	//   - *ErrNoMatches: No matches are found.
	func MatchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
		size := s.Size()
	
		if from < 0 || from >= size {
			reason = ue.NewErrInvalidParameter(
				"from",
				ue.NewErrOutOfBounds(from, 0, size),
			)
	
			return
		}
	
		subSet, err := s.Get(from, size)
		if err != nil {
			panic(err)
		}
	
		for i, p := range ps {
			matched := p.Match(from, subSet)
			if matched != nil {
				matches = append(matches, gr.NewMatchResult(matched, i))
			}
		}
	
		if len(matches) == 0 {
			reason = NewErrNoMatches()
		}
	
		return
	}

// Lex is the main function of the lexer.
//
// Parameters:
//   - source: The source to lex.
//
// Returns:
//   - []*gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//
// Example:
//
//	err := Lex([]byte("1 + 2")))
//	if err != nil {
//	    // Handle error.
//	}
//
//	tokenBranches, err := lexer.GetTokens()
//	if err != nil {
//	    // Handle error.
//	} else if len(tokenBranches) == 0 {
//	    // No tokens found.
//	}
//
//	tokenBranches = lexer.RemoveToSkipTokens(tokenBranches) // prepare for parsing
//
// // DEBUG: Print tokens.
//
//	for _, branch := range tokenBranches {
//	    for _, token := range branch {
//	        fmt.Println(token)
//	    }
//	}
//
// // Continue with parsing.
func Lex(source []byte) ([]*com.TokenStream, error) {
	lexer := teval.NewTreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher](
		RemoveToSkipTokens,
	)

	matcher := &LexerMatcher{
		source: cds.NewStream(source),
	}

	err := lexer.Evaluate(matcher, gr.NewRootToken())

	branches, err2 := lexer.GetBranches()

	var result []*com.TokenStream
	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)

		result = append(result, conv)
	}

	if err != nil {
		return result, err
	}

	return result, err2
}

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc[O any] func(branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], error)

// MatchResult is an interface that represents a match result.
type MatchResulter[O any] interface {
	// GetMatch returns the match.
	//
	// Returns:
	//   - O: The match.
	GetMatch() O
}

// Matcher is an interface that represents a matcher.
type Matcher[R MatchResulter[O], O any] interface {
	// IsDone is a function that checks if the matcher is done.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - bool: True if the matcher is done, false otherwise.
	IsDone(from int) bool

	// Match is a function that matches the element.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - []R: The list of matched results.
	//   - error: An error if the matchers cannot be created.
	Match(from int) ([]R, error)

	// SelectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []T: The best matches.
	SelectBestMatches(matches []R) []R

	// GetNext is a function that returns the next position of an element.
	//
	// Parameters:
	//   - elem: The element to get the next position of.
	//
	// Returns:
	//   - int: The next position of the element.
	GetNext(elem O) int
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves[O any](h *CurrentEval[O]) bool {
	return h == nil || h.Status == EvalError
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches[O any](branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], int) {
	branches, ok := us.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	} else if len(branches) == 0 {
		return nil, -1
	}

	// Return the longest branch.
	weights := us.ApplyWeightFunc(branches, HelperWeightFunc)
	weights = us.FilterByPositiveWeight(weights)

	elems := weights[0].GetData().First

	return [][]*CurrentEval[O]{elems}, len(elems)
}


		`
	default:
		content = `package Lexer
	
	import (
		gr "github.com/PlayerR9/LyneParser/Grammar"
		cds "github.com/PlayerR9/MyGoLib/CustomData/Stream"
		teval "github.com/PlayerR9/MyGoLib/TreeLike/Explorer"
		ue "github.com/PlayerR9/MyGoLib/Units/errors"
		us "github.com/PlayerR9/MyGoLib/Units/slice"
	)
	
	var (
		// MatchWeightFunc is a weight function that returns the length of the match.
		//
		// Parameters:
		//   - match: The match to weigh.
		//
		// Returns:
		//   - float64: The weight of the match.
		//   - bool: True if the weight is valid, false otherwise.
		MatchWeightFunc us.WeightFunc[*gr.MatchedResult[*gr.LeafToken]]
	
		// FilterEmptyBranch is a filter that filters out empty branches.
		//
		// Parameters:
		//   - branch: The branch to filter.
		//
		// Returns:
		//   - bool: True if the branch is not empty, false otherwise.
		FilterEmptyBranch us.PredicateFilter[[]*teval.CurrentEval[*gr.LeafToken]]
	
		// RemoveToSkipTokens removes tokens that are marked as to skip in the grammar.
		//
		// Parameters:
		//   - branches: The branches to remove tokens from.
		//
		// Returns:
		//   - []gr.TokenStream: The branches with the tokens removed.
		RemoveToSkipTokens teval.FilterBranchesFunc[*gr.LeafToken]
	)
	
	func init() {
		MatchWeightFunc = func(elem *gr.MatchedResult[*gr.LeafToken]) (float64, bool) {
			return float64(len(elem.Matched.Data)), true
		}
	
		FilterEmptyBranch = func(branch []*teval.CurrentEval[*gr.LeafToken]) bool {
			return len(branch) != 0
		}
	
		RemoveToSkipTokens = func(branches [][]*teval.CurrentEval[*gr.LeafToken]) ([][]*teval.CurrentEval[*gr.LeafToken], error) {
			var newBranches [][]*teval.CurrentEval[*gr.LeafToken]
			var reason error

			for _, branch := range branches {
				if len(branch) != 0 {
					newBranches = append(newBranches, branch[1:])
				}
			}

			for _, toSkip := range ToSkip {
				newBranches = us.SliceFilter(newBranches, FilterEmptyBranch)
				if len(newBranches) == 0 {
					reason = teval.NewErrAllMatchesFailed()

					return newBranches, reason
				}

				filterTokenDifferentID := func(h *teval.CurrentEval[*gr.LeafToken]) bool {
					return h.GetElem().ID != toSkip
				}

				for i := 0; i < len(newBranches); i++ {
					newBranches[i] = us.SliceFilter(newBranches[i], filterTokenDifferentID)
				}
			}

			return newBranches, reason
		}
	}
	
	// MatchFrom matches the source stream from a given index with a list of production rules.
	//
	// Parameters:
	//   - s: The source stream to match.
	//   - from: The index to start matching from.
	//   - ps: The production rules to match.
	//
	// Returns:
	//   - matches: A slice of MatchedResult that match the input token.
	//   - reason: An error if no matches are found.
	//
	// Errors:
	//   - *ue.ErrInvalidParameter: The from index is out of bounds.
	//   - *ErrNoMatches: No matches are found.
	func MatchFrom(s *cds.Stream[byte], from int, ps []*gr.RegProduction) (matches []*gr.MatchedResult[*gr.LeafToken], reason error) {
		size := s.Size()
	
		if from < 0 || from >= size {
			reason = ue.NewErrInvalidParameter(
				"from",
				ue.NewErrOutOfBounds(from, 0, size),
			)
	
			return
		}
	
		subSet, err := s.Get(from, size)
		if err != nil {
			panic(err)
		}
	
		for i, p := range ps {
			matched := p.Match(from, subSet)
			if matched != nil {
				matches = append(matches, gr.NewMatchResult(matched, i))
			}
		}
	
		if len(matches) == 0 {
			reason = NewErrNoMatches()
		}
	
		return
	}

// Lex is the main function of the lexer.
//
// Parameters:
//   - source: The source to lex.
//
// Returns:
//   - []*gr.TokenStream: The tokens that have been lexed.
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrNoTokensToLex: There are no tokens to lex.
//   - *ErrNoMatches: No matches are found in the source.
//   - *ErrAllMatchesFailed: All matches failed.
//
// Example:
//
//	err := Lex([]byte("1 + 2")))
//	if err != nil {
//	    // Handle error.
//	}
//
//	tokenBranches, err := lexer.GetTokens()
//	if err != nil {
//	    // Handle error.
//	} else if len(tokenBranches) == 0 {
//	    // No tokens found.
//	}
//
//	tokenBranches = lexer.RemoveToSkipTokens(tokenBranches) // prepare for parsing
//
// // DEBUG: Print tokens.
//
//	for _, branch := range tokenBranches {
//	    for _, token := range branch {
//	        fmt.Println(token)
//	    }
//	}
//
// // Continue with parsing.
func Lex(source []byte) ([]*com.TokenStream, error) {
	lexer := teval.NewTreeEvaluator[*gr.MatchedResult[*gr.LeafToken], *LexerMatcher](
		RemoveToSkipTokens,
	)

	matcher := &LexerMatcher{
		source: cds.NewStream(source),
	}

	err := lexer.Evaluate(matcher, gr.NewRootToken())

	branches, err2 := lexer.GetBranches()

	var result []*com.TokenStream
	for _, branch := range branches {
		conv := convertBranchToTokenStream(branch)

		result = append(result, conv)
	}

	if err != nil {
		return result, err
	}

	return result, err2
}

// FilterBranchesFunc is a function that filters branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]*CurrentEval: The filtered branches.
//   - error: An error if the branches are invalid.
type FilterBranchesFunc[O any] func(branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], error)

// MatchResult is an interface that represents a match result.
type MatchResulter[O any] interface {
	// GetMatch returns the match.
	//
	// Returns:
	//   - O: The match.
	GetMatch() O
}

// Matcher is an interface that represents a matcher.
type Matcher[R MatchResulter[O], O any] interface {
	// IsDone is a function that checks if the matcher is done.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - bool: True if the matcher is done, false otherwise.
	IsDone(from int) bool

	// Match is a function that matches the element.
	//
	// Parameters:
	//   - from: The starting position of the match.
	//
	// Returns:
	//   - []R: The list of matched results.
	//   - error: An error if the matchers cannot be created.
	Match(from int) ([]R, error)

	// SelectBestMatches selects the best matches from the list of matches.
	// Usually, the best matches' euristic is the longest match.
	//
	// Parameters:
	//   - matches: The list of matches.
	//
	// Returns:
	//   - []T: The best matches.
	SelectBestMatches(matches []R) []R

	// GetNext is a function that returns the next position of an element.
	//
	// Parameters:
	//   - elem: The element to get the next position of.
	//
	// Returns:
	//   - int: The next position of the element.
	GetNext(elem O) int
}

// FilterErrorLeaves is a filter that filters out leaves that are in error.
//
// Parameters:
//   - leaf: The leaf to filter.
//
// Returns:
//   - bool: True if the leaf is in error, false otherwise.
func FilterErrorLeaves[O any](h *CurrentEval[O]) bool {
	return h == nil || h.Status == EvalError
}

// filterInvalidBranches filters out invalid branches.
//
// Parameters:
//   - branches: The branches to filter.
//
// Returns:
//   - [][]helperToken: The filtered branches.
//   - int: The index of the last invalid token. -1 if no invalid token is found.
func filterInvalidBranches[O any](branches [][]*CurrentEval[O]) ([][]*CurrentEval[O], int) {
	branches, ok := us.SFSeparateEarly(branches, FilterIncompleteTokens)
	if ok {
		return branches, -1
	} else if len(branches) == 0 {
		return nil, -1
	}

	// Return the longest branch.
	weights := us.ApplyWeightFunc(branches, HelperWeightFunc)
	weights = us.FilterByPositiveWeight(weights)

	elems := weights[0].GetData().First

	return [][]*CurrentEval[O]{elems}, len(elems)
}

		`
	}

	CommonFile := &gw.GoFile{
		FileName: "common.go",
		Content:  content,
	}

	return CommonFile
}
