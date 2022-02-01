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

import "errors"

// QueryBuilder is the Query builder interface
type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	ForUpdate() QueryBuilder
	From(tables ...string) QueryBuilder
	InnerJoin(table string) QueryBuilder
	LeftJoin(table string) QueryBuilder
	RightJoin(table string) QueryBuilder
	On(cond string) QueryBuilder
	Where(cond string) QueryBuilder
	And(cond string) QueryBuilder
	Or(cond string) QueryBuilder
	In(vals ...string) QueryBuilder
	OrderBy(fields ...string) QueryBuilder
	Asc() QueryBuilder
	Desc() QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Having(cond string) QueryBuilder
	Update(tables ...string) QueryBuilder
	Set(kv ...string) QueryBuilder
	Delete(tables ...string) QueryBuilder
	InsertInto(table string, fields ...string) QueryBuilder
	Values(vals ...string) QueryBuilder
	Subquery(sub string, alias string) string
	String() string
}

// NewQueryBuilder return the QueryBuilder
func NewQueryBuilder(driver string) (qb QueryBuilder, err error) {
	if driver == "mysql" {
		qb = new(MySQLQueryBuilder)
	} else if driver == "tidb" {
		qb = new(TiDBQueryBuilder)
	} else if driver == "postgres" {
		qb = new(PostgresQueryBuilder)
	} else if driver == "sqlite" {
		err = errors.New("sqlite query builder is not supported yet")
	} else {
		err = errors.New("unknown driver for query builder")
	}
	return
}
