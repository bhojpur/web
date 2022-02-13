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
	"bytes"
	"io"
	"net"
	"sync"
)

const (
	initVirtualHostBufSize = 1024 // allocate 1 KB up front to try to avoid resizing
)

type sharedConn struct {
	sync.Mutex
	net.Conn               // the raw connection
	vhostBuf *bytes.Buffer // all of the initial data that has to be read in order to virtual host a connection is saved here
}

func newShared(conn net.Conn) (*sharedConn, io.Reader) {
	c := &sharedConn{
		Conn:     conn,
		vhostBuf: bytes.NewBuffer(make([]byte, 0, initVirtualHostBufSize)),
	}

	return c, io.TeeReader(conn, c.vhostBuf)
}

func (c *sharedConn) Read(p []byte) (n int, err error) {
	c.Lock()
	if c.vhostBuf == nil {
		c.Unlock()
		return c.Conn.Read(p)
	}
	n, err = c.vhostBuf.Read(p)

	// end of the request buffer
	if err == io.EOF {
		// let the request buffer get garbage collected
		// and make sure we don't read from it again
		c.vhostBuf = nil

		// continue reading from the connection
		var n2 int
		n2, err = c.Conn.Read(p[n:])

		// update total read
		n += n2
	}
	c.Unlock()
	return
}
