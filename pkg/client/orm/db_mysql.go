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
	"reflect"
	"strings"
)

// mysql operators.
var mysqlOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "LIKE ?",
	"strictexact": "= BINARY ?",
	"contains":    "LIKE BINARY ?",
	"icontains":   "LIKE ?",
	// "regex":       "REGEXP BINARY ?",
	// "iregex":      "REGEXP ?",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"eq":          "= ?",
	"ne":          "!= ?",
	"startswith":  "LIKE BINARY ?",
	"endswith":    "LIKE BINARY ?",
	"istartswith": "LIKE ?",
	"iendswith":   "LIKE ?",
}

// mysql column field types.
var mysqlTypes = map[string]string{
	"auto":                "AUTO_INCREMENT NOT NULL PRIMARY KEY",
	"pk":                  "NOT NULL PRIMARY KEY",
	"bool":                "bool",
	"string":              "varchar(%d)",
	"string-char":         "char(%d)",
	"string-text":         "longtext",
	"time.Time-date":      "date",
	"time.Time":           "datetime",
	"int8":                "tinyint",
	"int16":               "smallint",
	"int32":               "integer",
	"int64":               "bigint",
	"uint8":               "tinyint unsigned",
	"uint16":              "smallint unsigned",
	"uint32":              "integer unsigned",
	"uint64":              "bigint unsigned",
	"float64":             "double precision",
	"float64-decimal":     "numeric(%d, %d)",
	"time.Time-precision": "datetime(%d)",
}

// mysql dbBaser implementation.
type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

// get mysql operator.
func (d *dbBaseMysql) OperatorSQL(operator string) string {
	return mysqlOperators[operator]
}

// get mysql table field types.
func (d *dbBaseMysql) DbTypes() map[string]string {
	return mysqlTypes
}

// show table sql for mysql.
func (d *dbBaseMysql) ShowTablesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = DATABASE()"
}

// show columns sql of table for mysql.
func (d *dbBaseMysql) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.columns "+
		"WHERE table_schema = DATABASE() AND table_name = '%s'", table)
}

// execute sql to check index exist.
func (d *dbBaseMysql) IndexExists(db dbQuerier, table string, name string) bool {
	row := db.QueryRow("SELECT count(*) FROM information_schema.statistics "+
		"WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", table, name)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

// InsertOrUpdate a row
// If your primary key or unique column conflict will update
// If no will insert
// Add "`" for mysql sql building
func (d *dbBaseMysql) InsertOrUpdate(q dbQuerier, mi *modelInfo, ind reflect.Value, a *alias, args ...string) (int64, error) {
	var iouStr string
	argsMap := map[string]string{}

	iouStr = "ON DUPLICATE KEY UPDATE"

	//Get on the key-value pairs
	for _, v := range args {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			argsMap[strings.ToLower(kv[0])] = kv[1]
		}
	}

	isMulti := false
	names := make([]string, 0, len(mi.fields.dbcols)-1)
	Q := d.ins.TableQuote()
	values, _, err := d.collectValues(mi, ind, mi.fields.dbcols, true, true, &names, a.TZ)

	if err != nil {
		return 0, err
	}

	marks := make([]string, len(names))
	updateValues := make([]interface{}, 0)
	updates := make([]string, len(names))

	for i, v := range names {
		marks[i] = "?"
		valueStr := argsMap[strings.ToLower(v)]
		if valueStr != "" {
			updates[i] = "`" + v + "`" + "=" + valueStr
		} else {
			updates[i] = "`" + v + "`" + "=?"
			updateValues = append(updateValues, values[i])
		}
	}

	values = append(values, updateValues...)

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	qupdates := strings.Join(updates, ", ")
	columns := strings.Join(names, sep)

	multi := len(values) / len(names)

	if isMulti {
		qmarks = strings.Repeat(qmarks+"), (", multi-1) + qmarks
	}
	//conflitValue maybe is a int,can`t use fmt.Sprintf
	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s) %s "+qupdates, Q, mi.table, Q, Q, columns, Q, qmarks, iouStr)

	d.ins.ReplaceMarks(&query)

	if isMulti || !d.ins.HasReturningID(mi, &query) {
		res, err := q.Exec(query, values...)
		if err == nil {
			if isMulti {
				return res.RowsAffected()
			}

			lastInsertId, err := res.LastInsertId()
			if err != nil {
				DebugLog.Println(ErrLastInsertIdUnavailable, ':', err)
				return lastInsertId, ErrLastInsertIdUnavailable
			} else {
				return lastInsertId, nil
			}
		}
		return 0, err
	}

	row := q.QueryRow(query, values...)
	var id int64
	err = row.Scan(&id)
	return id, err
}

// create new mysql dbBaser.
func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
