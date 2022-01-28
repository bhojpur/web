package authz

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

// It provides handlers to enable ACL, RBAC, ABAC authorization support.
// Simple Usage:
//	import(
//		"github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/filter/plugins/authz"
//		"github.com/casbin/casbin"
//	)
//
//	func main(){
//		// mediate the access for every request
//		bhojpur.InsertFilter("*", bhojpur.BeforeRouter, authz.NewAuthorizer(casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")))
//		bhojpur.Run()
//	}
//
// Advanced Usage:
//
//	func main(){
//		e := casbin.NewEnforcer("authz_model.conf", "")
//		e.AddRoleForUser("alice", "admin")
//		e.AddPolicy(...)
//
//		bhojpur.InsertFilter("*", bhojpur.BeforeRouter, authz.NewAuthorizer(e))
//		bhojpur.Run()
//	}

import (
	"net/http"

	"github.com/casbin/casbin"

	ctxsvr "github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// NewAuthorizer returns the authorizer.
// Use a casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) websvr.FilterFunc {
	return func(ctx *ctxsvr.Context) {
		a := &BasicAuthorizer{enforcer: e}

		if !a.CheckPermission(ctx.Request) {
			a.RequirePermission(ctx.ResponseWriter)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	user := a.GetUserName(r)
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(user, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
	w.WriteHeader(403)
	w.Write([]byte("403 Forbidden\n"))
}
