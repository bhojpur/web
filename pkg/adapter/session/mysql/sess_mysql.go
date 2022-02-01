package mysql

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

// MySQL for session provider
//
// depends on github.com/go-sql-driver/mysql:
//
// go install github.com/go-sql-driver/mysql
//
// mysql session support need create table as sql:
//	CREATE TABLE `session` (
//	`session_key` char(64) NOT NULL,
//	`session_data` blob,
//	`session_expiry` int(11) unsigned NOT NULL,
//	PRIMARY KEY (`session_key`)
//	) ENGINE=MyISAM DEFAULT CHARSET=utf8;
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/mysql"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("mysql", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]"}``)
//		go globalSessions.GC()
//	}

import (
	"context"
	"net/http"

	"github.com/bhojpur/session/pkg/provider/mysql"
	"github.com/bhojpur/web/pkg/adapter/session"

	// import MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	// TableName store the session in MySQL
	TableName = mysql.TableName
	mysqlpder = &Provider{}
)

// SessionStore mysql session store
type SessionStore mysql.SessionStore

// Set value in mysql session.
// it is temp value in map.
func (st *SessionStore) Set(key, value interface{}) error {
	return (*mysql.SessionStore)(st).Set(context.Background(), key, value)
}

// Get value from mysql session
func (st *SessionStore) Get(key interface{}) interface{} {
	return (*mysql.SessionStore)(st).Get(context.Background(), key)
}

// Delete value in mysql session
func (st *SessionStore) Delete(key interface{}) error {
	return (*mysql.SessionStore)(st).Delete(context.Background(), key)
}

// Flush clear all values in mysql session
func (st *SessionStore) Flush() error {
	return (*mysql.SessionStore)(st).Flush(context.Background())
}

// SessionID get session id of this mysql session store
func (st *SessionStore) SessionID() string {
	return (*mysql.SessionStore)(st).SessionID(context.Background())
}

// SessionRelease save mysql session values to database.
// must call this method to save values to database.
func (st *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*mysql.SessionStore)(st).SessionRelease(context.Background(), w)
}

// Provider mysql session provider
type Provider mysql.Provider

// SessionInit init mysql session.
// savepath is the connection string of mysql.
func (mp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*mysql.Provider)(mp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead get mysql session by sid
func (mp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*mysql.Provider)(mp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check mysql session exist
func (mp *Provider) SessionExist(sid string) bool {
	res, _ := (*mysql.Provider)(mp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for mysql session
func (mp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*mysql.Provider)(mp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete mysql session by sid
func (mp *Provider) SessionDestroy(sid string) error {
	return (*mysql.Provider)(mp).SessionDestroy(context.Background(), sid)
}

// SessionGC delete expired values in mysql session
func (mp *Provider) SessionGC() {
	(*mysql.Provider)(mp).SessionGC(context.Background())
}

// SessionAll count values in mysql session
func (mp *Provider) SessionAll() int {
	return (*mysql.Provider)(mp).SessionAll(context.Background())
}
