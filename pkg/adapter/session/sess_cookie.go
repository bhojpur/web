package session

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

	session "github.com/bhojpur/session/pkg/engine"
)

// CookieSessionStore Cookie SessionStore
type CookieSessionStore session.CookieSessionStore

// Set value to cookie session.
// the value are encoded as gob with hash block string.
func (st *CookieSessionStore) Set(key, value interface{}) error {
	return (*session.CookieSessionStore)(st).Set(context.Background(), key, value)
}

// Get value from cookie session
func (st *CookieSessionStore) Get(key interface{}) interface{} {
	return (*session.CookieSessionStore)(st).Get(context.Background(), key)
}

// Delete value in cookie session
func (st *CookieSessionStore) Delete(key interface{}) error {
	return (*session.CookieSessionStore)(st).Delete(context.Background(), key)
}

// Flush Clean all values in cookie session
func (st *CookieSessionStore) Flush() error {
	return (*session.CookieSessionStore)(st).Flush(context.Background())
}

// SessionID Return id of this cookie session
func (st *CookieSessionStore) SessionID() string {
	return (*session.CookieSessionStore)(st).SessionID(context.Background())
}

// SessionRelease Write cookie session to http response cookie
func (st *CookieSessionStore) SessionRelease(w http.ResponseWriter) {
	(*session.CookieSessionStore)(st).SessionRelease(context.Background(), w)
}

// CookieProvider Cookie session provider
type CookieProvider session.CookieProvider

// SessionInit Init cookie session provider with max lifetime and config json.
// maxlifetime is ignored.
// json config:
// 	securityKey - hash string
// 	blockKey - gob encode hash string. it's saved as aes crypto.
// 	securityName - recognized name in encoded cookie string
// 	cookieName - cookie name
// 	maxage - cookie max life time.
func (pder *CookieProvider) SessionInit(maxlifetime int64, config string) error {
	return (*session.CookieProvider)(pder).SessionInit(context.Background(), maxlifetime, config)
}

// SessionRead Get SessionStore in cooke.
// decode cooke string to map and put into SessionStore with sid.
func (pder *CookieProvider) SessionRead(sid string) (Store, error) {
	s, err := (*session.CookieProvider)(pder).SessionRead(context.Background(), sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionExist Cookie session is always existed
func (pder *CookieProvider) SessionExist(sid string) bool {
	res, _ := (*session.CookieProvider)(pder).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate Implement method, no used.
func (pder *CookieProvider) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := (*session.CookieProvider)(pder).SessionRegenerate(context.Background(), oldsid, sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionDestroy Implement method, no used.
func (pder *CookieProvider) SessionDestroy(sid string) error {
	return (*session.CookieProvider)(pder).SessionDestroy(context.Background(), sid)
}

// SessionGC Implement method, no used.
func (pder *CookieProvider) SessionGC() {
	(*session.CookieProvider)(pder).SessionGC(context.Background())
}

// SessionAll Implement method, return 0.
func (pder *CookieProvider) SessionAll() int {
	return (*session.CookieProvider)(pder).SessionAll(context.Background())
}

// SessionUpdate Implement method, no used.
func (pder *CookieProvider) SessionUpdate(sid string) error {
	return (*session.CookieProvider)(pder).SessionUpdate(context.Background(), sid)
}
