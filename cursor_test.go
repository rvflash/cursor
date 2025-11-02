package cursor_test

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/rvflash/cursor"
)

const (
	limit    = 2
	prev     = 1
	next     = 3
	total    = 10
	issuedAt = 1762101336
)

func TestFirst(t *testing.T) {
	t.Parallel()

	var (
		prv = cursor.Int64(prev)
		nxt = cursor.Int64(next)
	)
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out *cursor.Cursor[cursor.Int64]
	}{
		"Default": {},
		"Nothing before": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
		},
		"Nothing after": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
		"First page": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
		"Last page": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
		"OK": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := cursor.First(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestLast(t *testing.T) {
	t.Parallel()

	var (
		prv = cursor.Int64(prev)
		nxt = cursor.Int64(next)
	)
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out *cursor.Cursor[cursor.Int64]
	}{
		"Default": {},
		"Nothing before": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
		"Nothing after": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
		},
		"First page": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
		"Last page": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
		"OK": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := cursor.Last(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestNext(t *testing.T) {
	t.Parallel()

	var (
		prv = cursor.Int64(prev)
		nxt = cursor.Int64(next)
	)
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out *cursor.Cursor[cursor.Int64]
	}{
		"Default": {},
		"Nothing before": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
		},
		"Nothing after": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
		},
		"First page": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
		"Last page": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
		"OK": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := cursor.Next(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestPrev(t *testing.T) {
	t.Parallel()

	var (
		prv = cursor.Int64(prev)
		nxt = cursor.Int64(next)
	)
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out *cursor.Cursor[cursor.Int64]
	}{
		"Default": {},
		"Nothing before": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: &nxt,
			},
		},
		"Nothing after": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
		},
		"First page": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: new(cursor.Int64),
			},
		},
		"Last page": {
			in: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
		},
		"OK": {
			in: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
				Next: &nxt,
			},
			out: &cursor.Cursor[cursor.Int64]{
				Prev: &prv,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := cursor.Prev(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	sum := total
	for name, tc := range map[string]struct {
		// inputs
		limit int
		total int
		// outputs
		out *cursor.Cursor[cursor.Int64]
	}{
		"Default": {out: &cursor.Cursor[cursor.Int64]{
			Filters: make(url.Values),
		}},
		"Limit only": {
			limit: limit,
			out: &cursor.Cursor[cursor.Int64]{
				Limit:   limit,
				Filters: make(url.Values),
			},
		},
		"Complete": {
			limit: limit,
			total: sum,
			out: &cursor.Cursor[cursor.Int64]{
				Limit:   limit,
				Total:   &sum,
				Filters: make(url.Values),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := cursor.New[cursor.Int64](tc.limit, tc.total)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestCursor_Add(t *testing.T) {
	t.Parallel()

	var (
		prv = cursor.Int64(prev)
		nxt = cursor.Int64(next)
	)
	for name, tc := range map[string]struct {
		// inputs
		in  cursor.Cursor[cursor.Int64]
		num int
		// outputs
		out cursor.Cursor[cursor.Int64]
	}{
		"Default": {num: next},
		"No iteration": {
			in: cursor.Cursor[cursor.Int64]{
				Limit: limit,
			},
			out: cursor.Cursor[cursor.Int64]{
				Limit: limit,
			},
		},
		"No more": {
			in: cursor.Cursor[cursor.Int64]{
				Limit: limit,
			},
			num: prev,
			out: cursor.Cursor[cursor.Int64]{
				Limit: limit,
				Prev:  &prv,
			},
		},
		"OK": {
			in: cursor.Cursor[cursor.Int64]{
				Limit: limit,
			},
			num: next,
			out: cursor.Cursor[cursor.Int64]{
				Limit: limit,
				Prev:  &prv,
				Next:  &nxt,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for i := 1; i <= tc.num; i++ {
				tc.in.Add(cursor.Int64(i))
			}
			if !reflect.DeepEqual(tc.in.Prev, tc.out.Prev) {
				t.Errorf("\ngot %#v\nexp %#v", tc.in.Prev, tc.out.Prev)
			}
			if !reflect.DeepEqual(tc.in.Next, tc.out.Next) {
				t.Errorf("\ngot %#v\nexp %#v", tc.in.Next, tc.out.Next)
			}
		})
	}
}

func TestCursor_Decode(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		// inputs
		got cursor.Cursor[cursor.Int64]
		in  []byte
		// outputs
		exp cursor.Cursor[cursor.Int64]
		msg string
	}{
		"Default": {
			msg: "unmarshalling: unexpected end",
		},
		"Partial": {
			in: []byte("eyJuZXh0IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjowfQ"),
			exp: cursor.Cursor[cursor.Int64]{
				Next:     new(cursor.Int64),
				IssuedAt: issuedAt,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.got.Decode(tc.in)
			if err != nil || tc.msg != "" {
				checkErr(t, err, tc.msg)
			}
			if !reflect.DeepEqual(tc.got, tc.exp) {
				t.Errorf("\ngot %#v\nexp %#v", tc.got, tc.exp)
			}
		})
	}
}

func TestCursor_IsExpired(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		// inputs
		in     *cursor.Cursor[cursor.Int64]
		maxAge time.Duration
		// outputs
		out bool
	}{
		"Default": {out: true},
		"Blank":   {in: &cursor.Cursor[cursor.Int64]{}, out: true},
		"Expired": {
			in: &cursor.Cursor[cursor.Int64]{
				IssuedAt: issuedAt,
			},
			maxAge: time.Second,
			out:    true,
		},
		"Not expired": {
			in: &cursor.Cursor[cursor.Int64]{
				IssuedAt: issuedAt,
			},
			maxAge: 12 * time.Hour,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.IsExpired(tc.maxAge)
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestCursor_Reset(t *testing.T) {
	t.Parallel()

	var (
		sum = total
		nxt = cursor.Int64(next)
		got = cursor.Cursor[cursor.Int64]{
			Prev:     new(cursor.Int64),
			Next:     &nxt,
			IssuedAt: issuedAt,
			Limit:    limit,
			Total:    &sum,
			Filters:  url.Values{"new": []string{"true"}},
		}
		exp = cursor.Cursor[cursor.Int64]{
			Limit:   limit,
			Total:   &sum,
			Filters: url.Values{"new": []string{"true"}},
		}
	)
	got.Reset()

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\ngot %#v\nexp %#v", got, exp)
	}
}

func TestCursor_TotalItems(t *testing.T) {
	t.Parallel()

	sum := total
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out int
	}{
		"Default": {out: -1},
		"Blank":   {in: &cursor.Cursor[cursor.Int64]{}, out: -1},
		"OK":      {in: &cursor.Cursor[cursor.Int64]{Total: &sum}, out: total},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.TotalItems()
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestCursor_TotalPages(t *testing.T) {
	t.Parallel()

	sum := total
	for name, tc := range map[string]struct {
		in  *cursor.Cursor[cursor.Int64]
		out int
	}{
		"Default":       {out: -1},
		"Blank":         {in: &cursor.Cursor[cursor.Int64]{}, out: -1},
		"Total missing": {in: &cursor.Cursor[cursor.Int64]{Limit: limit}, out: -1},
		"OK":            {in: &cursor.Cursor[cursor.Int64]{Limit: limit, Total: &sum}, out: total / 2},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.TotalPages()
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func checkErr(t *testing.T, err error, substr string) {
	t.Helper()

	if substr == "" {
		t.Errorf("unexpected error: %s", err.Error())
	} else if err == nil {
		t.Errorf("got = nil, exp error: %s", substr)
	} else if !strings.Contains(err.Error(), substr) {
		t.Errorf("got = %s, exp = %s", err.Error(), substr)
	}
}
