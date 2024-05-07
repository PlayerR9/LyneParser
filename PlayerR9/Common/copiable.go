package Interfaces

// Copier is an interface that provides a method to create a deep copy of an object.
type Copier interface {
	// Copy creates a shallow copy of the object.
	//
	// Returns:
	//
	//   - Copier: A shallow copy of the object.
	Copy() Copier
}

// CopyOf creates a copy of the element.
//
// Parameters:
//   - elem: The element to copy.
//
// Returns:
//   - any: A copy of the element.
func CopyOf(elem any) any {
	switch elem := elem.(type) {
	case Copier:
		return elem.Copy()
	default:
		return elem
	}
}
