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

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type (
	ClientOption             func(client *Client)
	BhojpurHTTPRequestOption func(request *BhojpurHTTPRequest)
)

// WithEnableCookie will enable cookie in all subsequent request
func WithEnableCookie(enable bool) ClientOption {
	return func(client *Client) {
		client.Setting.EnableCookie = enable
	}
}

// WithEnableCookie will adds UA in all subsequent request
func WithUserAgent(userAgent string) ClientOption {
	return func(client *Client) {
		client.Setting.UserAgent = userAgent
	}
}

// WithTLSClientConfig will adds tls config in all subsequent request
func WithTLSClientConfig(config *tls.Config) ClientOption {
	return func(client *Client) {
		client.Setting.TLSClientConfig = config
	}
}

// WithTransport will set transport field in all subsequent request
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(client *Client) {
		client.Setting.Transport = transport
	}
}

// WithProxy will set http proxy field in all subsequent request
func WithProxy(proxy func(*http.Request) (*url.URL, error)) ClientOption {
	return func(client *Client) {
		client.Setting.Proxy = proxy
	}
}

// WithCheckRedirect will specifies the policy for handling redirects in all subsequent request
func WithCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *Client) {
		client.Setting.CheckRedirect = redirect
	}
}

// WithHTTPSetting can replace BhojpurHTTPSeting
func WithHTTPSetting(setting BhojpurHTTPSettings) ClientOption {
	return func(client *Client) {
		client.Setting = setting
	}
}

// WithEnableGzip will enable gzip in all subsequent request
func WithEnableGzip(enable bool) ClientOption {
	return func(client *Client) {
		client.Setting.Gzip = enable
	}
}

// BhojpurHttpRequestOption

// WithTimeout sets connect time out and read-write time out for BhojpurRequest.
func WithTimeout(connectTimeout, readWriteTimeout time.Duration) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.SetTimeout(connectTimeout, readWriteTimeout)
	}
}

// WithHeader adds header item string in request.
func WithHeader(key, value string) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.Header(key, value)
	}
}

// WithCookie adds a cookie to the request.
func WithCookie(cookie *http.Cookie) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.Header("Cookie", cookie.String())
	}
}

// Withtokenfactory adds a custom function to set Authorization
func WithTokenFactory(tokenFactory func() string) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		t := tokenFactory()

		request.Header("Authorization", t)
	}
}

// WithBasicAuth adds a custom function to set basic auth
func WithBasicAuth(basicAuth func() (string, string)) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		username, password := basicAuth()
		request.SetBasicAuth(username, password)
	}
}

// WithFilters will use the filter as the invocation filters
func WithFilters(fcs ...FilterChain) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.SetFilters(fcs...)
	}
}

// WithContentType adds ContentType in header
func WithContentType(contentType string) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.Header(contentTypeKey, contentType)
	}
}

// WithParam adds query param in to request.
func WithParam(key, value string) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.Param(key, value)
	}
}

// WithRetry set retry times and delay for the request
// default is 0 (never retry)
// -1 retry indefinitely (forever)
// Other numbers specify the exact retry amount
func WithRetry(times int, delay time.Duration) BhojpurHTTPRequestOption {
	return func(request *BhojpurHTTPRequest) {
		request.Retries(times)
		request.RetryDelay(delay)
	}
}
