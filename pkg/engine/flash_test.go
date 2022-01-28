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
	"strings"
	"testing"
)

type TestFlashController struct {
	Controller
}

func (t *TestFlashController) TestWriteFlash() {
	flash := NewFlash()
	flash.Notice("TestFlashString")
	flash.Store(&t.Controller)
	// we choose to serve json because we don't want to load a template html file
	t.ServeJSON(true)
}

func TestFlashHeader(t *testing.T) {
	// create fake GET request
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// setup the handler
	handler := NewControllerRegister()
	handler.Add("/", &TestFlashController{}, "get:TestWriteFlash")
	handler.ServeHTTP(w, r)

	// get the Set-Cookie value
	sc := w.Header().Get("Set-Cookie")
	// match for the expected header
	res := strings.Contains(sc, "BHOJPUR_FLASH=%00notice%23BHOJPURFLASH%23TestFlashString%00")
	// validate the assertion
	if !res {
		t.Errorf("TestFlashHeader() unable to validate flash message")
	}
}
