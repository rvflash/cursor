// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor_test

import (
	"reflect"
	"testing"

	"github.com/rvflash/cursor"
)

const secret = "ThisIsAnInsecureSecret!"

func TestDecrypt(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		// inputs
		in []byte
		// outputs
		out *cursor.Cursor[cursor.Int64]
		msg string
	}{
		"Default": {msg: "parsing: invalid cursor format"},
		"Invalid": {
			in: []byte(
				"eyJuZXh0IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjowfQ.icJN",
			),
			msg: "signature mismatch",
		},
		"OK": {
			in: []byte(
				"eyJuZXh0IjowLCJpc3N1ZWRfYXQiOjE3NjIxMDEzMzYsImxpbWl0IjowfQ.icJNmFSIVfkw77vuW9fLZAr_L9j2e-s2HYI-SiflMRU",
			),
			out: &cursor.Cursor[cursor.Int64]{
				Next:     new(cursor.Int64),
				IssuedAt: issuedAt,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out, err := cursor.Decrypt[cursor.Int64](tc.in, []byte(secret))
			if err != nil || tc.msg != "" {
				checkErr(t, err, tc.msg)
			}
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("\ngot %#v\nexp %#v", out, tc.out)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		// inputs
		cursor *cursor.Cursor[cursor.Int64]
		secret []byte
		// outputs
		size int
		msg  string
	}{
		"Default": {},
		"OK": {
			cursor: &cursor.Cursor[cursor.Int64]{
				Next: new(cursor.Int64),
			},
			secret: []byte(secret),
			size:   102,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out, err := cursor.Encrypt(tc.cursor, tc.secret)
			if err != nil || tc.msg != "" {
				checkErr(t, err, tc.msg)
			}
			if n := len(out); n != tc.size {
				t.Errorf("\ngot %d\nexp %d", n, tc.size)
			}
		})
	}
}
