package redis

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

// Redis for session provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/redis"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
// 	func init() {
// 		globalSessions, _ = session.NewManager("redis", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:7070"}``)
// 		go globalSessions.GC()
// 	}

import (
	"context"
	"net/http"

	"github.com/bhojpur/web/pkg/adapter/session"

	bhojpurRedis "github.com/bhojpur/session/pkg/provider/redis"
)

// MaxPoolSize redis max pool size
var MaxPoolSize = bhojpurRedis.MaxPoolSize

// SessionStore redis session store
type SessionStore bhojpurRedis.SessionStore

// Set value in redis session
func (rs *SessionStore) Set(key, value interface{}) error {
	return (*bhojpurRedis.SessionStore)(rs).Set(context.Background(), key, value)
}

// Get value in redis session
func (rs *SessionStore) Get(key interface{}) interface{} {
	return (*bhojpurRedis.SessionStore)(rs).Get(context.Background(), key)
}

// Delete value in redis session
func (rs *SessionStore) Delete(key interface{}) error {
	return (*bhojpurRedis.SessionStore)(rs).Delete(context.Background(), key)
}

// Flush clear all values in redis session
func (rs *SessionStore) Flush() error {
	return (*bhojpurRedis.SessionStore)(rs).Flush(context.Background())
}

// SessionID get redis session id
func (rs *SessionStore) SessionID() string {
	return (*bhojpurRedis.SessionStore)(rs).SessionID(context.Background())
}

// SessionRelease save session values to redis
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*bhojpurRedis.SessionStore)(rs).SessionRelease(context.Background(), w)
}

// Provider redis session provider
type Provider bhojpurRedis.Provider

// SessionInit init redis session
// savepath like redis server addr,pool size,password,dbnum,IdleTimeout second
// e.g. 127.0.0.1:6379,100,bhojpur,0,30
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*bhojpurRedis.Provider)(rp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read redis session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*bhojpurRedis.Provider)(rp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check redis session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	res, _ := (*bhojpurRedis.Provider)(rp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for redis session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*bhojpurRedis.Provider)(rp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	return (*bhojpurRedis.Provider)(rp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
	(*bhojpurRedis.Provider)(rp).SessionGC(context.Background())
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return (*bhojpurRedis.Provider)(rp).SessionAll(context.Background())
}
