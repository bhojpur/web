package cors

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
	"time"

	"github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// HTTPHeaderGuardRecorder is httptest.ResponseRecorder with own http.Header
type HTTPHeaderGuardRecorder struct {
	*httptest.ResponseRecorder
	savedHeaderMap http.Header
}

// NewRecorder return HttpHeaderGuardRecorder
func NewRecorder() *HTTPHeaderGuardRecorder {
	return &HTTPHeaderGuardRecorder{httptest.NewRecorder(), nil}
}

func (gr *HTTPHeaderGuardRecorder) WriteHeader(code int) {
	gr.ResponseRecorder.WriteHeader(code)
	gr.savedHeaderMap = gr.ResponseRecorder.Header()
}

func (gr *HTTPHeaderGuardRecorder) Header() http.Header {
	if gr.savedHeaderMap != nil {
		// headers were written. clone so we don't get updates
		clone := make(http.Header)
		for k, v := range gr.savedHeaderMap {
			clone[k] = v
		}
		return clone
	}
	return gr.ResponseRecorder.Header()
}

func Test_AllowAll(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowAllOrigins: true,
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	r, _ := http.NewRequest("PUT", "/foo", nil)
	handler.ServeHTTP(recorder, r)

	if recorder.HeaderMap.Get(headerAllowOrigin) != "*" {
		t.Errorf("Allow-Origin header should be *")
	}
}

func Test_AllowRegexMatch(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowOrigins: []string{"https://aaa.com", "https://*.bhojpur.net"},
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	origin := "https://bar.bhojpur.net"
	r, _ := http.NewRequest("PUT", "/foo", nil)
	r.Header.Add("Origin", origin)
	handler.ServeHTTP(recorder, r)

	headerValue := recorder.HeaderMap.Get(headerAllowOrigin)
	if headerValue != origin {
		t.Errorf("Allow-Origin header should be %v, found %v", origin, headerValue)
	}
}

func Test_AllowRegexNoMatch(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowOrigins: []string{"https://*.bhojpur.net"},
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	origin := "https://ww.bhojpur.net.evil.com"
	r, _ := http.NewRequest("PUT", "/foo", nil)
	r.Header.Add("Origin", origin)
	handler.ServeHTTP(recorder, r)

	headerValue := recorder.HeaderMap.Get(headerAllowOrigin)
	if headerValue != "" {
		t.Errorf("Allow-Origin header should not exist, found %v", headerValue)
	}
}

func Test_OtherHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PATCH", "GET"},
		AllowHeaders:     []string{"Origin", "X-whatever"},
		ExposeHeaders:    []string{"Content-Length", "Hello"},
		MaxAge:           5 * time.Minute,
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	r, _ := http.NewRequest("PUT", "/foo", nil)
	handler.ServeHTTP(recorder, r)

	credentialsVal := recorder.HeaderMap.Get(headerAllowCredentials)
	methodsVal := recorder.HeaderMap.Get(headerAllowMethods)
	headersVal := recorder.HeaderMap.Get(headerAllowHeaders)
	exposedHeadersVal := recorder.HeaderMap.Get(headerExposeHeaders)
	maxAgeVal := recorder.HeaderMap.Get(headerMaxAge)

	if credentialsVal != "true" {
		t.Errorf("Allow-Credentials is expected to be true, found %v", credentialsVal)
	}

	if methodsVal != "PATCH,GET" {
		t.Errorf("Allow-Methods is expected to be PATCH,GET; found %v", methodsVal)
	}

	if headersVal != "Origin,X-whatever" {
		t.Errorf("Allow-Headers is expected to be Origin,X-whatever; found %v", headersVal)
	}

	if exposedHeadersVal != "Content-Length,Hello" {
		t.Errorf("Expose-Headers are expected to be Content-Length,Hello. Found %v", exposedHeadersVal)
	}

	if maxAgeVal != "300" {
		t.Errorf("Max-Age is expected to be 300, found %v", maxAgeVal)
	}
}

func Test_DefaultAllowHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowAllOrigins: true,
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})

	r, _ := http.NewRequest("PUT", "/foo", nil)
	handler.ServeHTTP(recorder, r)

	headersVal := recorder.HeaderMap.Get(headerAllowHeaders)
	if headersVal != "Origin,Accept,Content-Type,Authorization" {
		t.Errorf("Allow-Headers is expected to be Origin,Accept,Content-Type,Authorization; found %v", headersVal)
	}
}

func Test_Preflight(t *testing.T) {
	recorder := NewRecorder()
	handler := websvr.NewControllerRegister()
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH"},
		AllowHeaders:    []string{"Origin", "X-whatever", "X-CaseSensitive"},
	}))

	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	r, _ := http.NewRequest("OPTIONS", "/foo", nil)
	r.Header.Add(headerRequestMethod, "PUT")
	r.Header.Add(headerRequestHeaders, "X-whatever, x-casesensitive")
	handler.ServeHTTP(recorder, r)

	headers := recorder.Header()
	methodsVal := headers.Get(headerAllowMethods)
	headersVal := headers.Get(headerAllowHeaders)
	originVal := headers.Get(headerAllowOrigin)

	if methodsVal != "PUT,PATCH" {
		t.Errorf("Allow-Methods is expected to be PUT,PATCH, found %v", methodsVal)
	}

	if !strings.Contains(headersVal, "X-whatever") {
		t.Errorf("Allow-Headers is expected to contain X-whatever, found %v", headersVal)
	}

	if !strings.Contains(headersVal, "x-casesensitive") {
		t.Errorf("Allow-Headers is expected to contain x-casesensitive, found %v", headersVal)
	}

	if originVal != "*" {
		t.Errorf("Allow-Origin is expected to be *, found %v", originVal)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Status code is expected to be 200, found %d", recorder.Code)
	}
}

func Benchmark_WithoutCORS(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	websvr.BConfig.RunMode = websvr.PROD
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(recorder, r)
	}
}

func Benchmark_WithCORS(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	websvr.BConfig.RunMode = websvr.PROD
	handler.InsertFilter("*", websvr.BeforeRouter, Allow(&Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PATCH", "GET"},
		AllowHeaders:     []string{"Origin", "X-whatever"},
		MaxAge:           5 * time.Minute,
	}))
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(recorder, r)
	}
}
