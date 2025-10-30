package cursor_test

import (
	"fmt"
	"time"

	"github.com/rvflash/cursor"
)

func ExampleStatement_WithOrderBy() {
	// First cursor
	var (
		s = cursor.Statement[time.Time]{
			Cursor: &cursor.Cursor[time.Time]{
				Prev:  &time.Time{},
				Limit: 20,
			},
			ColumnName:      "creation_date",
			DescendingOrder: true,
		}
		q = `WITH d AS (SELECT * FROM request t %sORDER BY %s LIMIT %d) SELECT * FROM d ORDER BY %s`
	)
	fmt.Printf(q, s.WithWhereCondition(), s.WithOrderBy(), s.Limit(), s.OrderBy())

	// Output:
	// WITH d AS (SELECT * FROM request t ORDER BY creation_date ASC LIMIT 21) SELECT * FROM d ORDER BY creation_date DESC
}

func ExampleStatement_WithWhereCondition() {
	// Next cursor
	var (
		d = time.Unix(1761689506, 0)
		s = cursor.Statement[time.Time]{
			Cursor: &cursor.Cursor[time.Time]{
				Next:  &d,
				Limit: 20,
			},
			ColumnName:      "creation_date",
			DescendingOrder: true,
		}
		q = `WITH d AS (SELECT * FROM request t WHERE %s ORDER BY %s LIMIT %d) SELECT * FROM d ORDER BY %s`
	)
	fmt.Printf(q, s.WithWhereCondition(), s.WithOrderBy(), s.Limit(), s.OrderBy())

	// Output:
	// WITH d AS (SELECT * FROM request t WHERE creation_date > ? ORDER BY creation_date DESC LIMIT 21) SELECT * FROM d ORDER BY creation_date DESC
}
