// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor

// Pointer must be implemented by any cursor point.
type Pointer interface {
	// Args returns the arguments to use in a statement.
	Args() []any
	// IsZero returns true if the pointer is a zero value.
	IsZero() bool
}

// Int64 mangers int64 pointer.
type Int64 int64

// Args implements the Pointer interface.
func (n Int64) Args() []any {
	if n.IsZero() {
		return nil
	}
	return []any{n}
}

// IsZero implements the Pointer interface.
func (n Int64) IsZero() bool {
	return n == 0
}

// List allows manipulation of a list of pointer data as pointer.
type List []Pointer

// Args implements the Pointer interface.
func (l List) Args() []any {
	n := len(l)
	if n == 0 {
		return nil
	}
	var (
		a = make([]any, n)
		c []any
	)
	for k := range l {
		if c = l[k].Args(); len(c) == 1 {
			a[k] = c[0]
		}
	}
	return a
}

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

// Args implements the Pointer interface.
func (s String) Args() []any {
	if s.IsZero() {
		return nil
	}
	return []any{s}
}

// IsZero implements the Pointer interface.
func (s String) IsZero() bool {
	return s == ""
}
