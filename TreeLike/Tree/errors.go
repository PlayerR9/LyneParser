package Tree

// ErrMissingRoot is an error that is returned when the root of a tree is missing.
type ErrMissingRoot struct{}

// Error returns the error message: "missing root".
//
// Returns:
//   - string: The error message.
func (e *ErrMissingRoot) Error() string {
	return "missing root"
}

// NewErrMissingRoot creates a new ErrMissingRoot.
//
// Returns:
//   - *ErrMissingRoot: The newly created error.
func NewErrMissingRoot() *ErrMissingRoot {
	return &ErrMissingRoot{}
}

// ErrNodeNotPartOfTree is an error that is returned when a node is not part of a tree.
type ErrNodeNotPartOfTree struct{}

// Error returns the error message: "node is not part of the tree".
//
// Returns:
//   - string: The error message.
func (e *ErrNodeNotPartOfTree) Error() string {
	return "node is not part of the tree"
}

// NewErrNodeNotPartOfTree creates a new ErrNodeNotPartOfTree.
//
// Returns:
//   - *ErrNodeNotPartOfTree: The newly created error.
func NewErrNodeNotPartOfTree() *ErrNodeNotPartOfTree {
	return &ErrNodeNotPartOfTree{}
}
