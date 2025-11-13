// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor

import "fmt"

// Paginate generations all cursors to navigate from a cursor.
func Paginate[T Pointer](c *Cursor[T], secret []byte) (*Pagination, error) {
	if len(secret) == 0 {
		return &Pagination{
			First: First(c).String(),
			Prev:  Prev(c).String(),
			Last:  Last(c).String(),
			Next:  Next(c).String(),
		}, nil
	}
	var (
		p   Pagination
		err error
	)
	p.First, err = encryptString(First(c), secret)
	if err != nil {
		return nil, fmt.Errorf("first: %w", err)
	}
	p.Prev, err = encryptString(Prev(c), secret)
	if err != nil {
		return nil, fmt.Errorf("prev: %w", err)
	}
	p.Next, err = encryptString(Next(c), secret)
	if err != nil {
		return nil, fmt.Errorf("next: %w", err)
	}
	p.Last, err = encryptString(Last(c), secret)
	if err != nil {
		return nil, fmt.Errorf("last: %w", err)
	}
	return &p, nil
}

// Pagination contains all cursors to navigate from a cursor.
type Pagination struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
}

func encryptString[T Pointer](c *Cursor[T], secret []byte) (string, error) {
	b, err := Encrypt(c, secret)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
