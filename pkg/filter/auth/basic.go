package auth

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

// Package auth provides handlers to enable basic auth support.
// Simple Usage:
//	import(
//		"github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/filter/auth"
//	)
//
//	func main(){
//		// authenticate every request
//		bhojpur.InsertFilter("*", bhojpur.BeforeRouter,auth.Basic("username","secretpassword"))
//		bhojpur.Run()
//	}
//
//
// Advanced Usage:
//
//	func SecretAuth(username, password string) bool {
//		return username == "bhojpur" && password == "welcome"
//	}
//	authPlugin := auth.NewBasicAuthenticator(SecretAuth, "Authorization Required")
//	bhojpur.InsertFilter("*", bhojpur.BeforeRouter,authPlugin)

import (
	"encoding/base64"
	"net/http"
	"strings"

	ctxsvr "github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

var defaultRealm = "Authorization Required"

// Basic is the http basic auth
func Basic(username string, password string) websvr.FilterFunc {
	secrets := func(user, pass string) bool {
		return user == username && pass == password
	}
	return NewBasicAuthenticator(secrets, defaultRealm)
}

// NewBasicAuthenticator return the BasicAuth
func NewBasicAuthenticator(secrets SecretProvider, Realm string) websvr.FilterFunc {
	return func(ctx *ctxsvr.Context) {
		a := &BasicAuth{Secrets: secrets, Realm: Realm}
		if username := a.CheckAuth(ctx.Request); username == "" {
			a.RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	}
}

// SecretProvider is the SecretProvider function
type SecretProvider func(user, pass string) bool

// BasicAuth store the SecretProvider and Realm
type BasicAuth struct {
	Secrets SecretProvider
	Realm   string
}

// CheckAuth Checks the username/password combination from the request. Returns
// either an empty string (authentication failed) or the name of the
// authenticated user.
// Supports MD5 and SHA1 password entries
func (a *BasicAuth) CheckAuth(r *http.Request) string {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return ""
	}

	if a.Secrets(pair[0], pair[1]) {
		return pair[0]
	}
	return ""
}

// RequireAuth http.Handler for BasicAuth which initiates the authentication process
// (or requires reauthentication).
func (a *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+a.Realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}
