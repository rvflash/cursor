# Cursor

[![GoDoc](https://godoc.org/github.com/rvflash/cursor?status.svg)](https://godoc.org/github.com/rvflash/cursor)
[![Build Status](https://github.com/rvflash/cursor/workflows/build/badge.svg)](https://github.com/rvflash/cursor/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/cursor/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/cursor)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/cursor?)](https://goreportcard.com/report/github.com/rvflash/cursor)


A lightweight, generic cursor-based pagination package for Go.
Designed for MySQL and MariaDB, `cursor` lets you build encrypted, 
stateless cursors that encode pagination state and query parameters safely.

Cursor-based pagination is preferred to OFFSET for better performance on large tables.
It stores the last seen ID or timestamp in the cursor to build efficient WHERE clauses.
It is recommended to check the expiration date with `IsExpired(maxAge)` to avoid using outdated cursors.

3 main types:
- `Cursor` allows computation of data necessary for pagination.
- `Statement` builds based on a Cursor SQL query parts, to use to perform a SELECT statement.
- `Pointer` defines the data types that can be used as a cursor to filter the query.
Such as `Int64` to manage the auto-increment field.
See also `String` or List to manage a set of `Pointer` as `Pointer`.
Finally, `Limit` can be used as `Pointer` to transform the cursor into a standard LIMIT statement, with offset and row count.


### Cursor Encoding Format

- Internally serialized as JSON, then encoded as Base64 (URL-safe).
- Optionally encrypted or signed using HMAC for integrity.
- Fully stateless — no server session needed.


## Features

- 🔒 Encrypted cursors — opaque Base64 tokens with or not HMAC signing.
- 📜 Cursor-based pagination — no offset drift, efficient for large datasets.
- 🧠 Stateless by design — all state is encoded in the cursor.
- 💡 Generic — supports any data type T.
- 🧩 SQL helpers for LIMIT, ORDER BY, and conditional pagination queries.
- ⏱️ Expiration support — cursors can self-expire based on max age.


## Installation

```go
go get github.com/rvflash/cursor
```

## Example Usage

Codes with multiple shortcuts for demonstration purposes only.

### SQL statement


```go
func ListFromDatabase(ctx context.Context, cur *cursor.Cursor[cursor.Int64]) ([]User, error) {
    // Create a Statement based on the Cursor.
    var st = cursor.Statement[cursor.Int64]{
        Cursor:          cur,
        DescendingOrder: false,
    }
    // WHERE uses the cursor semantics (e.g., "id < ?") under descending order
    where, args := st.WhereCondition("id")
    if len(args) > 0 {
        where = " WHERE " + where
    }
    // LIMIT +1 to check if there is a next page.
    args = append(args, st.Limit())
    // ORDER BY applies the desired ordering for the limited page (e.g., "ORDER BY id DESC")
    query := `SELECT id, name FROM users` + where + " ORDER BY" + st.OrderBy("id") + " LIMIT ?"
    
    rows, err := DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer func() { _ = rows.Close() }()
    
    var (
		res []User
        u User
	)
    for rows.Next() {
        if err = rows.Scan(&u.ID, &u.Name); err != nil {
            return nil, err
        }
        res = append(res, u)
    }
    return res[:min(len(res), st.Cursor.Limit)], rows.Err()
}

type User struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
}
```
 
### Integrating with an HTTP API

Example of returning paginated results in a REST response:
```go
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
    var (
        secret = []byte(os.Getenv("CURSOR_SECRET"))
        cur    *cursor.Cursor[cursor.Int64]
        err    error
    )
    if tok := r.URL.Query().Get("cursor"); tok != "" {
        // Decrypt verifies HMAC and returns the cursor state
        cur, err = cursor.Decrypt[cursor.Int64]([]byte(tok), secret)
        if err != nil || cur.IsExpired(time.Hour) {
         http.Error(w, "invalid or expired cursor", http.StatusBadRequest)
            return
        }
	} else {
        // New(limit, total). If you don’t know total, you can pass 0 (or compute it separately)
        cur = cursor.New[cursor.Int64](20, 0)
    }
    // SQL query
    rows, err := ListFromDatabase(r.Context(), cur)
    if err != nil {
        http.Error(w, "query error", http.StatusInternalServerError)
        return
    }
    // Reset allows to reuse the current cursor to build the next ones.
    cur.Reset()
    for _, u := range rows {
        cur.Add(cursor.Int64(u.ID)) // we’re pointing by ID in this example
    }
    // Build pagination tokens: first/prev/next/last.
    pg, err := cursor.Paginate(cur, secret)
    if err != nil {
        http.Error(w, "pagination error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    _ = json.NewEncoder(w).Encode(usersResponse{
        Data:       rows,
        Pagination: pg,
    })
}
```