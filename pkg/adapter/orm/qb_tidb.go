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
	"github.com/bhojpur/web/pkg/client/orm"
)

// TiDBQueryBuilder is the SQL build
type TiDBQueryBuilder orm.TiDBQueryBuilder

// Select will join the fields
func (qb *TiDBQueryBuilder) Select(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Select(fields...)
}

// ForUpdate add the FOR UPDATE clause
func (qb *TiDBQueryBuilder) ForUpdate() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).ForUpdate()
}

// From join the tables
func (qb *TiDBQueryBuilder) From(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).From(tables...)
}

// InnerJoin INNER JOIN the table
func (qb *TiDBQueryBuilder) InnerJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).InnerJoin(table)
}

// LeftJoin LEFT JOIN the table
func (qb *TiDBQueryBuilder) LeftJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).LeftJoin(table)
}

// RightJoin RIGHT JOIN the table
func (qb *TiDBQueryBuilder) RightJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).RightJoin(table)
}

// On join with on cond
func (qb *TiDBQueryBuilder) On(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).On(cond)
}

// Where join the Where cond
func (qb *TiDBQueryBuilder) Where(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Where(cond)
}

// And join the and cond
func (qb *TiDBQueryBuilder) And(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).And(cond)
}

// Or join the or cond
func (qb *TiDBQueryBuilder) Or(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Or(cond)
}

// In join the IN (vals)
func (qb *TiDBQueryBuilder) In(vals ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).In(vals...)
}

// OrderBy join the Order by fields
func (qb *TiDBQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).OrderBy(fields...)
}

// Asc join the asc
func (qb *TiDBQueryBuilder) Asc() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Asc()
}

// Desc join the desc
func (qb *TiDBQueryBuilder) Desc() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Desc()
}

// Limit join the limit num
func (qb *TiDBQueryBuilder) Limit(limit int) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Limit(limit)
}

// Offset join the offset num
func (qb *TiDBQueryBuilder) Offset(offset int) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Offset(offset)
}

// GroupBy join the Group by fields
func (qb *TiDBQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).GroupBy(fields...)
}

// Having join the Having cond
func (qb *TiDBQueryBuilder) Having(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Having(cond)
}

// Update join the update table
func (qb *TiDBQueryBuilder) Update(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Update(tables...)
}

// Set join the set kv
func (qb *TiDBQueryBuilder) Set(kv ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Set(kv...)
}

// Delete join the Delete tables
func (qb *TiDBQueryBuilder) Delete(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Delete(tables...)
}

// InsertInto join the insert SQL
func (qb *TiDBQueryBuilder) InsertInto(table string, fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).InsertInto(table, fields...)
}

// Values join the Values(vals)
func (qb *TiDBQueryBuilder) Values(vals ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Values(vals...)
}

// Subquery join the sub as alias
func (qb *TiDBQueryBuilder) Subquery(sub string, alias string) string {
	return (*orm.TiDBQueryBuilder)(qb).Subquery(sub, alias)
}

// String join all Tokens
func (qb *TiDBQueryBuilder) String() string {
	return (*orm.TiDBQueryBuilder)(qb).String()
}
