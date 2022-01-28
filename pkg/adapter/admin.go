package adapter

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
	"time"

	_ "github.com/bhojpur/web/pkg/core/admin"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// FilterMonitorFunc is default monitor filter when admin module is enable.
// if this func returns, admin module records qps for this request by condition of this function logic.
// usage:
// 	func MyFilterMonitor(method, requestPath string, t time.Duration, pattern string, statusCode int) bool {
//	 	if method == "POST" {
//			return false
//	 	}
//	 	if t.Nanoseconds() < 100 {
//			return false
//	 	}
//	 	if strings.HasPrefix(requestPath, "/bhojpur") {
//			return false
//	 	}
//	 	return true
// 	}
// 	bhojpur.FilterMonitorFunc = MyFilterMonitor.
var FilterMonitorFunc func(string, string, time.Duration, string, int) bool

// PrintTree prints all registered routers.
func PrintTree() M {
	return (M)(websvr.WebEngine.PrintTree())
}
