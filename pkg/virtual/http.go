package virtual

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
)

type HTTPConn struct {
	*sharedConn
	Request *http.Request
}

// HTTP parses the header of the first HTTP request on conn and returns
// a new, unread connection with metadata for virtual host multiplexing
func HTTP(conn net.Conn) (httpConn *HTTPConn, err error) {
	c, rd := newShared(conn)

	httpConn = &HTTPConn{sharedConn: c}
	if httpConn.Request, err = http.ReadRequest(bufio.NewReader(rd)); err != nil {
		return
	}

	// You probably don't need access to the request body and this makes the API
	// simpler by allowing you to call Free() optionally
	httpConn.Request.Body.Close()

	return
}

// Free sets Request to nil so that it can be garbage collected
func (c *HTTPConn) Free() {
	c.Request = nil
}

func (c *HTTPConn) Host() string {
	if c.Request == nil {
		return ""
	}

	return c.Request.Host
}
