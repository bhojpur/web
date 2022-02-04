package httplib

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

// It is used as http.Client
// Usage:
//
// import "github.com/bhojpur/web/pkg/client/httplib"
//
//	b := httplib.Post("http://bhojpur.net/")
//	b.Param("username","bhojpur")
//	b.Param("password","123456")
//	b.PostFile("uploadfile1", "httplib.pdf")
//	b.PostFile("uploadfile2", "httplib.txt")
//	str, err := b.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(str)

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/bhojpur/web/pkg/client/httplib"
)

// SetDefaultSetting Overwrite default settings
func SetDefaultSetting(setting BhojpurHTTPSettings) {
	httplib.SetDefaultSetting(httplib.BhojpurHTTPSettings(setting))
}

// NewBhojpurRequest return *BhojpurHttpRequest with specific method
func NewBhojpurRequest(rawurl, method string) *BhojpurHTTPRequest {
	return &BhojpurHTTPRequest{
		delegate: httplib.NewBhojpurRequest(rawurl, method),
	}
}

// Get returns *BhojpurHttpRequest with GET method.
func Get(url string) *BhojpurHTTPRequest {
	return NewBhojpurRequest(url, "GET")
}

// Post returns *BhojpurHttpRequest with POST method.
func Post(url string) *BhojpurHTTPRequest {
	return NewBhojpurRequest(url, "POST")
}

// Put returns *BhojpurHttpRequest with PUT method.
func Put(url string) *BhojpurHTTPRequest {
	return NewBhojpurRequest(url, "PUT")
}

// Delete returns *BhojpurHttpRequest DELETE method.
func Delete(url string) *BhojpurHTTPRequest {
	return NewBhojpurRequest(url, "DELETE")
}

// Head returns *BhojpurHttpRequest with HEAD method.
func Head(url string) *BhojpurHTTPRequest {
	return NewBhojpurRequest(url, "HEAD")
}

// BhojpurHTTPSettings is the http.Client setting
type BhojpurHTTPSettings httplib.BhojpurHTTPSettings

// BhojpurHTTPRequest provides more useful methods for requesting one url than http.Request.
type BhojpurHTTPRequest struct {
	delegate *httplib.BhojpurHTTPRequest
}

// GetRequest return the request object
func (b *BhojpurHTTPRequest) GetRequest() *http.Request {
	return b.delegate.GetRequest()
}

// Setting Change request settings
func (b *BhojpurHTTPRequest) Setting(setting BhojpurHTTPSettings) *BhojpurHTTPRequest {
	b.delegate.Setting(httplib.BhojpurHTTPSettings(setting))
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *BhojpurHTTPRequest) SetBasicAuth(username, password string) *BhojpurHTTPRequest {
	b.delegate.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *BhojpurHTTPRequest) SetEnableCookie(enable bool) *BhojpurHTTPRequest {
	b.delegate.SetEnableCookie(enable)
	return b
}

// SetUserAgent sets User-Agent header field
func (b *BhojpurHTTPRequest) SetUserAgent(useragent string) *BhojpurHTTPRequest {
	b.delegate.SetUserAgent(useragent)
	return b
}

// Retries sets Retries times.
// default is 0 means no retried.
// -1 means retried forever.
// others means retried times.
func (b *BhojpurHTTPRequest) Retries(times int) *BhojpurHTTPRequest {
	b.delegate.Retries(times)
	return b
}

func (b *BhojpurHTTPRequest) RetryDelay(delay time.Duration) *BhojpurHTTPRequest {
	b.delegate.RetryDelay(delay)
	return b
}

// SetTimeout sets connect time out and read-write time out for BhojpurRequest.
func (b *BhojpurHTTPRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *BhojpurHTTPRequest {
	b.delegate.SetTimeout(connectTimeout, readWriteTimeout)
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *BhojpurHTTPRequest) SetTLSClientConfig(config *tls.Config) *BhojpurHTTPRequest {
	b.delegate.SetTLSClientConfig(config)
	return b
}

// Header add header item string in request.
func (b *BhojpurHTTPRequest) Header(key, value string) *BhojpurHTTPRequest {
	b.delegate.Header(key, value)
	return b
}

// SetHost set the request host
func (b *BhojpurHTTPRequest) SetHost(host string) *BhojpurHTTPRequest {
	b.delegate.SetHost(host)
	return b
}

// SetProtocolVersion Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *BhojpurHTTPRequest) SetProtocolVersion(vers string) *BhojpurHTTPRequest {
	b.delegate.SetProtocolVersion(vers)
	return b
}

// SetCookie add cookie into request.
func (b *BhojpurHTTPRequest) SetCookie(cookie *http.Cookie) *BhojpurHTTPRequest {
	b.delegate.SetCookie(cookie)
	return b
}

// SetTransport set the setting transport
func (b *BhojpurHTTPRequest) SetTransport(transport http.RoundTripper) *BhojpurHTTPRequest {
	b.delegate.SetTransport(transport)
	return b
}

// SetProxy set the http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (b *BhojpurHTTPRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *BhojpurHTTPRequest {
	b.delegate.SetProxy(proxy)
	return b
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (b *BhojpurHTTPRequest) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *BhojpurHTTPRequest {
	b.delegate.SetCheckRedirect(redirect)
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *BhojpurHTTPRequest) Param(key, value string) *BhojpurHTTPRequest {
	b.delegate.Param(key, value)
	return b
}

// PostFile add a post file to the request
func (b *BhojpurHTTPRequest) PostFile(formname, filename string) *BhojpurHTTPRequest {
	b.delegate.PostFile(formname, filename)
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *BhojpurHTTPRequest) Body(data interface{}) *BhojpurHTTPRequest {
	b.delegate.Body(data)
	return b
}

// XMLBody adds request raw body encoding by XML.
func (b *BhojpurHTTPRequest) XMLBody(obj interface{}) (*BhojpurHTTPRequest, error) {
	_, err := b.delegate.XMLBody(obj)
	return b, err
}

// YAMLBody adds request raw body encoding by YAML.
func (b *BhojpurHTTPRequest) YAMLBody(obj interface{}) (*BhojpurHTTPRequest, error) {
	_, err := b.delegate.YAMLBody(obj)
	return b, err
}

// JSONBody adds request raw body encoding by JSON.
func (b *BhojpurHTTPRequest) JSONBody(obj interface{}) (*BhojpurHTTPRequest, error) {
	_, err := b.delegate.JSONBody(obj)
	return b, err
}

// DoRequest will do the client.Do
func (b *BhojpurHTTPRequest) DoRequest() (resp *http.Response, err error) {
	return b.delegate.DoRequest()
}

// String returns the body string in response.
// it calls Response inner.
func (b *BhojpurHTTPRequest) String() (string, error) {
	return b.delegate.String()
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *BhojpurHTTPRequest) Bytes() ([]byte, error) {
	return b.delegate.Bytes()
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *BhojpurHTTPRequest) ToFile(filename string) error {
	return b.delegate.ToFile(filename)
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *BhojpurHTTPRequest) ToJSON(v interface{}) error {
	return b.delegate.ToJSON(v)
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *BhojpurHTTPRequest) ToXML(v interface{}) error {
	return b.delegate.ToXML(v)
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (b *BhojpurHTTPRequest) ToYAML(v interface{}) error {
	return b.delegate.ToYAML(v)
}

// Response executes request client gets response mannually.
func (b *BhojpurHTTPRequest) Response() (*http.Response, error) {
	return b.delegate.Response()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return httplib.TimeoutDialer(cTimeout, rwTimeout)
}
