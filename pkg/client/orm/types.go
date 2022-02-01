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
	"database/sql"
	"reflect"
	"time"

	"github.com/bhojpur/web/pkg/core/utils"
)

// TableNaming is usually used by model
// when you custom your table name, please implement this interfaces
// for example:
// type User struct {
//   ...
// }
// func (u *User) TableName() string {
//    return "USER_TABLE"
// }
type TableNameI interface {
	TableName() string
}

// TableEngineI is usually used by model
// when you want to use specific engine, like myisam, you can implement this interface
// for example:
// type User struct {
//   ...
// }
// func (u *User) TableEngine() string {
//    return "myisam"
// }
type TableEngineI interface {
	TableEngine() string
}

// TableIndexI is usually used by model
// when you want to create indexes, you can implement this interface
// for example:
// type User struct {
//   ...
// }
// func (u *User) TableIndex() [][]string {
//    return [][]string{{"Name"}}
// }
type TableIndexI interface {
	TableIndex() [][]string
}

// TableUniqueI is usually used by model
// when you want to create unique indexes, you can implement this interface
// for example:
// type User struct {
//   ...
// }
// func (u *User) TableUnique() [][]string {
//    return [][]string{{"Email"}}
// }
type TableUniqueI interface {
	TableUnique() [][]string
}

// IsApplicableTableForDB if return false, we won't create table to this db
type IsApplicableTableForDB interface {
	IsApplicableTableForDB(db string) bool
}

// Driver define database driver
type Driver interface {
	Name() string
	Type() DriverType
}

// Fielder define field info
type Fielder interface {
	String() string
	FieldType() int
	SetRaw(interface{}) error
	RawValue() interface{}
}

type TxBeginner interface {
	//self control transaction
	Begin() (TxOrmer, error)
	BeginWithCtx(ctx context.Context) (TxOrmer, error)
	BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error)
	BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error)

	//closure control transaction
	DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
}

type TxCommitter interface {
	Commit() error
	Rollback() error
}

//Data Manipulation Language
type DML interface {
	// insert model data to database
	// for example:
	//  user := new(User)
	//  id, err = Ormer.Insert(user)
	//  user must be a pointer and Insert will set user's pk field
	Insert(md interface{}) (int64, error)
	InsertWithCtx(ctx context.Context, md interface{}) (int64, error)
	// mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
	// if colu type is integer : can use(+-*/), string : convert(colu,"value")
	// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
	// if colu type is integer : can use(+-*/), string : colu || "value"
	InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
	InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error)
	// insert some models to database
	InsertMulti(bulk int, mds interface{}) (int64, error)
	InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error)
	// update model to database.
	// cols set the columns those want to update.
	// find model by Id(pk) field and update columns specified by fields, if cols is null then update all columns
	// for example:
	// user := User{Id: 2}
	//	user.Langs = append(user.Langs, "zh-CN", "en-US")
	//	user.Extra.Name = "bhojpur"
	//	user.Extra.Data = "orm"
	//	num, err = Ormer.Update(&user, "Langs", "Extra")
	Update(md interface{}, cols ...string) (int64, error)
	UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)
	// delete model in database
	Delete(md interface{}, cols ...string) (int64, error)
	DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)

	// return a raw query seter for raw sql string.
	// for example:
	//	 ormer.Raw("UPDATE `user` SET `user_name` = ? WHERE `user_name` = ?", "pramila", "testing").Exec()
	//	// update user testing's name to pramila
	Raw(query string, args ...interface{}) RawSetter
	RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSetter
}

// Data Query Language
type DQL interface {
	// read data to model
	// for example:
	//	this will find User by Id field
	// 	u = &User{Id: user.Id}
	// 	err = Ormer.Read(u)
	//	this will find User by UserName field
	// 	u = &User{UserName: "bhojpur", Password: "pass"}
	//	err = Ormer.Read(u, "UserName")
	Read(md interface{}, cols ...string) error
	ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error

	// Like Read(), but with "FOR UPDATE" clause, useful in transaction.
	// Some databases are not support this feature.
	ReadForUpdate(md interface{}, cols ...string) error
	ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error

	// Try to read a row from the database, or insert one if it doesn't exist
	ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error)
	ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error)

	// load related models to md model.
	// args are limit, offset int and order string.
	//
	// example:
	// 	Ormer.LoadRelated(post,"Tags")
	// 	for _,tag := range post.Tags{...}
	// hints.DefaultRelDepth useDefaultRelsDepth ; or depth 0
	// hints.RelDepth loadRelationDepth
	// hints.Limit limit default limit 1000
	// hints.Offset int offset default offset 0
	// hints.OrderBy string order  for example : "-Id"
	// make sure the relation is defined in model struct tags.
	LoadRelated(md interface{}, name string, args ...utils.KV) (int64, error)
	LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils.KV) (int64, error)

	// create a models to models queryer
	// for example:
	// 	post := Post{Id: 4}
	// 	m2m := Ormer.QueryM2M(&post, "Tags")
	QueryM2M(md interface{}, name string) QueryM2Mer
	QueryM2MWithCtx(ctx context.Context, md interface{}, name string) QueryM2Mer

	// return a QuerySetter for table operations.
	// table name can be string or struct.
	// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
	QueryTable(ptrStructOrTableName interface{}) QuerySetter
	QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) QuerySetter

	DBStats() *sql.DBStats
}

type DriverGetter interface {
	Driver() Driver
}

type ormer interface {
	DQL
	DML
	DriverGetter
}

type Ormer interface {
	ormer
	TxBeginner
}

type TxOrmer interface {
	ormer
	TxCommitter
}

// Inserter insert prepared statement
type Inserter interface {
	Insert(interface{}) (int64, error)
	Close() error
}

// QuerySetter query seter
type QuerySetter interface {
	// add condition expression to QuerySeter.
	// for example:
	//	filter by UserName == 'pramila'
	//	qs.Filter("UserName", "pramila")
	//	sql : left outer join profile on t0.id1==t1.id2 where t1.age == 28
	//	Filter("profile__Age", 28)
	// 	 // time compare
	//	qs.Filter("created", time.Now())
	Filter(string, ...interface{}) QuerySetter
	// add raw sql to querySeter.
	// for example:
	// qs.FilterRaw("user_id IN (SELECT id FROM profile WHERE age>=18)")
	// //sql-> WHERE user_id IN (SELECT id FROM profile WHERE age>=18)
	FilterRaw(string, string) QuerySetter
	// add NOT condition to querySeter.
	// have the same usage as Filter
	Exclude(string, ...interface{}) QuerySetter
	// set condition to QuerySeter.
	// sql's where condition
	//	cond := orm.NewCondition()
	//	cond1 := cond.And("profile__isnull", false).AndNot("status__in", 1).Or("profile__age__gt", 2000)
	//	//sql-> WHERE T0.`profile_id` IS NOT NULL AND NOT T0.`Status` IN (?) OR T1.`age` >  2000
	//	num, err := qs.SetCond(cond1).Count()
	SetCond(*Condition) QuerySetter
	// get condition from QuerySeter.
	// sql's where condition
	//  cond := orm.NewCondition()
	//  cond = cond.And("profile__isnull", false).AndNot("status__in", 1)
	//  qs = qs.SetCond(cond)
	//  cond = qs.GetCond()
	//  cond := cond.Or("profile__age__gt", 2000)
	//  //sql-> WHERE T0.`profile_id` IS NOT NULL AND NOT T0.`Status` IN (?) OR T1.`age` >  2000
	//  num, err := qs.SetCond(cond).Count()
	GetCond() *Condition
	// add LIMIT value.
	// args[0] means offset, e.g. LIMIT num,offset.
	// if Limit <= 0 then Limit will be set to default limit ,eg 1000
	// if QuerySeter doesn't call Limit, the sql's Limit will be set to default limit, eg 1000
	//  for example:
	//	qs.Limit(10, 2)
	//	// sql-> limit 10 offset 2
	Limit(limit interface{}, args ...interface{}) QuerySetter
	// add OFFSET value
	// same as Limit function's args[0]
	Offset(offset interface{}) QuerySetter
	// add GROUP BY expression
	// for example:
	//	qs.GroupBy("id")
	GroupBy(exprs ...string) QuerySetter
	// add ORDER expression.
	// "column" means ASC, "-column" means DESC.
	// for example:
	//	qs.OrderBy("-status")
	OrderBy(exprs ...string) QuerySetter
	// add FORCE INDEX expression.
	// for example:
	//	qs.ForceIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	ForceIndex(indexes ...string) QuerySetter
	// add USE INDEX expression.
	// for example:
	//	qs.UseIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	UseIndex(indexes ...string) QuerySetter
	// add IGNORE INDEX expression.
	// for example:
	//	qs.IgnoreIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	IgnoreIndex(indexes ...string) QuerySetter
	// set relation model to query together.
	// it will query relation models and assign to parent model.
	// for example:
	//	// will load all related fields use left join .
	// 	qs.RelatedSel().One(&user)
	//	// will  load related field only profile
	//	qs.RelatedSel("profile").One(&user)
	//	user.Profile.Age = 32
	RelatedSel(params ...interface{}) QuerySetter
	// Set Distinct
	// for example:
	//  o.QueryTable("policy").Filter("Groups__Group__Users__User", user).
	//    Distinct().
	//    All(&permissions)
	Distinct() QuerySetter
	// set FOR UPDATE to query.
	// for example:
	//  o.QueryTable("user").Filter("uid", uid).ForUpdate().All(&users)
	ForUpdate() QuerySetter
	// return QuerySeter execution result number
	// for example:
	//	num, err = qs.Filter("profile__age__gt", 28).Count()
	Count() (int64, error)
	// check result empty or not after QuerySeter executed
	// the same as QuerySeter.Count > 0
	Exist() bool
	// execute update with parameters
	// for example:
	//	num, err = qs.Filter("user_name", "pramila").Update(Params{
	//		"Nums": ColValue(Col_Minus, 50),
	//	}) // user pramila's Nums will minus 50
	//	num, err = qs.Filter("UserName", "pramila").Update(Params{
	//		"user_name": "pramila2"
	//	}) // user pramila's  name will change to pramila2
	Update(values Params) (int64, error)
	// delete from table
	// for example:
	//	num ,err = qs.Filter("user_name__in", "testing1", "testing2").Delete()
	// 	//delete two user  who's name is testing1 or testing2
	Delete() (int64, error)
	// return a insert queryer.
	// it can be used in times.
	// example:
	// 	i,err := sq.PrepareInsert()
	// 	num, err = i.Insert(&user1) // user table will add one record user1 at once
	//	num, err = i.Insert(&user2) // user table will add one record user2 at once
	//	err = i.Close() //don't forget call Close
	PrepareInsert() (Inserter, error)
	// query all data and map to containers.
	// cols means the columns when querying.
	// for example:
	//	var users []*User
	//	qs.All(&users) // users[0],users[1],users[2] ...
	All(container interface{}, cols ...string) (int64, error)
	// query one row data and map to containers.
	// cols means the columns when querying.
	// for example:
	//	var user User
	//	qs.One(&user) //user.UserName == "pramila"
	One(container interface{}, cols ...string) error
	// query all data and map to []map[string]interface.
	// expres means condition expression.
	// it converts data to []map[column]value.
	// for example:
	//	var maps []Params
	//	qs.Values(&maps) //maps[0]["UserName"]=="pramila"
	Values(results *[]Params, exprs ...string) (int64, error)
	// query all data and map to [][]interface
	// it converts data to [][column_index]value
	// for example:
	//	var list []ParamsList
	//	qs.ValuesList(&list) // list[0][1] == "pramila"
	ValuesList(results *[]ParamsList, exprs ...string) (int64, error)
	// query all data and map to []interface.
	// it's designed for one column record set, auto change to []value, not [][column]value.
	// for example:
	//	var list ParamsList
	//	qs.ValuesFlat(&list, "UserName") // list[0] == "pramila"
	ValuesFlat(result *ParamsList, expr string) (int64, error)
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
	RowsToMap(result *Params, keyCol, valueCol string) (int64, error)
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
	RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)
}

// QueryM2Mer model to model query struct
// all operations are on the m2m table only, will not affect the origin model table
type QueryM2Mer interface {
	// add models to origin models when creating queryM2M.
	// example:
	// 	m2m := orm.QueryM2M(post,"Tag")
	// 	m2m.Add(&Tag1{},&Tag2{})
	//  	for _,tag := range post.Tags{}{ ... }
	// param could also be any of the follow
	// 	[]*Tag{{Id:3,Name: "TestTag1"}, {Id:4,Name: "TestTag2"}}
	//	&Tag{Id:5,Name: "TestTag3"}
	//	[]interface{}{&Tag{Id:6,Name: "TestTag4"}}
	// insert one or more rows to m2m table
	// make sure the relation is defined in post model struct tag.
	Add(...interface{}) (int64, error)
	// remove models following the origin model relationship
	// only delete rows from m2m table
	// for example:
	// tag3 := &Tag{Id:5,Name: "TestTag3"}
	// num, err = m2m.Remove(tag3)
	Remove(...interface{}) (int64, error)
	// check model is existed in relationship of origin model
	Exist(interface{}) bool
	// clean all models in related of origin model
	Clear() (int64, error)
	// count all related models of origin model
	Count() (int64, error)
}

// RawPreparer raw query statement
type RawPreparer interface {
	Exec(...interface{}) (sql.Result, error)
	Close() error
}

// RawSetter raw query seter
// create From Ormer.Raw
// for example:
//  sql := fmt.Sprintf("SELECT %sid%s,%sname%s FROM %suser%s WHERE id = ?",Q,Q,Q,Q,Q,Q)
//  rs := Ormer.Raw(sql, 1)
type RawSetter interface {
	// execute sql and get result
	Exec() (sql.Result, error)
	// query data and map to container
	// for example:
	//	var name string
	//	var id int
	//	rs.QueryRow(&id,&name) // id==2 name=="pramila"
	QueryRow(containers ...interface{}) error

	// query data rows and map to container
	//	var ids []int
	//	var names []int
	//	query = fmt.Sprintf("SELECT 'id','name' FROM %suser%s", Q, Q)
	//	num, err = dORM.Raw(query).QueryRows(&ids,&names) // ids=>{1,2},names=>{"nobody","pramila"}
	QueryRows(containers ...interface{}) (int64, error)
	SetArgs(...interface{}) RawSetter
	// query data to []map[string]interface
	// see QuerySeter's Values
	Values(container *[]Params, cols ...string) (int64, error)
	// query data to [][]interface
	// see QuerySeter's ValuesList
	ValuesList(container *[]ParamsList, cols ...string) (int64, error)
	// query data to []interface
	// see QuerySeter's ValuesFlat
	ValuesFlat(container *ParamsList, cols ...string) (int64, error)
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
	RowsToMap(result *Params, keyCol, valueCol string) (int64, error)
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
	RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)

	// return prepared raw statement for used in times.
	// for example:
	// 	pre, err := dORM.Raw("INSERT INTO tag (name) VALUES (?)").Prepare()
	// 	r, err := pre.Exec("name1") // INSERT INTO tag (name) VALUES (`name1`)
	Prepare() (RawPreparer, error)
}

// stmtQuerier statement querier
type stmtQuerier interface {
	Close() error
	Exec(args ...interface{}) (sql.Result, error)
	// ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	// QueryContext(args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
	// QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
}

// db querier
type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// type DB interface {
// 	Begin() (*sql.Tx, error)
// 	Prepare(query string) (stmtQuerier, error)
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	Query(query string, args ...interface{}) (*sql.Rows, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// }

// transaction beginner
type txer interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// transaction ending
type txEnder interface {
	Commit() error
	Rollback() error
}

// base database struct
type dbBaser interface {
	Read(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string, bool) error
	ReadBatch(dbQuerier, *querySet, *modelInfo, *Condition, interface{}, *time.Location, []string) (int64, error)
	Count(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)
	ReadValues(dbQuerier, *querySet, *modelInfo, *Condition, []string, interface{}, *time.Location) (int64, error)

	Insert(dbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	InsertOrUpdate(dbQuerier, *modelInfo, reflect.Value, *alias, ...string) (int64, error)
	InsertMulti(dbQuerier, *modelInfo, reflect.Value, int, *time.Location) (int64, error)
	InsertValue(dbQuerier, *modelInfo, bool, []string, []interface{}) (int64, error)
	InsertStmt(stmtQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)

	Update(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) (int64, error)
	UpdateBatch(dbQuerier, *querySet, *modelInfo, *Condition, Params, *time.Location) (int64, error)

	Delete(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) (int64, error)
	DeleteBatch(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)

	SupportUpdateJoin() bool
	OperatorSQL(string) string
	GenerateOperatorSQL(*modelInfo, *fieldInfo, string, []interface{}, *time.Location) (string, []interface{})
	GenerateOperatorLeftCol(*fieldInfo, string, *string)
	PrepareInsert(dbQuerier, *modelInfo) (stmtQuerier, string, error)
	MaxLimit() uint64
	TableQuote() string
	ReplaceMarks(*string)
	HasReturningID(*modelInfo, *string) bool
	TimeFromDB(*time.Time, *time.Location)
	TimeToDB(*time.Time, *time.Location)
	DbTypes() map[string]string
	GetTables(dbQuerier) (map[string]bool, error)
	GetColumns(dbQuerier, string) (map[string][3]string, error)
	ShowTablesQuery() string
	ShowColumnsQuery(string) string
	IndexExists(dbQuerier, string, string) bool
	collectFieldValue(*modelInfo, *fieldInfo, reflect.Value, bool, *time.Location) (interface{}, error)
	setval(dbQuerier, *modelInfo, []string) error

	GenerateSpecifyIndex(tableName string, useIndex int, indexes []string) string
}
