//go:build go1.17
// +build go1.17

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

package orm

// ORM provides Object Relationship Mapping for MySQL/PostgreSQL/sqlite
// Simple Usage
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/bhojpur/web/pkg/orm"
//		_ "github.com/go-sql-driver/mysql" // import your used driver
//	)
//
//	// Model Struct
//	type User struct {
//		Id   int    `orm:"auto"`
//		Name string `orm:"size(100)"`
//	}
//
//	func init() {
//		orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
//	}
//
//	func main() {
//		o := orm.NewOrm()
//		user := User{Name: "pramila"}
//		// insert
//		id, err := o.Insert(&user)
//		// update
//		user.Name = "bhojpur"
//		num, err := o.Update(&user)
//		// read one
//		u := User{Id: user.Id}
//		err = o.Read(&u)
//		// delete
//		num, err = o.Delete(&u)
//	}

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bhojpur/web/pkg/client/orm"
	"github.com/bhojpur/web/pkg/client/orm/hints"
	"github.com/bhojpur/web/pkg/core/utils"
)

// DebugQueries define the debug
const (
	DebugQueries = iota
)

// Define common vars
var (
	Debug            = orm.Debug
	DebugLog         = orm.DebugLog
	DefaultRowsLimit = orm.DefaultRowsLimit
	DefaultRelsDepth = orm.DefaultRelsDepth
	DefaultTimeLoc   = orm.DefaultTimeLoc
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
)

type ormer struct {
	delegate   orm.Ormer
	txDelegate orm.TxOrmer
	isTx       bool
}

var _ Ormer = new(ormer)

// read data to model
func (o *ormer) Read(md interface{}, cols ...string) error {
	if o.isTx {
		return o.txDelegate.Read(md, cols...)
	}
	return o.delegate.Read(md, cols...)
}

// read data to model, like Read(), but use "SELECT FOR UPDATE" form
func (o *ormer) ReadForUpdate(md interface{}, cols ...string) error {
	if o.isTx {
		return o.txDelegate.ReadForUpdate(md, cols...)
	}
	return o.delegate.ReadForUpdate(md, cols...)
}

// Try to read a row from the database, or insert one if it doesn't exist
func (o *ormer) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	if o.isTx {
		return o.txDelegate.ReadOrCreate(md, col1, cols...)
	}
	return o.delegate.ReadOrCreate(md, col1, cols...)
}

// insert model data to database
func (o *ormer) Insert(md interface{}) (int64, error) {
	if o.isTx {
		return o.txDelegate.Insert(md)
	}
	return o.delegate.Insert(md)
}

// insert some models to database
func (o *ormer) InsertMulti(bulk int, mds interface{}) (int64, error) {
	if o.isTx {
		return o.txDelegate.InsertMulti(bulk, mds)
	}
	return o.delegate.InsertMulti(bulk, mds)
}

// InsertOrUpdate data to database
func (o *ormer) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	if o.isTx {
		return o.txDelegate.InsertOrUpdate(md, colConflitAndArgs...)
	}
	return o.delegate.InsertOrUpdate(md, colConflitAndArgs...)
}

// update model to database.
// cols set the columns those want to update.
func (o *ormer) Update(md interface{}, cols ...string) (int64, error) {
	if o.isTx {
		return o.txDelegate.Update(md, cols...)
	}
	return o.delegate.Update(md, cols...)
}

// delete model in database
// cols shows the delete conditions values read from. default is pk
func (o *ormer) Delete(md interface{}, cols ...string) (int64, error) {
	if o.isTx {
		return o.txDelegate.Delete(md, cols...)
	}
	return o.delegate.Delete(md, cols...)
}

// create a models to models queryer
func (o *ormer) QueryM2M(md interface{}, name string) QueryM2Mer {
	if o.isTx {
		return o.txDelegate.QueryM2M(md, name)
	}
	return o.delegate.QueryM2M(md, name)
}

// load related models to md model.
// args are limit, offset int and order string.
//
// example:
// 	orm.LoadRelated(post,"Tags")
// 	for _,tag := range post.Tags{...}
//
// make sure the relation is defined in model struct tags.
func (o *ormer) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	kvs := make([]utils.KV, 0, 4)
	for i, arg := range args {
		switch i {
		case 0:
			if v, ok := arg.(bool); ok {
				if v {
					kvs = append(kvs, hints.DefaultRelDepth())
				}
			} else if v, ok := arg.(int); ok {
				kvs = append(kvs, hints.RelDepth(v))
			}
		case 1:
			kvs = append(kvs, hints.Limit(orm.ToInt64(arg)))
		case 2:
			kvs = append(kvs, hints.Offset(orm.ToInt64(arg)))
		case 3:
			kvs = append(kvs, hints.Offset(orm.ToInt64(arg)))
		}
	}
	if o.isTx {
		return o.txDelegate.LoadRelated(md, name, kvs...)
	}
	return o.delegate.LoadRelated(md, name, kvs...)
}

// return a QuerySeter for table operations.
// table name can be string or struct.
// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
func (o *ormer) QueryTable(ptrStructOrTableName interface{}) (qs QuerySetter) {
	if o.isTx {
		return o.txDelegate.QueryTable(ptrStructOrTableName)
	}
	return o.delegate.QueryTable(ptrStructOrTableName)
}

// switch to another registered database driver by given name.
func (o *ormer) Using(name string) error {
	if o.isTx {
		return ErrTxHasBegan
	}
	o.delegate = orm.NewOrmUsingDB(name)
	return nil
}

// begin transaction
func (o *ormer) Begin() error {
	if o.isTx {
		return ErrTxHasBegan
	}
	return o.BeginTx(context.Background(), nil)
}

func (o *ormer) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	if o.isTx {
		return ErrTxHasBegan
	}
	txOrmer, err := o.delegate.BeginWithCtxAndOpts(ctx, opts)
	if err != nil {
		return err
	}
	o.txDelegate = txOrmer
	o.isTx = true
	return nil
}

// commit transaction
func (o *ormer) Commit() error {
	if !o.isTx {
		return ErrTxDone
	}
	err := o.txDelegate.Commit()
	if err == nil {
		o.isTx = false
		o.txDelegate = nil
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

// rollback transaction
func (o *ormer) Rollback() error {
	if !o.isTx {
		return ErrTxDone
	}
	err := o.txDelegate.Rollback()
	if err == nil {
		o.isTx = false
		o.txDelegate = nil
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

// return a raw query setter for raw sql string.
func (o *ormer) Raw(query string, args ...interface{}) RawSetter {
	if o.isTx {
		return o.txDelegate.Raw(query, args...)
	}
	return o.delegate.Raw(query, args...)
}

// return current using database Driver
func (o *ormer) Driver() Driver {
	if o.isTx {
		return o.txDelegate.Driver()
	}
	return o.delegate.Driver()
}

// return sql.DBStats for current database
func (o *ormer) DBStats() *sql.DBStats {
	if o.isTx {
		return o.txDelegate.DBStats()
	}
	return o.delegate.DBStats()
}

// NewOrm create new orm
func NewOrm() Ormer {
	o := orm.NewOrm()
	return &ormer{
		delegate: o,
	}
}

// NewOrmWithDB create a new ormer object with specify *sql.DB for query
func NewOrmWithDB(driverName, aliasName string, db *sql.DB) (Ormer, error) {
	o, err := orm.NewOrmWithDB(driverName, aliasName, db)
	if err != nil {
		return nil, err
	}
	return &ormer{
		delegate: o,
	}, nil
}
