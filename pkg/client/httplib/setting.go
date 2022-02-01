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
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

// BhojpurHTTPSettings is the http.Client setting
type BhojpurHTTPSettings struct {
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	Retries          int // if set to -1 means will retry forever
	RetryDelay       time.Duration
	FilterChains     []FilterChain
	EscapeHTML       bool // if set to false means will not escape escape HTML special characters during processing, default true
}

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// SetDefaultSetting overwrites default settings
// Keep in mind that when you invoke the SetDefaultSetting
// some methods invoked before SetDefaultSetting
func SetDefaultSetting(setting BhojpurHTTPSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultSetting = setting
}

// GetDefaultSetting return current default setting
func GetDefaultSetting() BhojpurHTTPSettings {
	return defaultSetting
}

var defaultSetting = BhojpurHTTPSettings{
	UserAgent:        "Bhojpur Webctl",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	FilterChains:     make([]FilterChain, 0, 4),
	EscapeHTML:       true,
}

var (
	defaultCookieJar http.CookieJar
	settingMutex     sync.Mutex
)

// AddDefaultFilter add a new filter into defaultSetting
// Be careful about using this method if you invoke SetDefaultSetting somewhere
func AddDefaultFilter(fc FilterChain) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	if defaultSetting.FilterChains == nil {
		defaultSetting.FilterChains = make([]FilterChain, 0, 4)
	}
	defaultSetting.FilterChains = append(defaultSetting.FilterChains, fc)
}
