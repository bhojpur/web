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

// Package authz provides handlers to enable ACL, RBAC, ABAC authorization support.
// Simple Usage:
//	import(
//		websvr "github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/plugins/authz"
//		plcsvr "github.com/bhojpur/policy/pkg/engine"
//	)
//
//	func main(){
//		// mediate the access for every request
//		websvr.InsertFilter("*", websvr.BeforeRouter, authz.NewAuthorizer(plcsvr.NewEnforcer("authz_model.conf", "authz_policy.csv")))
//		websvr.Run()
//	}
//
//
// Advanced Usage:
//
//	func main(){
//		e := plcsvr.NewEnforcer("authz_model.conf", "")
//		e.AddRoleForUser("alice", "admin")
//		e.AddPolicy(...)
//
//		websvr.InsertFilter("*", websvr.BeforeRouter, authz.NewAuthorizer(e))
//		websvr.Run()
//	}

import (
	"net/http"

	plcsvr "github.com/bhojpur/policy/pkg/engine"

	bhojpur "github.com/bhojpur/web/pkg/adapter"
	"github.com/bhojpur/web/pkg/adapter/context"
	ctxsvr "github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/filter/authz"
)

// NewAuthorizer returns the authorizer.
// Use a Bhojpur Policy enforcer as input
func NewAuthorizer(e *plcsvr.Enforcer) bhojpur.FilterFunc {
	f := authz.NewAuthorizer(e)
	return func(context *context.Context) {
		f((*ctxsvr.Context)(context))
	}
}

// BasicAuthorizer stores the Bhojpur Policy handler
type BasicAuthorizer authz.BasicAuthorizer

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	return (*authz.BasicAuthorizer)(a).GetUserName(r)
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) (bool, error) {
	return (*authz.BasicAuthorizer)(a).CheckPermission(r)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
	(*authz.BasicAuthorizer)(a).RequirePermission(w)
}
