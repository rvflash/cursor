package cursor

import (
	"fmt"
	"strings"
)

const (
	mysqlQueryArg = "?"
	beforeExpr    = "<"
	afterExpr     = ">"
	equalExpr     = "="
)

// Statement allows building of SQL query for MySQL or MariaDB.
// The idea is to build a SQL statement like this one to going forward and
// returns results based on the cursor with descending order.
//
// WITH d AS (
//   SELECT * FROM table t WHERE cursor > ? ORDER BY cursor ASC LIMIT ?
// )
// SELECT * FROM p ORDER BY cursor DESC
//
// Limit statement adds one to the cursor's limit in order to know the start of the next cursor
// and if there is more data.
type Statement[T Pointer] struct {
	// Cursor is the cursor of pagination.
	Cursor *Cursor[T]
	// DescendingOrder defines the result's order by default.
	DescendingOrder bool
}

// Limit returns the row count to restrict the number of returned rows.
// The value is incremented by one to check if there is more to fetch.
func (s Statement[T]) Limit() int {
	if s.Cursor == nil || s.Cursor.Limit == 0 {
		return 0
	}
	return s.Cursor.Limit + 1
}

// WhereCondition returns the condition that rows must satisfy to be selected.
func (s Statement[T]) WhereCondition(columns ...string) string {
	if s.Cursor.isEmpty() {
		return ""
	}
	var p Pointer
	if s.Cursor.Next != nil {
		p = *s.Cursor.Next
	} else {
		p = *s.Cursor.Prev
	}
	if p.IsZero() {
		return ""
	}
	if len(columns) > 0 {
		buf := new(strings.Builder)
		for k := range columns {
			_, _ = fmt.Fprintf(buf, " AND %s %s %s", columns[k], s.expr(), mysqlQueryArg)
		}
		return buf.String()
	}
	return fmt.Sprintf(" %s %s", s.expr(), mysqlQueryArg)
}

// OrderBy returns the clause to order the selected and limited resultset.
// It differs from OrderBy to limit its scope to the WITH statement, also known as data source.
func (s Statement[T]) OrderBy(columns ...string) string {
	if s.Cursor.isEmpty() {
		return ""
	}
	var desc bool
	if s.Cursor.Next != nil {
		desc = s.DescendingOrder
	} else {
		desc = !s.DescendingOrder
	}
	if len(columns) > 0 {
		buf := new(strings.Builder)
		for k := range columns {
			if k > 0 {
				_, _ = fmt.Fprint(buf, ", ")
			}
			_, _ = fmt.Fprintf(buf, "%s %s", columns[k], s.orderBy(desc))
		}
		return buf.String()
	}
	return s.orderBy(desc)
}

func (s Statement[T]) expr() string {
	if s.Cursor.Next == nil {
		if s.DescendingOrder {
			return beforeExpr + equalExpr
		}
		return afterExpr + equalExpr
	}
	if s.DescendingOrder {
		return afterExpr
	}
	return beforeExpr
}

func (s Statement[T]) orderBy(desc bool) string {
	if desc {
		return " DESC"
	}
	return " ASC"
}
