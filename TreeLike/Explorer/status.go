package TreeExplorer

// EvalStatus represents the status of an evaluation.
type EvalStatus int

const (
	// EvalComplete represents a completed evaluation.
	EvalComplete EvalStatus = iota

	// EvalIncomplete represents an incomplete evaluation.
	EvalIncomplete

	// EvalError represents an evaluation that has an error.
	EvalError
)

// String implements fmt.Stringer interface.
func (s EvalStatus) String() string {
	return [...]string{
		"complete",
		"incomplete",
		"error",
	}[s]
}
