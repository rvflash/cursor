package cursor_test

import (
	"reflect"
	"testing"

	"github.com/rvflash/cursor"
)

// Statement building, step by step.
//
// Base on the following dataset, with a cursor limit set to 2, so a statement limit of 3.
//
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52405 | 2025-10-31 11:38:37 |
// | 52404 | 2025-10-30 16:58:32 |
// | 52352 | 2025-10-30 16:17:12 |
// | 52351 | 2025-10-30 16:15:46 |
// | 52350 | 2025-10-30 16:13:50 |
// | 52349 | 2025-10-30 16:11:47 |
// | 52348 | 2025-10-30 11:36:49 |
// | 52320 | 2025-10-29 17:24:13 |
// | 52319 | 2025-10-29 17:11:53 |
// | 52318 | 2025-10-29 16:35:25 |
// ...
// +-------+---------------------+
//
// 3 items per page, descending order.
//
// First page
// > SELECT id, creation_date FROM data ORDER BY id DESC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52405 | 2025-10-31 11:38:37 |
// | 52404 | 2025-10-30 16:58:32 |
// | 52352 | 2025-10-30 16:17:12 |
// +-------+---------------------+
//
// Prev page (no content)
//
// Next page
// > SELECT id, creation_date FROM data WHERE id <= 52352 ORDER BY id DESC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52352 | 2025-10-30 16:17:12 |
// | 52351 | 2025-10-30 16:15:46 |
// | 52350 | 2025-10-30 16:13:50 |
// +-------+---------------------+
//
// Last page
// > SELECT id, creation_date FROM data ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10000 | 2024-02-19 12:14:42 |
// | 10001 | 2024-02-19 13:40:14 |
// | 10012 | 2024-02-19 14:27:15 |
// +-------+---------------------+
//
// Based one page 2
//
// First page
// > SELECT id, creation_date FROM data ORDER BY id DESC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52405 | 2025-10-31 11:38:37 |
// | 52404 | 2025-10-30 16:58:32 |
// | 52352 | 2025-10-30 16:17:12 |
// +-------+---------------------+
//
// Prev page
// > SELECT id, creation_date FROM data WHERE id > 52352 ORDER BY id DESC LIMIT 3
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52405 | 2025-10-31 11:38:37 |
// | 52404 | 2025-10-30 16:58:32 |
// +-------+---------------------+
//
// Next page
// > SELECT id, creation_date FROM data WHERE id <= 52350 ORDER BY id DESC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52350 | 2025-10-30 16:13:50 |
// | 52349 | 2025-10-30 16:11:47 |
// | 52348 | 2025-10-30 11:36:49 |
// +-------+---------------------+
//
// Last page
// > SELECT id, creation_date FROM data ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10000 | 2024-02-19 12:14:42 |
// | 10001 | 2024-02-19 13:40:14 |
// | 10012 | 2024-02-19 14:27:15 |
// +-------+---------------------+
//
// And with ascending sort order:
//
// > SELECT id, creation_date FROM data ORDER BY id ASC LIMIT 10;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10000 | 2024-02-19 12:14:42 |
// | 10001 | 2024-02-19 13:40:14 |
// | 10012 | 2024-02-19 14:27:15 |
// | 10013 | 2024-02-19 14:50:28 |
// | 10046 | 2024-02-19 15:01:46 |
// | 10071 | 2024-02-19 15:11:26 |
// | 10072 | 2024-02-19 15:27:37 |
// | 10073 | 2024-02-19 15:33:38 |
// | 10074 | 2024-02-19 15:49:29 |
// | 10097 | 2024-02-19 16:17:30 |
// +-------+---------------------+
//
// Based one page 2:
//
// > SELECT id, creation_date FROM data WHERE id >= 10012 ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10012 | 2024-02-19 14:27:15 |
// | 10013 | 2024-02-19 14:50:28 |
// | 10046 | 2024-02-19 15:01:46 |
// +-------+---------------------+
//
// First page
// > SELECT id, creation_date FROM data ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10000 | 2024-02-19 12:14:42 |
// | 10001 | 2024-02-19 13:40:14 |
// | 10012 | 2024-02-19 14:27:15 |
// +-------+---------------------+
//
// Prev page
// > SELECT id, creation_date FROM data WHERE id < 10012 ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10000 | 2024-02-19 12:14:42 |
// | 10001 | 2024-02-19 13:40:14 |
// +-------+---------------------+
//
// Next page
// > SELECT id, creation_date FROM data WHERE id >= 10046 ORDER BY id ASC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 10046 | 2024-02-19 15:01:46 |
// | 10071 | 2024-02-19 15:11:26 |
// | 10072 | 2024-02-19 15:27:37 |
// +-------+---------------------+
//
// Last page
// > SELECT id, creation_date FROM data ORDER BY id DESC LIMIT 3;
// +-------+---------------------+
// | id    | creation_date       |
// +-------+---------------------+
// | 52406 | 2025-11-03 11:31:16 |
// | 52405 | 2025-10-31 11:38:37 |
// | 52404 | 2025-10-30 16:58:32 |
// +-------+---------------------+

const (
	p1DescKey = 52405
	p3DescKey = 52352
	p1AscKey  = 10000
	p3AscKey  = 10012
)

func TestStatement_Limit(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		in  cursor.Statement[cursor.Int64]
		out int
	}{
		"Default": {},
		"Blank": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{},
			},
		},
		"OK": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
				},
			},
			out: limit + 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.Limit()
			if out != tc.out {
				t.Errorf("\ngot %d\nexp %d", out, tc.out)
			}
		})
	}
}

func TestStatement_OrderBy(t *testing.T) {
	t.Parallel()

	var (
		descP1K = cursor.Int64(p1DescKey)
		descP3K = cursor.Int64(p3DescKey)
		ascP1K  = cursor.Int64(p1AscKey)
		ascP3K  = cursor.Int64(p3AscKey)
	)
	for name, tc := range map[string]struct {
		// inputs
		in   cursor.Statement[cursor.Int64]
		cols []string
		// outputs
		out string
	}{
		"Default": {out: " ASC"},
		"Descending - No cursor": {
			in: cursor.Statement[cursor.Int64]{
				DescendingOrder: true,
			},
			out: " DESC",
		},
		"Descending - First page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			out: " DESC",
		},
		"Descending - First page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01"},
			out:  " col01 DESC",
		},
		"Descending - First page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 DESC, a.col02 DESC",
		},
		"Descending - Prev page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &descP1K,
				},
				DescendingOrder: true,
			},
			out: " DESC",
		},
		"Descending - Prev page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &descP1K,
				},
				DescendingOrder: true,
			},
			cols: []string{"col01"},
			out:  " col01 DESC",
		},
		"Descending - Prev page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &descP1K,
				},
				DescendingOrder: true,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 DESC, a.col02 DESC",
		},
		"Descending - Next page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  &descP3K,
				},
				DescendingOrder: true,
			},
			out: " DESC",
		},
		"Descending - Last page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Next: new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			out: " ASC",
		},
		"Descending - Last page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01"},
			out:  " col01 ASC",
		},
		"Descending - Last page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 ASC, a.col02 ASC",
		},
		"Ascending - No cursor": {
			in: cursor.Statement[cursor.Int64]{
				DescendingOrder: false,
			},
			out: " ASC",
		},
		"Ascending - First page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			out: " ASC",
		},
		"Ascending - First page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01"},
			out:  " col01 ASC",
		},
		"Ascending - First page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 ASC, a.col02 ASC",
		},
		"Ascending - Prev page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &ascP1K,
				},
				DescendingOrder: false,
			},
			out: " ASC",
		},
		"Ascending - Prev page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &ascP1K,
				},
				DescendingOrder: false,
			},
			cols: []string{"col01"},
			out:  " col01 ASC",
		},
		"Ascending - Prev page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &ascP1K,
				},
				DescendingOrder: false,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 ASC, a.col02 ASC",
		},
		"Ascending - Next page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  &ascP3K,
				},
				DescendingOrder: false,
			},
			out: " ASC",
		},
		"Ascending - Last page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			out: " DESC",
		},
		"Ascending - Last page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01"},
			out:  " col01 DESC",
		},
		"Ascending - Last page - Some columns": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01", "a.col02"},
			out:  " col01 DESC, a.col02 DESC",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := tc.in.OrderBy(tc.cols...)
			if out != tc.out {
				t.Errorf("\ngot %s\nexp %s", out, tc.out)
			}
		})
	}
}

func TestStatement_WhereCondition(t *testing.T) {
	t.Parallel()

	var (
		descP1K = cursor.Int64(p1DescKey)
		descP3K = cursor.Int64(p3DescKey)
		ascP1K  = cursor.Int64(p1AscKey)
		ascP3K  = cursor.Int64(p3AscKey)
	)
	for name, tc := range map[string]struct {
		// inputs
		in   cursor.Statement[cursor.Int64]
		cols []string
		// outputs
		query string
		args  []any
	}{
		"Default": {},
		"Descending - First page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
		},
		"Descending - First page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01"},
		},
		"Descending - Prev page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &descP1K,
				},
				DescendingOrder: true,
			},
			query: " > ?",
			args:  []any{descP1K},
		},
		"Descending - Prev page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &descP1K,
				},
				DescendingOrder: true,
			},
			cols:  []string{"col01"},
			query: " AND col01 > ?",
			args:  []any{descP1K},
		},
		"Descending - Next page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  &descP3K,
				},
				DescendingOrder: true,
			},
			query: " <= ?",
			args:  []any{descP3K},
		},
		"Descending - Last page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
		},
		"Descending - Last page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: true,
			},
			cols: []string{"col01"},
		},
		"Ascending - First page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
		},
		"Ascending - First page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01"},
		},
		"Ascending - Prev page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &ascP1K,
				},
				DescendingOrder: false,
			},
			query: " < ?",
			args:  []any{ascP1K},
		},
		"Ascending - Prev page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Prev:  &ascP1K,
				},
				DescendingOrder: false,
			},
			cols:  []string{"col01"},
			query: " AND col01 < ?",
			args:  []any{ascP1K},
		},
		"Ascending - Next page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  &ascP3K,
				},
				DescendingOrder: false,
			},
			query: " >= ?",
			args:  []any{ascP3K},
		},
		"Ascending - Next page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  &ascP3K,
				},
				DescendingOrder: false,
			},
			cols:  []string{"col01"},
			query: " AND col01 >= ?",
			args:  []any{ascP3K},
		},
		"Ascending - Last page - No column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
		},
		"Ascending - Last page - One column": {
			in: cursor.Statement[cursor.Int64]{
				Cursor: &cursor.Cursor[cursor.Int64]{
					Limit: limit,
					Next:  new(cursor.Int64),
				},
				DescendingOrder: false,
			},
			cols: []string{"col01"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, args := tc.in.WhereCondition(tc.cols...)
			if query != tc.query {
				t.Errorf("\ngot %s\nexp %s", query, tc.query)
			}
			if !reflect.DeepEqual(args, tc.args) {
				t.Errorf("\ngot %#v\nexp %#v", args, tc.args)
			}
		})
	}
}
