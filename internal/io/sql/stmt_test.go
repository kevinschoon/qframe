package sql

import (
	"testing"

	"github.com/tobgu/qframe/types"
)

func TestInsert(t *testing.T) {
	// Unescaped
	query := Insert([]string{"COL1", "COL2"}, SQLConfig{Table: "test"})
	expected := `INSERT INTO test (COL1,COL2) VALUES (?,?);`
	assertEqual(t, expected, query)

	// Double quote escaped
	query = Insert([]string{"COL1", "COL2"}, SQLConfig{
		Table: "test", EscapeChar: '"'})
	expected = "INSERT INTO \"test\" (\"COL1\",\"COL2\") VALUES (?,?);"
	assertEqual(t, expected, query)

	// Backtick escaped
	query = Insert([]string{"COL1", "COL2"}, SQLConfig{
		Table: "test", EscapeChar: '`'})
	expected = "INSERT INTO `test` (`COL1`,`COL2`) VALUES (?,?);"
	assertEqual(t, expected, query)
}

func TestCreate(t *testing.T) {
	// Unescaped
	stmt := Create(
		[]string{"COL1", "COL2"},
		[]types.DataType{types.Int, types.Float},
		SQLConfig{
			Table: "test",
			TypeMap: map[types.DataType]string{
				types.Int:   "INT",
				types.Float: "FLOAT",
			},
		},
	)
	expected := `CREATE TABLE test (COL1 INT, COL2 FLOAT);`
	assertEqual(t, expected, stmt)

	// Double quote escaped
	stmt = Create(
		[]string{"COL1", "COL2"},
		[]types.DataType{types.Int, types.Float},
		SQLConfig{
			Table: "test",
			TypeMap: map[types.DataType]string{
				types.Int:   "INT",
				types.Float: "FLOAT",
			},
			EscapeChar: '"',
		},
	)
	expected = `CREATE TABLE "test" ("COL1" "INT", "COL2" "FLOAT");`
	assertEqual(t, expected, stmt)

	// Backtick escaped
	stmt = Create(
		[]string{"COL1", "COL2"},
		[]types.DataType{types.Int, types.Float},
		SQLConfig{
			Table: "test",
			TypeMap: map[types.DataType]string{
				types.Int:   "INT",
				types.Float: "FLOAT",
			},
			EscapeChar: '`',
		},
	)
	expected = "CREATE TABLE `test` (`COL1` `INT`, `COL2` `FLOAT`);"
	assertEqual(t, expected, stmt)
}
