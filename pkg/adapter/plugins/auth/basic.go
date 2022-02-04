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
//		websvr "github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/plugins/auth"
//	)
//
//	func main(){
//		// authenticate every request
//		websvr.InsertFilter("*", websvr.BeforeRouter, auth.Basic("username","secretpassword"))
//		websvr.Run()
//	}
//
// Advanced Usage:
//
//	func SecretAuth(username, password string) bool {
//		return username == "bhojpur" && password == "helloBhojpur"
//	}
//	authPlugin := auth.NewBasicAuthenticator(SecretAuth, "Authorization Required")
//	websvr.InsertFilter("*", websvr.BeforeRouter,authPlugin)

import (
	"net/http"

	bhojpur "github.com/bhojpur/web/pkg/adapter"
	"github.com/bhojpur/web/pkg/adapter/context"
	ctxsvr "github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/filter/auth"
)

// Basic is the http basic auth
func Basic(username string, password string) bhojpur.FilterFunc {
	return func(c *context.Context) {
		f := auth.Basic(username, password)
		f((*ctxsvr.Context)(c))
	}
}

// NewBasicAuthenticator return the BasicAuth
func NewBasicAuthenticator(secrets SecretProvider, realm string) bhojpur.FilterFunc {
	f := auth.NewBasicAuthenticator(auth.SecretProvider(secrets), realm)
	return func(c *context.Context) {
		f((*ctxsvr.Context)(c))
	}
}

// SecretProvider is the SecretProvider function
type SecretProvider auth.SecretProvider

// BasicAuth store the SecretProvider and Realm
type BasicAuth auth.BasicAuth

// CheckAuth Checks the username/password combination from the request. Returns
// either an empty string (authentication failed) or the name of the
// authenticated user.
// Supports MD5 and SHA1 password entries
func (a *BasicAuth) CheckAuth(r *http.Request) string {
	return (*auth.BasicAuth)(a).CheckAuth(r)
}

// RequireAuth http.Handler for BasicAuth which initiates the authentication process
// (or requires reauthentication).
func (a *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	(*auth.BasicAuth)(a).RequireAuth(w, r)
}
