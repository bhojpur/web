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
)

// an insert queryer struct
type insertSet struct {
	mi     *modelInfo
	orm    *ormBase
	stmt   stmtQuerier
	closed bool
}

var _ Inserter = new(insertSet)

// insert model ignore it's registered or not.
func (o *insertSet) Insert(md interface{}) (int64, error) {
	if o.closed {
		return 0, ErrStmtClosed
	}
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := getFullName(typ)
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Inserter.Insert> cannot use non-ptr model struct `%s`", name))
	}
	if name != o.mi.fullName {
		panic(fmt.Errorf("<Inserter.Insert> need model `%s` but found `%s`", o.mi.fullName, name))
	}
	id, err := o.orm.alias.DbBaser.InsertStmt(o.stmt, o.mi, ind, o.orm.alias.TZ)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if o.mi.fields.pk.auto {
			if o.mi.fields.pk.fieldType&IsPositiveIntegerField > 0 {
				ind.FieldByIndex(o.mi.fields.pk.fieldIndex).SetUint(uint64(id))
			} else {
				ind.FieldByIndex(o.mi.fields.pk.fieldIndex).SetInt(id)
			}
		}
	}
	return id, nil
}

// close insert queryer statement
func (o *insertSet) Close() error {
	if o.closed {
		return ErrStmtClosed
	}
	o.closed = true
	return o.stmt.Close()
}

// create new insert queryer.
func newInsertSet(orm *ormBase, mi *modelInfo) (Inserter, error) {
	bi := new(insertSet)
	bi.orm = orm
	bi.mi = mi
	st, query, err := orm.alias.DbBaser.PrepareInsert(orm.db, mi)
	if err != nil {
		return nil, err
	}
	if Debug {
		bi.stmt = newStmtQueryLog(orm.alias, st, query)
	} else {
		bi.stmt = st
	}
	return bi, nil
}
