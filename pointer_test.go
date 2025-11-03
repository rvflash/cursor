package cursor_test

import (
	"reflect"
	"testing"

	"github.com/rvflash/cursor"
)

func TestInt64_Args(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.Int64
		out []any
	}{
		"Default":  {},
		"Negative": {in: -1, out: []any{cursor.Int64(-1)}},
		"Positive": {in: 1, out: []any{cursor.Int64(1)}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.Args()
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestInt64_IsZero(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.Int64
		out bool
	}{
		"Default":  {out: true},
		"Negative": {in: -1},
		"Positive": {in: 1},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.IsZero()
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestList_Args(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.List
		out []any
	}{
		"Default": {},
		"Blank":   {in: cursor.List{}},
		"Zero":    {in: cursor.List{cursor.Int64(0)}, out: []any{nil}},
		"Partial": {in: cursor.List{cursor.Int64(1), cursor.Int64(0)}, out: []any{cursor.Int64(1), nil}},
		"OK":      {in: cursor.List{cursor.Int64(1)}, out: []any{cursor.Int64(1)}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.Args()
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestList_IsZero(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.List
		out bool
	}{
		"Default": {out: true},
		"Blank":   {in: cursor.List{}, out: true},
		"Zero":    {in: cursor.List{cursor.Int64(0)}, out: true},
		"OK":      {in: cursor.List{cursor.Int64(1)}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.IsZero()
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestString_Args(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.String
		out []any
	}{
		"Default": {},
		"Blank":   {in: cursor.String(" "), out: []any{cursor.String(" ")}},
		"OK":      {in: cursor.String("hi"), out: []any{cursor.String("hi")}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.Args()
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestString_IsZero(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.String
		out bool
	}{
		"Default": {out: true},
		"Blank":   {in: cursor.String(" ")},
		"OK":      {in: cursor.String("hi")},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.IsZero()
			if out != tc.out {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}
