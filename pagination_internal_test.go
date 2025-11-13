// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPaginate(t *testing.T) {
	// No parallelization here due to global variable overloading.
	now = fakeNow
	defer func() { now = time.Now }()

	const (
		limit = 2
		prev  = 1
		next  = 3
		total = 10
	)
	var (
		prv = Int64(prev)
		nxt = Int64(next)
		sum = total
	)
	for name, tc := range map[string]struct {
		// inputs
		cursor *Cursor[Int64]
		secret []byte
		// outputs
		out *Pagination
	}{
		"Default": {
			out: &Pagination{},
		},
		"Simple": {
			cursor: &Cursor[Int64]{
				Limit: limit,
				Prev:  &prv,
			},
			out: &Pagination{
				First: "eyJwcmV2IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyfQ",
				Prev:  "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyfQ",
			},
		},
		"OK": {
			cursor: &Cursor[Int64]{
				Limit: limit,
				Total: &sum,
				Prev:  &prv,
				Next:  &nxt,
				Filters: url.Values{
					"new": []string{"true"},
				},
			},
			out: &Pagination{
				First: "eyJwcmV2IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19",
				Prev:  "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19",
				Next:  "eyJuZXh0IjozLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19",
				Last:  "eyJuZXh0IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19",
			},
		},
		"Signed": {
			cursor: &Cursor[Int64]{
				Limit: limit,
				Total: &sum,
				Prev:  &prv,
				Next:  &nxt,
				Filters: url.Values{
					"new": []string{"true"},
				},
			},
			secret: []byte("ThisIsAnInsecureSecret!"),
			out: &Pagination{
				First: "eyJwcmV2IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19.",
				Prev:  "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19.",
				Next:  "eyJuZXh0IjozLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19.",
				Last:  "eyJuZXh0IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjoyLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19.",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := Paginate[Int64](tc.cursor, tc.secret)
			if err != nil {
				t.Fatal(err)
			}
			if len(tc.secret) == 0 {
				if !reflect.DeepEqual(out, tc.out) {
					t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
				}
			} else {
				if !strings.HasPrefix(out.First, tc.out.First) {
					t.Errorf("\ngot %q\nexp %q", out.First, tc.out.First)
				}
				if !strings.HasPrefix(out.Prev, tc.out.Prev) {
					t.Errorf("\ngot %q\nexp %q", out.Prev, tc.out.Prev)
				}
				if !strings.HasPrefix(out.Next, tc.out.Next) {
					t.Errorf("\ngot %q\nexp %q", out.Next, tc.out.Next)
				}
				if !strings.HasPrefix(out.Last, tc.out.Last) {
					t.Errorf("\ngot %q\nexp %q", out.Last, tc.out.Last)
				}
			}
		})
	}
}
