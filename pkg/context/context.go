package context

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

// Usage:
//
//	import "github.com/bhojpur/web/pkg/context"
//
//	ctx := context.Context{Request:req,ResponseWriter:rw}

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Commonly used mime-types
const (
	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
	ApplicationYAML = "application/x-yaml"
	TextXML         = "text/xml"
)

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return &Context{
		Input:  NewInput(),
		Output: NewOutput(),
	}
}

// Context Http request context struct including BhojpurInput, BhojpurOutput, http.Request and http.ResponseWriter.
// BhojpurInput and BhojpurOutput provides an api to operate request and response more easily.
type Context struct {
	Input          *BhojpurInput
	Output         *BhojpurOutput
	Request        *http.Request
	ResponseWriter *Response
	_xsrfToken     string
}

// Reset initializes Context, BhojpurInput and BhojpurOutput
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	ctx.Request = r
	if ctx.ResponseWriter == nil {
		ctx.ResponseWriter = &Response{}
	}
	ctx.ResponseWriter.reset(rw)
	ctx.Input.Reset(ctx)
	ctx.Output.Reset(ctx)
	ctx._xsrfToken = ""
}

// Redirect redirects to localurl with http header status code.
func (ctx *Context) Redirect(status int, localurl string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, localurl, status)
}

// Abort stops the request.
// If Bhojpur.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	ctx.Output.SetStatus(status)
	panic(body)
}

// WriteString writes a string to response body.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// GetCookie gets a cookie from a request for a given key.
// (Alias of BhojpurInput.Cookie)
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// SetCookie sets a cookie for a response.
// (Alias of BhojpurOutput.Cookie)
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}

// GetSecureCookie gets a secure cookie from a request for a given key.
func (ctx *Context) GetSecureCookie(Secret, key string) (string, bool) {
	val := ctx.Input.Cookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha256.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// SetSecureCookie sets a secure cookie for a response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha256.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.Output.Cookie(name, cookie, others...)
}

// XSRFToken creates and returns an xsrf token string
func (ctx *Context) XSRFToken(key string, expire int64) string {
	if ctx._xsrfToken == "" {
		token, ok := ctx.GetSecureCookie(key, "_xsrf")
		if !ok {
			token = string(token.RandomCreateBytes(32))
			// TODO make it configurable
			ctx.SetSecureCookie(key, "_xsrf", token, expire, "", "")
		}
		ctx._xsrfToken = token
	}
	return ctx._xsrfToken
}

// CheckXSRFCookie checks if the XSRF token in this request is valid or not.
// The token can be provided in the request header in the form "X-Xsrftoken" or "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) CheckXSRFCookie() bool {
	token := ctx.Input.Query("_xsrf")
	if token == "" {
		token = ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		ctx.Abort(422, "422")
		return false
	}
	if ctx._xsrfToken != token {
		ctx.Abort(417, "417")
		return false
	}
	return true
}

// RenderMethodResult renders the return value of a controller method to the output
func (ctx *Context) RenderMethodResult(result interface{}) {
	if result != nil {
		renderer, ok := result.(Renderer)
		if !ok {
			err, ok := result.(error)
			if ok {
				renderer = errorRenderer(err)
			} else {
				renderer = jsonRenderer(result)
			}
		}
		renderer.Render(ctx)
	}
}

// Response is a wrapper for the http.ResponseWriter
// Started:  if true, response was already written to so the other handler will not be executed
type Response struct {
	http.ResponseWriter
	Started bool
	Status  int
	Elapsed time.Duration
}

func (r *Response) reset(rw http.ResponseWriter) {
	r.ResponseWriter = rw
	r.Status = 0
	r.Started = false
}

// Write writes the data to the connection as part of a HTTP reply,
// and sets `Started` to true.
// Started:  if true, the response was already sent
func (r *Response) Write(p []byte) (int, error) {
	r.Started = true
	return r.ResponseWriter.Write(p)
}

// WriteHeader sends a HTTP response header with status code,
// and sets `Started` to true.
func (r *Response) WriteHeader(code int) {
	if r.Status > 0 {
		//prevent multiple response.WriteHeader calls
		return
	}
	r.Status = code
	r.Started = true
	r.ResponseWriter.WriteHeader(code)
}

// Hijack hijacker for http
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

// Flush http.Flusher
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify http.CloseNotifier
func (r *Response) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

// Pusher http.Pusher
func (r *Response) Pusher() (pusher http.Pusher) {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
