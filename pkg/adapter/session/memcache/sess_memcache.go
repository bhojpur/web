package memcache

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

// memcache for session provider
//
// depend on github.com/bradfitz/gomemcache/memcache
//
// go install github.com/bradfitz/gomemcache/memcache
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/memcache"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("memcache", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:11211"}``)
//		go globalSessions.GC()
//	}

import (
	"context"
	"net/http"

	"github.com/bhojpur/web/pkg/adapter/session"

	bhojpurmem "github.com/bhojpur/session/pkg/provider/memcache"
)

// SessionStore memcache session store
type SessionStore bhojpurmem.SessionStore

// Set value in memcache session
func (rs *SessionStore) Set(key, value interface{}) error {
	return (*bhojpurmem.SessionStore)(rs).Set(context.Background(), key, value)
}

// Get value in memcache session
func (rs *SessionStore) Get(key interface{}) interface{} {
	return (*bhojpurmem.SessionStore)(rs).Get(context.Background(), key)
}

// Delete value in memcache session
func (rs *SessionStore) Delete(key interface{}) error {
	return (*bhojpurmem.SessionStore)(rs).Delete(context.Background(), key)
}

// Flush clear all values in memcache session
func (rs *SessionStore) Flush() error {
	return (*bhojpurmem.SessionStore)(rs).Flush(context.Background())
}

// SessionID get memcache session id
func (rs *SessionStore) SessionID() string {
	return (*bhojpurmem.SessionStore)(rs).SessionID(context.Background())
}

// SessionRelease save session values to memcache
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*bhojpurmem.SessionStore)(rs).SessionRelease(context.Background(), w)
}

// MemProvider memcache session provider
type MemProvider bhojpurmem.MemProvider

// SessionInit init memcache session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemProvider) SessionInit(maxlifetime int64, savePath string) error {
	return (*bhojpurmem.MemProvider)(rp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read memcache session by sid
func (rp *MemProvider) SessionRead(sid string) (session.Store, error) {
	s, err := (*bhojpurmem.MemProvider)(rp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check memcache session exist by sid
func (rp *MemProvider) SessionExist(sid string) bool {
	res, _ := (*bhojpurmem.MemProvider)(rp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for memcache session
func (rp *MemProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*bhojpurmem.MemProvider)(rp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete memcache session by id
func (rp *MemProvider) SessionDestroy(sid string) error {
	return (*bhojpurmem.MemProvider)(rp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (rp *MemProvider) SessionGC() {
	(*bhojpurmem.MemProvider)(rp).SessionGC(context.Background())
}

// SessionAll return all activeSession
func (rp *MemProvider) SessionAll() int {
	return (*bhojpurmem.MemProvider)(rp).SessionAll(context.Background())
}
