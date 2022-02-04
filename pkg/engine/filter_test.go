package engine

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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhojpur/web/pkg/context"
)

var FilterUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am " + ctx.Input.Param(":last") + ctx.Input.Param(":first")))
}

func TestFilter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/person/bhoj/pur", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/person/:last/:first", BeforeRouter, FilterUser)
	handler.Add("/person/:last/:first", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am bhojpur" {
		t.Errorf("user define func can't run")
	}
}

var FilterAdminUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am admin"))
}

// Filter pattern /admin/:all
// all url like    /admin/    /admin/pur    will all get filter

func TestPatternTwo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/admin/", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/?:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/ can't run")
	}
}

func TestPatternThree(t *testing.T) {
	r, _ := http.NewRequest("GET", "/admin/bhojpur", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/bhojpur can't run")
	}
}
