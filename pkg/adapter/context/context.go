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

import (
	"bufio"
	"net"
	"net/http"

	"github.com/bhojpur/web/pkg/context"
)

// commonly used mime-types
const (
	ApplicationJSON = context.ApplicationJSON
	ApplicationXML  = context.ApplicationXML
	ApplicationYAML = context.ApplicationYAML
	TextXML         = context.TextXML
)

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return (*Context)(context.NewContext())
}

// Context Http request context struct including BhojpurInput, BhojpurOutput, http.Request and http.ResponseWriter.
// BhojpurInput and BhojpurOutput provides some api to operate request and response more easily.
type Context context.Context

// Reset init Context, BhojpurInput and BhojpurOutput
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	(*context.Context)(ctx).Reset(rw, r)
}

// Redirect does redirection to localurl with http header status code.
func (ctx *Context) Redirect(status int, localurl string) {
	(*context.Context)(ctx).Redirect(status, localurl)
}

// Abort stops this request.
// if Bhojpur.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	(*context.Context)(ctx).Abort(status, body)
}

// WriteString Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	(*context.Context)(ctx).WriteString(content)
}

// GetCookie Get cookie from request by a given key.
// It's alias of BhojpurInput.Cookie.
func (ctx *Context) GetCookie(key string) string {
	return (*context.Context)(ctx).GetCookie(key)
}

// SetCookie Set cookie for response.
// It's alias of BhojpurOutput.Cookie.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	(*context.Context)(ctx).SetCookie(name, value, others)
}

// GetSecureCookie Get secure cookie from request by a given key.
func (ctx *Context) GetSecureCookie(Secret, key string) (string, bool) {
	return (*context.Context)(ctx).GetSecureCookie(Secret, key)
}

// SetSecureCookie Set Secure cookie for response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	(*context.Context)(ctx).SetSecureCookie(Secret, name, value, others)
}

// XSRFToken creates a xsrf token string and returns.
func (ctx *Context) XSRFToken(key string, expire int64) string {
	return (*context.Context)(ctx).XSRFToken(key, expire)
}

// CheckXSRFCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) CheckXSRFCookie() bool {
	return (*context.Context)(ctx).CheckXSRFCookie()
}

// RenderMethodResult renders the return value of a controller method to the output
func (ctx *Context) RenderMethodResult(result interface{}) {
	(*context.Context)(ctx).RenderMethodResult(result)
}

// Response is a wrapper for the http.ResponseWriter
// started set to true if response was written to then don't execute other handler
type Response context.Response

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (r *Response) Write(p []byte) (int, error) {
	return (*context.Response)(r).Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (r *Response) WriteHeader(code int) {
	(*context.Response)(r).WriteHeader(code)
}

// Hijack hijacker for http
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return (*context.Response)(r).Hijack()
}

// Flush http.Flusher
func (r *Response) Flush() {
	(*context.Response)(r).Flush()
}

// CloseNotify http.CloseNotifier
func (r *Response) CloseNotify() <-chan bool {
	return (*context.Response)(r).CloseNotify()
}

// Pusher http.Pusher
func (r *Response) Pusher() (pusher http.Pusher) {
	return (*context.Response)(r).Pusher()
}
