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
	"time"

	"github.com/bhojpur/web/pkg/client/orm"
)

// DriverType database driver constant int.
type DriverType orm.DriverType

// Enum the Database driver
const (
	DRMySQL    = DriverType(orm.DRMySQL)
	DRSqlite   = DriverType(orm.DRSqlite)   // sqlite
	DROracle   = DriverType(orm.DROracle)   // oracle
	DRPostgres = DriverType(orm.DRPostgres) // pgsql
	DRTiDB     = DriverType(orm.DRTiDB)     // TiDB
)

type DB orm.DB

func (d *DB) Begin() (*sql.Tx, error) {
	return (*orm.DB)(d).Begin()
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return (*orm.DB)(d).BeginTx(ctx, opts)
}

func (d *DB) Prepare(query string) (*sql.Stmt, error) {
	return (*orm.DB)(d).Prepare(query)
}

func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return (*orm.DB)(d).PrepareContext(ctx, query)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return (*orm.DB)(d).Exec(query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return (*orm.DB)(d).ExecContext(ctx, query, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return (*orm.DB)(d).Query(query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return (*orm.DB)(d).QueryContext(ctx, query, args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return (*orm.DB)(d).QueryRow(query, args)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return (*orm.DB)(d).QueryRowContext(ctx, query, args...)
}

// AddAliasWthDB add a aliasName for the drivename
func AddAliasWthDB(aliasName, driverName string, db *sql.DB) error {
	return orm.AddAliasWthDB(aliasName, driverName, db)
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) error {
	opts := make([]orm.DBOption, 0, 2)
	if len(params) > 0 {
		opts = append(opts, orm.MaxIdleConnections(params[0]))
	}

	if len(params) > 1 {
		opts = append(opts, orm.MaxOpenConnections(params[1]))
	}
	return orm.RegisterDataBase(aliasName, driverName, dataSource, opts...)
}

// RegisterDriver Register a database driver use specify driver name, this can be definition the driver is which database type.
func RegisterDriver(driverName string, typ DriverType) error {
	return orm.RegisterDriver(driverName, orm.DriverType(typ))
}

// SetDataBaseTZ Change the database default used timezone
func SetDataBaseTZ(aliasName string, tz *time.Location) error {
	return orm.SetDataBaseTZ(aliasName, tz)
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
func SetMaxIdleConns(aliasName string, maxIdleConns int) {
	orm.SetMaxIdleConns(aliasName, maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
func SetMaxOpenConns(aliasName string, maxOpenConns int) {
	orm.SetMaxOpenConns(aliasName, maxOpenConns)
}

// GetDB Get *sql.DB from registered database by db alias name.
// Use "default" as alias name if you not set.
func GetDB(aliasNames ...string) (*sql.DB, error) {
	return orm.GetDB(aliasNames...)
}
