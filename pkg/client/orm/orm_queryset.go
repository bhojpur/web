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
	"context"
	"fmt"

	"github.com/bhojpur/web/pkg/client/orm/hints"
)

type colValue struct {
	value int64
	opt   operator
}

type operator int

// define Col operations
const (
	ColAdd operator = iota
	ColMinus
	ColMultiply
	ColExcept
	ColBitAnd
	ColBitRShift
	ColBitLShift
	ColBitXOR
	ColBitOr
)

// ColValue do the field raw changes. e.g Nums = Nums + 10. usage:
// 	Params{
// 		"Nums": ColValue(Col_Add, 10),
// 	}
func ColValue(opt operator, value interface{}) interface{} {
	switch opt {
	case ColAdd, ColMinus, ColMultiply, ColExcept, ColBitAnd, ColBitRShift,
		ColBitLShift, ColBitXOR, ColBitOr:
	default:
		panic(fmt.Errorf("orm.ColValue wrong operator"))
	}
	v, err := StrTo(ToStr(value)).Int64()
	if err != nil {
		panic(fmt.Errorf("orm.ColValue doesn't support non string/numeric type, %s", err))
	}
	var val colValue
	val.value = v
	val.opt = opt
	return val
}

// real query struct
type querySet struct {
	mi         *modelInfo
	cond       *Condition
	related    []string
	relDepth   int
	limit      int64
	offset     int64
	groups     []string
	orders     []string
	distinct   bool
	forUpdate  bool
	useIndex   int
	indexes    []string
	orm        *ormBase
	ctx        context.Context
	forContext bool
}

var _ QuerySetter = new(querySet)

// add condition expression to QuerySeter.
func (o querySet) Filter(expr string, args ...interface{}) QuerySetter {
	if o.cond == nil {
		o.cond = NewCondition()
	}
	o.cond = o.cond.And(expr, args...)
	return &o
}

// add raw sql to querySeter.
func (o querySet) FilterRaw(expr string, sql string) QuerySetter {
	if o.cond == nil {
		o.cond = NewCondition()
	}
	o.cond = o.cond.Raw(expr, sql)
	return &o
}

// add NOT condition to querySeter.
func (o querySet) Exclude(expr string, args ...interface{}) QuerySetter {
	if o.cond == nil {
		o.cond = NewCondition()
	}
	o.cond = o.cond.AndNot(expr, args...)
	return &o
}

// set offset number
func (o *querySet) setOffset(num interface{}) {
	o.offset = ToInt64(num)
}

// add LIMIT value.
// args[0] means offset, e.g. LIMIT num,offset.
func (o querySet) Limit(limit interface{}, args ...interface{}) QuerySetter {
	o.limit = ToInt64(limit)
	if len(args) > 0 {
		o.setOffset(args[0])
	}
	return &o
}

// add OFFSET value
func (o querySet) Offset(offset interface{}) QuerySetter {
	o.setOffset(offset)
	return &o
}

// add GROUP expression
func (o querySet) GroupBy(exprs ...string) QuerySetter {
	o.groups = exprs
	return &o
}

// add ORDER expression.
// "column" means ASC, "-column" means DESC.
func (o querySet) OrderBy(exprs ...string) QuerySetter {
	o.orders = exprs
	return &o
}

// add DISTINCT to SELECT
func (o querySet) Distinct() QuerySetter {
	o.distinct = true
	return &o
}

// add FOR UPDATE to SELECT
func (o querySet) ForUpdate() QuerySetter {
	o.forUpdate = true
	return &o
}

// ForceIndex force index for query
func (o querySet) ForceIndex(indexes ...string) QuerySetter {
	o.useIndex = hints.KeyForceIndex
	o.indexes = indexes
	return &o
}

// UseIndex use index for query
func (o querySet) UseIndex(indexes ...string) QuerySetter {
	o.useIndex = hints.KeyUseIndex
	o.indexes = indexes
	return &o
}

// IgnoreIndex ignore index for query
func (o querySet) IgnoreIndex(indexes ...string) QuerySetter {
	o.useIndex = hints.KeyIgnoreIndex
	o.indexes = indexes
	return &o
}

// set relation model to query together.
// it will query relation models and assign to parent model.
func (o querySet) RelatedSel(params ...interface{}) QuerySetter {
	if len(params) == 0 {
		o.relDepth = DefaultRelsDepth
	} else {
		for _, p := range params {
			switch val := p.(type) {
			case string:
				o.related = append(o.related, val)
			case int:
				o.relDepth = val
			default:
				panic(fmt.Errorf("<QuerySetter.RelatedSel> wrong param kind: %v", val))
			}
		}
	}
	return &o
}

// set condition to QuerySeter.
func (o querySet) SetCond(cond *Condition) QuerySetter {
	o.cond = cond
	return &o
}

// get condition from QuerySeter
func (o querySet) GetCond() *Condition {
	return o.cond
}

// return QuerySeter execution result number
func (o *querySet) Count() (int64, error) {
	return o.orm.alias.DbBaser.Count(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
}

// check result empty or not after QuerySeter executed
func (o *querySet) Exist() bool {
	cnt, _ := o.orm.alias.DbBaser.Count(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
	return cnt > 0
}

// execute update with parameters
func (o *querySet) Update(values Params) (int64, error) {
	return o.orm.alias.DbBaser.UpdateBatch(o.orm.db, o, o.mi, o.cond, values, o.orm.alias.TZ)
}

// execute delete
func (o *querySet) Delete() (int64, error) {
	return o.orm.alias.DbBaser.DeleteBatch(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
}

// return a insert queryer.
// it can be used in times.
// example:
// 	i,err := sq.PrepareInsert()
// 	i.Add(&user1{},&user2{})
func (o *querySet) PrepareInsert() (Inserter, error) {
	return newInsertSet(o.orm, o.mi)
}

// query all data and map to containers.
// cols means the columns when querying.
func (o *querySet) All(container interface{}, cols ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ, cols)
}

// query one row data and map to containers.
// cols means the columns when querying.
func (o *querySet) One(container interface{}, cols ...string) error {
	o.limit = 1
	num, err := o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ, cols)
	if err != nil {
		return err
	}
	if num == 0 {
		return ErrNoRows
	}

	if num > 1 {
		return ErrMultiRows
	}
	return nil
}

// query all data and map to []map[string]interface.
// expres means condition expression.
// it converts data to []map[column]value.
func (o *querySet) Values(results *[]Params, exprs ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results, o.orm.alias.TZ)
}

// query all data and map to [][]interface
// it converts data to [][column_index]value
func (o *querySet) ValuesList(results *[]ParamsList, exprs ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results, o.orm.alias.TZ)
}

// query all data and map to []interface.
// it's designed for one row record set, auto change to []value, not [][column]value.
func (o *querySet) ValuesFlat(result *ParamsList, expr string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, []string{expr}, result, o.orm.alias.TZ)
}

// query all rows into map[string]interface with specify key and value column name.
// keyCol = "name", valueCol = "value"
// table data
// name  | value
// total | 100
// found | 200
// to map[string]interface{}{
// 	"total": 100,
// 	"found": 200,
// }
func (o *querySet) RowsToMap(result *Params, keyCol, valueCol string) (int64, error) {
	panic(ErrNotImplement)
}

// query all rows into struct with specify key and value column name.
// keyCol = "name", valueCol = "value"
// table data
// name  | value
// total | 100
// found | 200
// to struct {
// 	Total int
// 	Found int
// }
func (o *querySet) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	panic(ErrNotImplement)
}

// set context to QuerySeter.
func (o querySet) WithContext(ctx context.Context) QuerySetter {
	o.ctx = ctx
	o.forContext = true
	return &o
}

// create new QuerySeter.
func newQuerySet(orm *ormBase, mi *modelInfo) QuerySetter {
	o := new(querySet)
	o.mi = mi
	o.orm = orm
	return o
}
