package redis_sentinel

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

// Redis Sentinel for session provider
//
// depend on github.com/go-redis/redis
//
// go install github.com/go-redis/redis
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/redis_sentinel"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("redis_sentinel", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:26379;127.0.0.2:26379"}``)
//		go globalSessions.GC()
//	}
//
// more detail about params: please check the notes on the function SessionInit in this package

import (
	"context"
	"net/http"

	"github.com/bhojpur/web/pkg/adapter/session"

	sentinel "github.com/bhojpur/session/pkg/provider/redis_sentinel"
)

// DefaultPoolSize redis_sentinel default pool size
var DefaultPoolSize = sentinel.DefaultPoolSize

// SessionStore redis_sentinel session store
type SessionStore sentinel.SessionStore

// Set value in redis_sentinel session
func (rs *SessionStore) Set(key, value interface{}) error {
	return (*sentinel.SessionStore)(rs).Set(context.Background(), key, value)
}

// Get value in redis_sentinel session
func (rs *SessionStore) Get(key interface{}) interface{} {
	return (*sentinel.SessionStore)(rs).Get(context.Background(), key)
}

// Delete value in redis_sentinel session
func (rs *SessionStore) Delete(key interface{}) error {
	return (*sentinel.SessionStore)(rs).Delete(context.Background(), key)
}

// Flush clear all values in redis_sentinel session
func (rs *SessionStore) Flush() error {
	return (*sentinel.SessionStore)(rs).Flush(context.Background())
}

// SessionID get redis_sentinel session id
func (rs *SessionStore) SessionID() string {
	return (*sentinel.SessionStore)(rs).SessionID(context.Background())
}

// SessionRelease save session values to redis_sentinel
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*sentinel.SessionStore)(rs).SessionRelease(context.Background(), w)
}

// Provider redis_sentinel session provider
type Provider sentinel.Provider

// SessionInit init redis_sentinel session
// savepath like redis sentinel addr,pool size,password,dbnum,masterName
// e.g. 127.0.0.1:26379;127.0.0.2:26379,100,1qaz2wsx,0,mymaster
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*sentinel.Provider)(rp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read redis_sentinel session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*sentinel.Provider)(rp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check redis_sentinel session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	res, _ := (*sentinel.Provider)(rp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for redis_sentinel session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*sentinel.Provider)(rp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	return (*sentinel.Provider)(rp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
	(*sentinel.Provider)(rp).SessionGC(context.Background())
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return (*sentinel.Provider)(rp).SessionAll(context.Background())
}
