package cursor_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/rvflash/cursor"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

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
	order := st.OrderBy("id")

	query := `SELECT id, name FROM users` + where + " ORDER BY" + order + " LIMIT ?"

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []User
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res[:min(len(res), st.Cursor.Limit)], rows.Err()
}

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
