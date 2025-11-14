// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package cursor uses a reference point (cursor) to fetch the next set of results.
// This reference point is typically a unique identifier that define the sort order.
package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"
)

const notFound = -1

var (
	// now is useful for tests purpose.
	now = time.Now
	b64 = base64.RawURLEncoding
	sep = []byte(".")
)

// First returns the cursor of the first page.
func First[T Pointer](c *Cursor[T]) *Cursor[T] {
	if c == nil || c.Prev == nil || c.Offset == 0 {
		return nil
	}
	if t := *c.Prev; t.IsZero() {
		return nil
	}
	return &Cursor[T]{
		Prev:    new(T),
		Offset:  0,
		Limit:   c.Limit,
		Total:   c.Total,
		Filters: c.Filters,
	}
}

// Last returns the cursor of the last page.
func Last[T Pointer](c *Cursor[T]) *Cursor[T] {
	if c == nil || c.Next == nil {
		return nil
	}
	if t := *c.Next; t.IsZero() {
		return nil
	}
	var (
		offset int
		total  = c.TotalPages()
	)
	if total > notFound {
		offset = c.Limit * (total - 1)
	}
	if c.Offset == offset-c.Limit {
		return Next(c)
	}
	return &Cursor[T]{
		Next:    new(T),
		Offset:  offset,
		Limit:   c.Limit,
		Total:   c.Total,
		Filters: c.Filters,
	}
}

// New creates a new cursor based on this limit and total.
func New[T Pointer](limit, total int) *Cursor[T] {
	c := Cursor[T]{Limit: limit}
	c.Filters = make(url.Values)
	if total > 0 {
		c.Total = &total
	}
	return &c
}

// Next returns the cursor of the next page.
func Next[T Pointer](c *Cursor[T]) *Cursor[T] {
	if c == nil || c.Next == nil {
		return nil
	}
	if t := *c.Next; t.IsZero() {
		return nil
	}
	return &Cursor[T]{
		Next:    c.Next,
		Offset:  c.Offset + c.Limit,
		Limit:   c.Limit,
		Total:   c.Total,
		Filters: c.Filters,
	}
}

// Prev returns the cursor of the previous page.
func Prev[T Pointer](c *Cursor[T]) *Cursor[T] {
	if c == nil || c.Prev == nil || c.Offset == 0 {
		return nil
	}
	if t := *c.Prev; t.IsZero() {
		return nil
	}
	if c.Offset == c.Limit {
		return First(c)
	}
	return &Cursor[T]{
		Prev:    c.Prev,
		Offset:  c.Offset - c.Limit,
		Limit:   c.Limit,
		Total:   c.Total,
		Filters: c.Filters,
	}
}

// Cursor contains elements required to paginate based on a cursor, a data pointed the start of the data to list.
type Cursor[T Pointer] struct {
	Prev     *T    `json:"prev,omitempty"`
	Next     *T    `json:"next,omitempty"`
	IssuedAt int64 `json:"issued_at,omitempty"` // epoch seconds
	Offset   int
	Limit    int        `json:"limit"`
	Total    *int       `json:"total,omitempty"`
	Filters  url.Values `json:"filters,omitempty"`

	cnt int
}

// Add notifies a new entry to the managed list of result.
func (c *Cursor[T]) Add(d T) {
	if c.Limit == 0 {
		return
	}
	switch c.cnt {
	case 0:
		c.Prev = &d
	case c.Limit:
		c.Next = &d
	}
	c.cnt++
}

// CurrentPage returns the current page number.
func (c *Cursor[T]) CurrentPage() int {
	if c == nil || c.Offset == 0 || c.Limit == 0 {
		return 1
	}
	return 1 + (c.Offset / c.Limit)
}

// Decode decodes a plain cursor.
func (c *Cursor[T]) Decode(text []byte) error {
	src, err := b64Decode(text)
	if err != nil {
		return fmt.Errorf("decoding: %w", err)
	}
	c2 := &Cursor[T]{}
	err = json.Unmarshal(src[:], &c2)
	if err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	*c = *c2
	return nil
}

// Encode encodes the cursor as plain data.
func (c *Cursor[T]) Encode() ([]byte, error) {
	if c == nil || (c.Prev == nil && c.Next == nil) {
		return nil, nil
	}
	c.IssuedAt = now().Unix()

	src, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}
	return b64Encode(src), nil
}

// IsExpired returns true if the issued timestamp exceeds the max age allowed.
func (c *Cursor[T]) IsExpired(maxAge time.Duration) bool {
	return c == nil || c.IssuedAt == 0 || time.Since(time.Unix(c.IssuedAt, 0)) > maxAge
}

// Reset resets the cursor allowing to reuse it in the same context.
func (c *Cursor[T]) Reset() {
	*c = Cursor[T]{
		Offset:  c.Offset,
		Limit:   c.Limit,
		Total:   c.Total,
		Filters: c.Filters,
	}
}

// String implements the fmt.Stringer interface.
func (c *Cursor[T]) String() string {
	b, _ := c.Encode()
	if len(b) == 0 {
		return ""
	}
	return string(b)
}

// TotalItems returns the total number of items, or -1 if unknown.
func (c *Cursor[T]) TotalItems() int {
	if c == nil || c.Total == nil {
		return notFound
	}
	return *c.Total
}

// TotalPages returns the total number of pages, or -1 if unknown.
func (c *Cursor[T]) TotalPages() int {
	if c == nil || c.Total == nil || c.Limit == 0 {
		return notFound
	}
	return int(math.Ceil(float64(*c.Total) / float64(c.Limit)))
}

func (c *Cursor[T]) isEmpty() bool {
	return c == nil || (c.Prev == nil && c.Next == nil)
}

func b64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, b64.DecodedLen(len(src)))
	n, err := b64.Decode(dst, src)
	return dst[:n], err
}

func b64Encode(src []byte) []byte {
	dst := make([]byte, b64.EncodedLen(len(src)))
	b64.Encode(dst, src)
	return dst
}
