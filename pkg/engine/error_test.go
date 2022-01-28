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
	"strconv"
	"strings"
	"testing"
)

type errorTestController struct {
	Controller
}

const parseCodeError = "parse code error"

func (ec *errorTestController) Get() {
	errorCode, err := ec.GetInt("code")
	if err != nil {
		ec.Abort(parseCodeError)
	}
	if errorCode != 0 {
		ec.CustomAbort(errorCode, ec.GetString("code"))
	}
	ec.Abort("404")
}

func TestErrorCode_01(t *testing.T) {
	registerDefaultErrorHandler()
	for k := range ErrorMaps {
		r, _ := http.NewRequest("GET", "/error?code="+k, nil)
		w := httptest.NewRecorder()

		handler := NewControllerRegister()
		handler.Add("/error", &errorTestController{})
		handler.ServeHTTP(w, r)
		code, _ := strconv.Atoi(k)
		if w.Code != code {
			t.Fail()
		}
		if !strings.Contains(w.Body.String(), http.StatusText(code)) {
			t.Fail()
		}
	}
}

func TestErrorCode_02(t *testing.T) {
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "/error?code=0", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/error", &errorTestController{})
	handler.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Fail()
	}
}

func TestErrorCode_03(t *testing.T) {
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "/error?code=panic", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/error", &errorTestController{})
	handler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fail()
	}
	if w.Body.String() != parseCodeError {
		t.Fail()
	}
}
