package ledis

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
	"net/http"

	bhojpurLedis "github.com/bhojpur/session/pkg/provider/ledis"
	"github.com/bhojpur/web/pkg/adapter/session"
)

// SessionStore ledis session store
type SessionStore bhojpurLedis.SessionStore

// Set value in ledis session
func (ls *SessionStore) Set(key, value interface{}) error {
	return (*bhojpurLedis.SessionStore)(ls).Set(context.Background(), key, value)
}

// Get value in ledis session
func (ls *SessionStore) Get(key interface{}) interface{} {
	return (*bhojpurLedis.SessionStore)(ls).Get(context.Background(), key)
}

// Delete value in ledis session
func (ls *SessionStore) Delete(key interface{}) error {
	return (*bhojpurLedis.SessionStore)(ls).Delete(context.Background(), key)
}

// Flush clear all values in ledis session
func (ls *SessionStore) Flush() error {
	return (*bhojpurLedis.SessionStore)(ls).Flush(context.Background())
}

// SessionID get ledis session id
func (ls *SessionStore) SessionID() string {
	return (*bhojpurLedis.SessionStore)(ls).SessionID(context.Background())
}

// SessionRelease save session values to ledis
func (ls *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*bhojpurLedis.SessionStore)(ls).SessionRelease(context.Background(), w)
}

// Provider ledis session provider
type Provider bhojpurLedis.Provider

// SessionInit init ledis session
// savepath like ledis server saveDataPath,pool size
// e.g. 127.0.0.1:6379,100,bhojpur
func (lp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*bhojpurLedis.Provider)(lp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read ledis session by sid
func (lp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*bhojpurLedis.Provider)(lp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check ledis session exist by sid
func (lp *Provider) SessionExist(sid string) bool {
	res, _ := (*bhojpurLedis.Provider)(lp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for ledis session
func (lp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*bhojpurLedis.Provider)(lp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete ledis session by id
func (lp *Provider) SessionDestroy(sid string) error {
	return (*bhojpurLedis.Provider)(lp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (lp *Provider) SessionGC() {
	(*bhojpurLedis.Provider)(lp).SessionGC(context.Background())
}

// SessionAll return all active session
func (lp *Provider) SessionAll() int {
	return (*bhojpurLedis.Provider)(lp).SessionAll(context.Background())
}
