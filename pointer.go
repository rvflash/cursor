package cursor

// Pointer must be implemented by any cursor point.
type Pointer interface {
	IsZero() bool
}

// Int64 mangers int64 pointer.
type Int64 int64

// IsZero implements the Pointer interface.
func (n Int64) IsZero() bool {
	return n == 0
}

// List allows manipulation of a list of pointer data as pointer.
type List []Pointer

// IsZero implements the Pointer interface.
func (l List) IsZero() bool {
	if len(l) == 0 {
		return true
	}
	for k := range l {
		if !l[k].IsZero() {
			return false
		}
	}
	return true
}

// String manages string pointer.
type String string

// IsZero implements the Pointer interface.
func (n String) IsZero() bool {
	return n == ""
}
