package multiplexer

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
	"fmt"
	"net"
)

// hasAddr is used to get the address from the underlying connection
type hasAddr interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

// muxAddr is used when we cannot get the underlying address
type muxAddr struct {
	Addr string
}

func (*muxAddr) Network() string {
	return "bhojpur"
}

func (m *muxAddr) String() string {
	return fmt.Sprintf("bhojpur web multiplexer:%s", m.Addr)
}

// Addr is used to get the address of the listener.
func (s *Session) Addr() net.Addr {
	return s.LocalAddr()
}

// LocalAddr is used to get the local address of the
// underlying connection.
func (s *Session) LocalAddr() net.Addr {
	addr, ok := s.conn.(hasAddr)
	if !ok {
		return &muxAddr{"local"}
	}
	return addr.LocalAddr()
}

// RemoteAddr is used to get the address of remote end
// of the underlying connection
func (s *Session) RemoteAddr() net.Addr {
	addr, ok := s.conn.(hasAddr)
	if !ok {
		return &muxAddr{"remote"}
	}
	return addr.RemoteAddr()
}

// LocalAddr returns the local address
func (s *Stream) LocalAddr() net.Addr {
	return s.session.LocalAddr()
}

// RemoteAddr returns the remote address
func (s *Stream) RemoteAddr() net.Addr {
	return s.session.RemoteAddr()
}
