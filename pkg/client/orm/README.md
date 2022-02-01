# Bhojpur Web - ORM Framework

A powerful Object Relational Mapping framework for the [Bhojpur Web](https://github.com/bhojpur/web). It is heavily influenced by Django ORM, SQLAlchemy.

**Support Database:**

* MySQL: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
* PostgreSQL: [github.com/lib/pq](https://github.com/lib/pq)
* Sqlite3: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

**Features:**

* full Go language type support
* easy for usage, simple CRUD operation
* auto join with relation table
* cross Database compatible query
* Raw SQL query / mapper without ORM model
* full test keep stable and strong

**Installation:**

	go get github.com/bhojpur/web/pkg/client/orm

## Quick Start

#### Simple Usage

```go
package main

import (
	"fmt"
	"github.com/bhojpur/web/pkg/client/orm"
	_ "github.com/go-sql-driver/mysql" // import your database driver
)

// Model Struct
type User struct {
	Id   int    `orm:"auto"`
	Name string `orm:"size(100)"`
}

func init() {
	// register model
	orm.RegisterModel(new(User))

	// set default database
	orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
	
	// create table
	orm.RunSyncdb("default", false, true)	
}

func main() {
	o := orm.NewOrm()

	user := User{Name: "pramila"}

	// insert
	id, err := o.Insert(&user)

	// update
	user.Name = "bhojpur"
	num, err := o.Update(&user)

	// read one
	u := User{Id: user.Id}
	err = o.Read(&u)

	// delete
	num, err = o.Delete(&u)	
}
```

#### Next with Relation

```go
type Post struct {
	Id    int    `orm:"auto"`
	Title string `orm:"size(100)"`
	User  *User  `orm:"rel(fk)"`
}

var posts []*Post
qs := o.QueryTable("post")
num, err := qs.Filter("User__Name", "pramila").All(&posts)
```

#### Use Raw sql

If you don't like ORM，use Raw SQL to query / mapping without ORM setting

```go
var maps []Params
num, err := o.Raw("SELECT id FROM user WHERE name = ?", "pramila").Values(&maps)
if num > 0 {
	fmt.Println(maps[0]["id"])
}
```

#### Transactions

```go
o.Begin()
...
user := User{Name: "pramila"}
id, err := o.Insert(&user)
if err == nil {
	o.Commit()
} else {
	o.Rollback()
}

```

#### Debug Log Queries

In development env, you can simple use

```go
func main() {
	orm.Debug = true
...
```

enable log queries.

output include all queries, such as exec / prepare / transaction.

like this:

```go
[ORM] - 2018-03-26 13:18:16 - [Queries/default] - [    db.Exec /     0.4ms] - [INSERT INTO `user` (`name`) VALUES (?)] - `pramila`
...
```

## Documentation

more details and examples in docs and test

[documents](http://docs.bhojpur.net/mvc/model/overview.md)
