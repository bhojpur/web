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
	"encoding/binary"
	"net"
)

var order = binary.BigEndian

type StreamType uint32

type TypedStream interface {
	Stream
	StreamType() StreamType
}

type TypedStreamSession interface {
	Session
	OpenTypedStream(stype StreamType) (Stream, error)
	AcceptTypedStream() (TypedStream, error)
}

func NewTypedStreamSession(s Session) TypedStreamSession {
	return &typedStreamSession{s}
}

type typedStreamSession struct {
	Session
}

func (s *typedStreamSession) Accept() (net.Conn, error) {
	return s.AcceptStream()
}

func (s *typedStreamSession) AcceptStream() (Stream, error) {
	return s.AcceptTypedStream()
}

func (s *typedStreamSession) AcceptTypedStream() (TypedStream, error) {
	str, err := s.Session.AcceptStream()
	if err != nil {
		return nil, err
	}
	var stype [4]byte
	_, err = str.Read(stype[:])
	if err != nil {
		str.Close()
		return nil, err
	}
	return &typedStream{str, StreamType(order.Uint32(stype[:]))}, nil
}

func (s *typedStreamSession) OpenTypedStream(st StreamType) (Stream, error) {
	str, err := s.Session.OpenStream()
	if err != nil {
		return nil, err
	}
	var stype [4]byte
	order.PutUint32(stype[:], uint32(st))
	_, err = str.Write(stype[:])
	if err != nil {
		return nil, err
	}
	return &typedStream{str, st}, nil
}

type typedStream struct {
	Stream
	streamType StreamType
}

func (s *typedStream) StreamType() StreamType {
	return s.streamType
}
