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
	"net"
	"net/http"
	"testing"
)

func TestHTTPHost(t *testing.T) {
	var testHostname string = "test.bhojpur.net"

	primary, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}
	defer primary.Close()

	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:12345")
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		req, err := http.NewRequest("GET", "http://"+testHostname+"/bar", nil)
		if err != nil {
			panic(err)
		}
		if err = req.Write(conn); err != nil {
			panic(err)
		}
	}()

	conn, err := primary.Accept()
	if err != nil {
		panic(err)
	}
	c, err := HTTP(conn)
	if err != nil {
		panic(err)
	}

	if c.Host() != testHostname {
		t.Errorf("Connection Host() is %s, expected %s", c.Host(), testHostname)
	}
}
