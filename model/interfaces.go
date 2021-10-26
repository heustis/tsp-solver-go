package model

// Deletable should be implemented by any struct that needs to be cleaned up (to prevent resource leaks).
type Deletable interface {
	Delete()
}

// Equal defines a boolean comparison method.
type Equal interface {
	// Equals should return true if this object is equal to the supplied object.
	Equals(other interface{}) bool
}

// Printable defines a method for printing the object as a string.
type Printable interface {
	// ToString converts the object into a human-readable string.
	ToString() string
}
