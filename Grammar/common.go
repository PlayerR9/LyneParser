package Grammar

const (
	// LeftToRight is the direction of a production from left to right.
	LeftToRight string = "->"

	// ArrowLen is the length of the arrow.
	ArrowLen int = 2

	// StartSymbolID is the identifier of the start symbol in the grammar.
	StartSymbolID string = "source"

	// EpsilonSymbolID is the identifier of the epsilon symbol in the grammar.
	EpsilonSymbolID string = "ε"
)

/*
LHS : [a-z][a-zA-Z0-9_]* ;
RHS : [a-zA-Z0-9_]+ ;
ARROW : '->' ;

STAR : '*' ;
PLUS : '+' ;
QUESTION : '?' ;

PIPE : '|' ;

OP_PAREN : '(' ;
CL_PAREN : ')' ;

WS : [ \t\n\r]+ -> skip ;
*/

/*
rule :
	LHS ARROW rhsCls
	;

rhsCls :
	rhs+
	;

rhs :
	(WORD | OP_PAREN rhsCls CL_PAREN) (STAR | PLUS | QUESTION)?
	| WORD PIPE rhsCls
	;


func hasStar(rhs string) (string, bool) {
	res := strings.TrimRight(rhs, "*")

	return res, len(res) != len(rhs)
}

func hasPlus(rhs string) (string, bool) {
	res := strings.TrimRight(rhs, "+")

	return res, len(res) != len(rhs)
}

func hasQuestion(rhs string) (string, bool) {
	res := strings.TrimRight(rhs, "?")

	return res, len(res) != len(rhs)
}

// ParseProductionRule parses a production rule.
//
// Parameters:
//   - rule: The rule to parse.
//
// Returns:
//   - *Production: The production.
//   - error: An error if there was a problem parsing the rule.
//
// Format:
//   - A -> B C D...
//   - A -> B C ... D?
func ParseProductionRule(rule string) (*Production, error) {
	fields := strings.Split(rule, " -> ")

	if len(fields) == 1 {
		return nil, errors.New("missing either LHS or RHS")
	} else if len(fields) > 2 {
		return nil, errors.New("too many ->")
	}

	var lhs, rhs string

	lhs = strings.TrimSpace(fields[0])
	rhs = strings.TrimSpace(fields[1])

	if lhs == "" {
		return nil, errors.New("missing LHS")
	} else if rhs == "" {
		return nil, errors.New("missing RHS")
	}

	rhss := strings.Fields(rhs)
	rhss = us.RemoveEmpty(rhss)

	if len(rhss) == 0 {
		return nil, errors.New("missing RHS")
	}

	return NewProduction(lhs, rhss), nil
}
*/
