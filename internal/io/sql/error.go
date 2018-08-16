package sql

import (
	"strings"
)

// IsTableNotFound checks if an error indicates
// that a table does not exist from all known
// SQL drivers. There is unfortunately no type
// included in any of the drivers to test against.
func IsTableNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// github.com/mattn/sqlite3
	// no such table: fuubar
	if strings.Contains(msg, "no such table:") {
		return true
	}
	// github.com/lib/pq
	// pq: table "fuubar" does not exist
	if strings.Contains(msg, "pg: table") &&
		strings.Contains(msg, "does not exist") {
		return true
	}
	// github.com/go-sql-driver/mysql
	// Error 1051: Unknown table fuubar
	if strings.Contains(msg, "Error 1051: Unknown table") {
		return true
	}
	return false
}
