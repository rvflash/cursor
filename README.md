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

 
### Integrating with an HTTP API

Example of returning paginated results in a REST response:
```go
```