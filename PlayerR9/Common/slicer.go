package Interfaces

// Slicer is an interface that defines a method to convert a data structure to a slice.
// It is implemented by data structures that can be converted to a slice of elements of
// the same type.
type Slicer[T any] interface {
	// The Slice method returns a slice containing all the elements in the data structure.
	Slice() []T
}
