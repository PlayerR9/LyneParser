package Interfaces

// Equaler is an interface that defines a method to compare two objects
// of the same type.
type Equaler[T any] interface {
	// Equals returns true if the object is equal to the other object.
	//
	// Parameters:
	//
	// 	- other: The other object to compare to.
	//
	// Returns:
	//
	// 	- bool: True if the object is equal to the other object.
	Equals(other T) bool
}
