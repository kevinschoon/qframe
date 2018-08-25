package sql

import (
	"bytes"
	"fmt"

	"github.com/tobgu/qframe/types"
)

func escape(s string, char rune, buf *bytes.Buffer) {
	if char == 0 {
		buf.WriteString(s)
		return
	}
	buf.WriteRune(char)
	buf.WriteString(s)
	buf.WriteRune(char)
}

// Insert generates a SQL insert statement
// for each colName. There are several variations
// of SQL that need to be produced for each driver.
// This has been tested with the following:
// PostgreSQL - github.com/lib/pq
// MySQL/MariaDB - github.com/go-sql-driver/mysql
// SQLite - github.com/mattn/go-sqlite3
//
// "Parameter markers" are used to specify placeholders
// for values scanned by the implementing driver:
// PostgreSQL accepts "incrementing" markers e.g. $1..$2
// While MySQL/MariaDB and SQLite accept ?..?.
func Insert(colNames []string, conf SQLConfig) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("INSERT INTO ")
	escape(conf.Table, conf.EscapeChar, buf)
	buf.WriteString(" (")
	for i, name := range colNames {
		escape(name, conf.EscapeChar, buf)
		if i+1 < len(colNames) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(") VALUES (")
	for i := range colNames {
		if conf.Incrementing {
			buf.WriteString(fmt.Sprintf("$%d", i+1))
		} else {
			buf.WriteString("?")
		}
		if i+1 < len(colNames) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(");")
	return buf.String()
}

// Create generates a SQL CREATE statement suitable for
// dynamically creating database tables based on QFrames.
//
// SQL Types can be optionally naively overwritten by
// specifying a TypeMap with the SQLConfig.
func Create(colNames []string, colTypes []types.DataType, conf SQLConfig) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("CREATE TABLE ")
	escape(conf.Table, conf.EscapeChar, buf)
	buf.WriteString(" (")
	for i, name := range colNames {
		escape(name, conf.EscapeChar, buf)
		buf.WriteString(" ")
		buf.WriteString(conf.getSQLType(colTypes[i]))
		if i+1 < len(colNames) {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(");")
	return buf.String()
}

// Clobber generates a SQL DROP TABLE statement to dynamically
// drop a table prior to inserting data.
func Clobber(conf SQLConfig) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("DROP TABLE ")
	escape(conf.Table, conf.EscapeChar, buf)
	buf.WriteString(";")
	return buf.String()
}
