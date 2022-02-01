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

// session provider
//
// Usage:
// import(
//   session "github.com/bhojpur/session/pkg/engine"
// )
//
//	func init() {
//      globalSessions, _ = session.NewManager("memory", `{"cookieName":"bsessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "cookieLifeTime": 3600, "providerConfig": ""}`)
//		go globalSessions.GC()
//	}

import (
	"io"
	"net/http"
	"os"

	session "github.com/bhojpur/session/pkg/engine"
)

// Store contains all data for one session process with specific id.
type Store interface {
	Set(key, value interface{}) error     // set session value
	Get(key interface{}) interface{}      // get session value
	Delete(key interface{}) error         // delete session value
	SessionID() string                    // back current sessionID
	SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
	Flush() error                         // delete all data
}

// Provider contains global session methods and saved SessionStores.
// it can operate a SessionStore by its id.
type Provider interface {
	SessionInit(gclifetime int64, config string) error
	SessionRead(sid string) (Store, error)
	SessionExist(sid string) bool
	SessionRegenerate(oldsid, sid string) (Store, error)
	SessionDestroy(sid string) error
	SessionAll() int // get all active session
	SessionGC()
}

// SLogger a helpful variable to log information about session
var SLogger = NewSessionLog(os.Stderr)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provide Provider) {
	session.Register(name, &oldToNewProviderAdapter{
		delegate: provide,
	})
}

// GetProvider
func GetProvider(name string) (Provider, error) {
	res, err := session.GetProvider(name)
	if adt, ok := res.(*oldToNewProviderAdapter); err == nil && ok {
		return adt.delegate, err
	}

	return &newToOldProviderAdapter{
		delegate: res,
	}, err
}

// ManagerConfig define the session config
type ManagerConfig session.ManagerConfig

// Manager contains Provider and its configuration.
type Manager session.Manager

// NewManager Create new Manager with provider name and json config string.
// provider name:
// 1. cookie
// 2. file
// 3. memory
// 4. redis
// 5. mysql
// json config:
// 1. is https  default false
// 2. hashfunc  default sha1
// 3. hashkey default bhojpursessionkey
// 4. maxage default is none
func NewManager(provideName string, cf *ManagerConfig) (*Manager, error) {
	m, err := session.NewManager(provideName, (*session.ManagerConfig)(cf))
	return (*Manager)(m), err
}

// GetProvider return current manager's provider
func (manager *Manager) GetProvider() Provider {
	return &newToOldProviderAdapter{
		delegate: (*session.Manager)(manager).GetProvider(),
	}
}

// SessionStart generate or read the session id from http request.
// if session id exists, return SessionStore with this id.
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Store, error) {
	s, err := (*session.Manager)(manager).SessionStart(w, r)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionDestroy Destroy session by its id in http request cookie.
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	(*session.Manager)(manager).SessionDestroy(w, r)
}

// GetSessionStore Get SessionStore by its id.
func (manager *Manager) GetSessionStore(sid string) (Store, error) {
	s, err := (*session.Manager)(manager).GetSessionStore(sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// GC Start session gc process.
// it can do gc in times after gc lifetime.
func (manager *Manager) GC() {
	(*session.Manager)(manager).GC()
}

// SessionRegenerateID Regenerate a session id for this SessionStore who's id is saving in http request.
func (manager *Manager) SessionRegenerateID(w http.ResponseWriter, r *http.Request) Store {
	s, _ := (*session.Manager)(manager).SessionRegenerateID(w, r)
	return &NewToOldStoreAdapter{
		delegate: s,
	}
}

// GetActiveSession Get all active sessions count number.
func (manager *Manager) GetActiveSession() int {
	return (*session.Manager)(manager).GetActiveSession()
}

// SetSecure Set cookie with https.
func (manager *Manager) SetSecure(secure bool) {
	(*session.Manager)(manager).SetSecure(secure)
}

// Log implement the log.Logger
type Log session.Log

// NewSessionLog set io.Writer to create a Logger for session.
func NewSessionLog(out io.Writer) *Log {
	return (*Log)(session.NewSessionLog(out))
}
