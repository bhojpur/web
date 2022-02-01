package orm

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"strconv"
)

// postgresql operators.
var postgresOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "= UPPER(?)",
	"contains":    "LIKE ?",
	"icontains":   "LIKE UPPER(?)",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"eq":          "= ?",
	"ne":          "!= ?",
	"startswith":  "LIKE ?",
	"endswith":    "LIKE ?",
	"istartswith": "LIKE UPPER(?)",
	"iendswith":   "LIKE UPPER(?)",
}

// postgresql column field types.
var postgresTypes = map[string]string{
	"auto":                "serial NOT NULL PRIMARY KEY",
	"pk":                  "NOT NULL PRIMARY KEY",
	"bool":                "bool",
	"string":              "varchar(%d)",
	"string-char":         "char(%d)",
	"string-text":         "text",
	"time.Time-date":      "date",
	"time.Time":           "timestamp with time zone",
	"int8":                `smallint CHECK("%COL%" >= -127 AND "%COL%" <= 128)`,
	"int16":               "smallint",
	"int32":               "integer",
	"int64":               "bigint",
	"uint8":               `smallint CHECK("%COL%" >= 0 AND "%COL%" <= 255)`,
	"uint16":              `integer CHECK("%COL%" >= 0)`,
	"uint32":              `bigint CHECK("%COL%" >= 0)`,
	"uint64":              `bigint CHECK("%COL%" >= 0)`,
	"float64":             "double precision",
	"float64-decimal":     "numeric(%d, %d)",
	"json":                "json",
	"jsonb":               "jsonb",
	"time.Time-precision": "timestamp(%d) with time zone",
}

// PostgreSQL dbBaser.
type dbBasePostgres struct {
	dbBase
}

var _ dbBaser = new(dbBasePostgres)

// get postgresql operator.
func (d *dbBasePostgres) OperatorSQL(operator string) string {
	return postgresOperators[operator]
}

// generate functioned sql string, such as contains(text).
func (d *dbBasePostgres) GenerateOperatorLeftCol(fi *fieldInfo, operator string, leftCol *string) {
	switch operator {
	case "contains", "startswith", "endswith":
		*leftCol = fmt.Sprintf("%s::text", *leftCol)
	case "iexact", "icontains", "istartswith", "iendswith":
		*leftCol = fmt.Sprintf("UPPER(%s::text)", *leftCol)
	}
}

// postgresql unsupports updating joined record.
func (d *dbBasePostgres) SupportUpdateJoin() bool {
	return false
}

func (d *dbBasePostgres) MaxLimit() uint64 {
	return 0
}

// postgresql quote is ".
func (d *dbBasePostgres) TableQuote() string {
	return `"`
}

// postgresql value placeholder is $n.
// replace default ? to $n.
func (d *dbBasePostgres) ReplaceMarks(query *string) {
	q := *query
	num := 0
	for _, c := range q {
		if c == '?' {
			num++
		}
	}
	if num == 0 {
		return
	}
	data := make([]byte, 0, len(q)+num)
	num = 1
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '?' {
			data = append(data, '$')
			data = append(data, []byte(strconv.Itoa(num))...)
			num++
		} else {
			data = append(data, c)
		}
	}
	*query = string(data)
}

// make returning sql support for postgresql.
func (d *dbBasePostgres) HasReturningID(mi *modelInfo, query *string) bool {
	fi := mi.fields.pk
	if fi.fieldType&IsPositiveIntegerField == 0 && fi.fieldType&IsIntegerField == 0 {
		return false
	}

	if query != nil {
		*query = fmt.Sprintf(`%s RETURNING "%s"`, *query, fi.column)
	}
	return true
}

// sync auto key
func (d *dbBasePostgres) setval(db dbQuerier, mi *modelInfo, autoFields []string) error {
	if len(autoFields) == 0 {
		return nil
	}

	Q := d.ins.TableQuote()
	for _, name := range autoFields {
		query := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', '%s'), (SELECT MAX(%s%s%s) FROM %s%s%s));",
			mi.table, name,
			Q, name, Q,
			Q, mi.table, Q)
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// show table sql for postgresql.
func (d *dbBasePostgres) ShowTablesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema NOT IN ('pg_catalog', 'information_schema')"
}

// show table columns sql for postgresql.
func (d *dbBasePostgres) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns where table_schema NOT IN ('pg_catalog', 'information_schema') and table_name = '%s'", table)
}

// get column types of postgresql.
func (d *dbBasePostgres) DbTypes() map[string]string {
	return postgresTypes
}

// check index exist in postgresql.
func (d *dbBasePostgres) IndexExists(db dbQuerier, table string, name string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM pg_indexes WHERE tablename = '%s' AND indexname = '%s'", table, name)
	row := db.QueryRow(query)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

// GenerateSpecifyIndex return a specifying index clause
func (d *dbBasePostgres) GenerateSpecifyIndex(tableName string, useIndex int, indexes []string) string {
	DebugLog.Println("[WARN] Not support any specifying index action, so that action is ignored")
	return ``
}

// create new postgresql dbBaser.
func newdbBasePostgres() dbBaser {
	b := new(dbBasePostgres)
	b.ins = b
	return b
}
