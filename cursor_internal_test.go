// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

const (
	b64Simple = "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsIk9mZnNldCI6MCwibGltaXQiOjN9"
	b64New    = "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsIk9mZnNldCI6MCwibGltaXQiOjMsInRvdGFsIjoxMCwiZmlsdGVycyI6eyJuZXciOlsidHJ1ZSJdfX0"
)

func TestCursor_Encode(t *testing.T) {
	// No parallelization here due to global variable overloading.
	now = fakeNow
	defer func() { now = time.Now }()

	const (
		limit = 3
		prev  = 1
		total = 10
	)
	var (
		prv = Int64(prev)
		sum = total
	)
	for name, tc := range map[string]struct {
		in  Cursor[Int64]
		out []byte
	}{
		"Default": {},
		"Simple": {
			in: Cursor[Int64]{
				Limit: limit,
				Prev:  &prv,
			},
			out: []byte(b64Simple),
		},
		"OK": {
			in: Cursor[Int64]{
				Limit: limit,
				Total: &sum,
				Prev:  &prv,
				Filters: url.Values{
					"new": []string{"true"},
				},
			},
			out: []byte(b64New),
		},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := tc.in.Encode()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %q\nexp %q", out, tc.out)
			}
		})
	}
}

func TestCursor_String(t *testing.T) {
	// No parallelization here due to global variable overloading.
	now = fakeNow
	defer func() { now = time.Now }()

	const (
		limit = 3
		prev  = 1
		total = 10
	)
	var (
		prv = Int64(prev)
		sum = total
	)
	for name, tc := range map[string]struct {
		in  Cursor[Int64]
		out string
	}{
		"Default": {},
		"Simple": {
			in: Cursor[Int64]{
				Limit: limit,
				Prev:  &prv,
			},
			out: b64Simple,
		},
		"OK": {
			in: Cursor[Int64]{
				Limit: limit,
				Total: &sum,
				Prev:  &prv,
				Filters: url.Values{
					"new": []string{"true"},
				},
			},
			out: b64New,
		},
	} {
		t.Run(name, func(t *testing.T) {
			out := tc.in.String()
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %q\nexp %q", out, tc.out)
			}
		})
	}
}

func fakeNow() time.Time {
	return time.Unix(1762101336, 0)
}
