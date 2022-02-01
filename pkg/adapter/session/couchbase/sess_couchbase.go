package couchbase

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

// couchbase for session provider
//
// depend on github.com/couchbaselabs/go-couchbasee
//
// go install github.com/couchbaselabs/go-couchbase
//
// Usage:
// import(
//   _ "github.com/bhojpur/session/pkg/provider/couchbase"
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("couchbase", ``{"cookieName":"bsessionid","gclifetime":3600,"ProviderConfig":"http://host:port/, Pool, Bucket"}``)
//		go globalSessions.GC()
//	}

import (
	"context"
	"net/http"

	bhojpurcb "github.com/bhojpur/session/pkg/provider/couchbase"
	"github.com/bhojpur/web/pkg/adapter/session"
)

// SessionStore store each session
type SessionStore bhojpurcb.SessionStore

// Provider couchabse provided
type Provider bhojpurcb.Provider

// Set value to couchabse session
func (cs *SessionStore) Set(key, value interface{}) error {
	return (*bhojpurcb.SessionStore)(cs).Set(context.Background(), key, value)
}

// Get value from couchabse session
func (cs *SessionStore) Get(key interface{}) interface{} {
	return (*bhojpurcb.SessionStore)(cs).Get(context.Background(), key)
}

// Delete value in couchbase session by given key
func (cs *SessionStore) Delete(key interface{}) error {
	return (*bhojpurcb.SessionStore)(cs).Delete(context.Background(), key)
}

// Flush Clean all values in couchbase session
func (cs *SessionStore) Flush() error {
	return (*bhojpurcb.SessionStore)(cs).Flush(context.Background())
}

// SessionID Get couchbase session store id
func (cs *SessionStore) SessionID() string {
	return (*bhojpurcb.SessionStore)(cs).SessionID(context.Background())
}

// SessionRelease Write couchbase session with Gob string
func (cs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*bhojpurcb.SessionStore)(cs).SessionRelease(context.Background(), w)
}

// SessionInit init couchbase session
// savepath like couchbase server REST/JSON URL
// e.g. http://host:port/, Pool, Bucket
func (cp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*bhojpurcb.Provider)(cp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read couchbase session by sid
func (cp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*bhojpurcb.Provider)(cp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist Check couchbase session exist.
// it checkes sid exist or not.
func (cp *Provider) SessionExist(sid string) bool {
	res, _ := (*bhojpurcb.Provider)(cp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate remove oldsid and use sid to generate new session
func (cp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*bhojpurcb.Provider)(cp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy Remove bucket in this couchbase
func (cp *Provider) SessionDestroy(sid string) error {
	return (*bhojpurcb.Provider)(cp).SessionDestroy(context.Background(), sid)
}

// SessionGC Recycle
func (cp *Provider) SessionGC() {
	(*bhojpurcb.Provider)(cp).SessionGC(context.Background())
}

// SessionAll return all active session
func (cp *Provider) SessionAll() int {
	return (*bhojpurcb.Provider)(cp).SessionAll(context.Background())
}
