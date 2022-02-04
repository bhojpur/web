package migration

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

// The table structure is as follow:
//
//	CREATE TABLE `migrations` (
//		`id_migration` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
//		`name` varchar(255) DEFAULT NULL COMMENT 'migration name, unique',
//		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
//		`statements` longtext COMMENT 'SQL statements for this migration',
//		`rollback_statements` longtext,
//		`status` enum('update','rollback') DEFAULT NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
//		PRIMARY KEY (`id_migration`)
//	) ENGINE=InnoDB DEFAULT CHARSET=utf8;

import (
	"github.com/bhojpur/web/pkg/client/orm/migration"
)

// const the data format for the Bhojpur.NET Platform generate migration datatype
const (
	DateFormat   = "20060102_150405"
	DBDateFormat = "2006-01-02 15:04:05"
)

// Migrationer is an interface for all Migration struct
type Migrationer interface {
	Up()
	Down()
	Reset()
	Exec(name, status string) error
	GetCreated() int64
}

// Migration defines the migrations by either SQL or DDL
type Migration migration.Migration

// Up implement in the Inheritance struct for upgrade
func (m *Migration) Up() {
	(*migration.Migration)(m).Up()
}

// Down implement in the Inheritance struct for down
func (m *Migration) Down() {
	(*migration.Migration)(m).Down()
}

// Migrate adds the SQL to the execution list
func (m *Migration) Migrate(migrationType string) {
	(*migration.Migration)(m).Migrate(migrationType)
}

// SQL add sql want to execute
func (m *Migration) SQL(sql string) {
	(*migration.Migration)(m).SQL(sql)
}

// Reset the sqls
func (m *Migration) Reset() {
	(*migration.Migration)(m).Reset()
}

// Exec execute the sql already add in the sql
func (m *Migration) Exec(name, status string) error {
	return (*migration.Migration)(m).Exec(name, status)
}

// GetCreated get the unixtime from the Created
func (m *Migration) GetCreated() int64 {
	return (*migration.Migration)(m).GetCreated()
}

// Register register the Migration in the map
func Register(name string, m Migrationer) error {
	return migration.Register(name, m)
}

// Upgrade upgrade the migration from lasttime
func Upgrade(lasttime int64) error {
	return migration.Upgrade(lasttime)
}

// Rollback rollback the migration by the name
func Rollback(name string) error {
	return migration.Rollback(name)
}

// Reset reset all migration
// run all migration's down function
func Reset() error {
	return migration.Reset()
}

// Refresh first Reset, then Upgrade
func Refresh() error {
	return migration.Refresh()
}
