package ssdb

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

	"github.com/bhojpur/web/pkg/adapter/session"

	bhojpurSsdb "github.com/bhojpur/session/pkg/provider/ssdb"
)

// Provider holds ssdb client and configs
type Provider bhojpurSsdb.Provider

// SessionInit init the ssdb with the config
func (p *Provider) SessionInit(maxLifetime int64, savePath string) error {
	return (*bhojpurSsdb.Provider)(p).SessionInit(context.Background(), maxLifetime, savePath)
}

// SessionRead return a ssdb client session Store
func (p *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*bhojpurSsdb.Provider)(p).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist judged whether sid is exist in session
func (p *Provider) SessionExist(sid string) bool {
	res, _ := (*bhojpurSsdb.Provider)(p).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate regenerate session with new sid and delete oldsid
func (p *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*bhojpurSsdb.Provider)(p).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy destroy the sid
func (p *Provider) SessionDestroy(sid string) error {
	return (*bhojpurSsdb.Provider)(p).SessionDestroy(context.Background(), sid)
}

// SessionGC not implemented
func (p *Provider) SessionGC() {
	(*bhojpurSsdb.Provider)(p).SessionGC(context.Background())
}

// SessionAll not implemented
func (p *Provider) SessionAll() int {
	return (*bhojpurSsdb.Provider)(p).SessionAll(context.Background())
}

// SessionStore holds the session information which stored in ssdb
type SessionStore bhojpurSsdb.SessionStore

// Set the key and value
func (s *SessionStore) Set(key, value interface{}) error {
	return (*bhojpurSsdb.SessionStore)(s).Set(context.Background(), key, value)
}

// Get return the value by the key
func (s *SessionStore) Get(key interface{}) interface{} {
	return (*bhojpurSsdb.SessionStore)(s).Get(context.Background(), key)
}

// Delete the key in session store
func (s *SessionStore) Delete(key interface{}) error {
	return (*bhojpurSsdb.SessionStore)(s).Delete(context.Background(), key)
}

// Flush delete all keys and values
func (s *SessionStore) Flush() error {
	return (*bhojpurSsdb.SessionStore)(s).Flush(context.Background())
}

// SessionID return the sessionID
func (s *SessionStore) SessionID() string {
	return (*bhojpurSsdb.SessionStore)(s).SessionID(context.Background())
}

// SessionRelease Store the keyvalues into ssdb
func (s *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*bhojpurSsdb.SessionStore)(s).SessionRelease(context.Background(), w)
}
