package cursor

import (
	"net/url"
	"reflect"
	"testing"
	"time"
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
			out: []byte("eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjozfQ"),
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
			out: []byte("eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjozLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19"),
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
			out: "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjozfQ",
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
			out: "eyJwcmV2IjoxLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjozLCJ0b3RhbCI6MTAsImZpbHRlcnMiOnsibmV3IjpbInRydWUiXX19",
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
