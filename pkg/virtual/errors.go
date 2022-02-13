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
	"errors"

	"github.com/bhojpur/web/pkg/virtual/frame"
)

// ErrorCode is a 32-bit integer indicating the type of an error condition
type ErrorCode uint32

const (
	NoError ErrorCode = iota
	ProtocolError
	InternalError
	FlowControlError
	StreamClosed
	StreamRefused
	StreamCancelled
	StreamReset
	FrameSizeError
	AcceptQueueFull
	EnhanceYourCalm
	RemoteGoneAway
	StreamsExhausted
	WriteTimeout
	SessionClosed
	PeerEOF

	ErrorUnknown ErrorCode = 0xFF
)

var (
	remoteGoneAway      = newErr(RemoteGoneAway, errors.New("remote gone away"))
	streamsExhausted    = newErr(StreamsExhausted, errors.New("streams exhuastated"))
	streamClosed        = newErr(StreamClosed, errors.New("stream closed"))
	writeTimeout        = newErr(WriteTimeout, errors.New("write timed out"))
	flowControlViolated = newErr(FlowControlError, errors.New("flow control violated"))
	sessionClosed       = newErr(SessionClosed, errors.New("session closed"))
	eofPeer             = newErr(PeerEOF, errors.New("read EOF from remote peer"))
)

func fromFrameError(err error) error {
	if e, ok := err.(*frame.Error); ok {
		switch e.Type() {
		case frame.ErrorFrameSize:
			return &virtualError{FrameSizeError, err}
		case frame.ErrorProtocol, frame.ErrorProtocolStream:
			return &virtualError{ProtocolError, err}
		}
	}
	return err
}

type virtualError struct {
	ErrorCode
	error
}

func (e *virtualError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}
	return "<nil>"
}

func newErr(code ErrorCode, err error) error {
	return &virtualError{code, err}
}

func GetError(err error) (ErrorCode, error) {
	if err == nil {
		return NoError, nil
	}
	if e, ok := err.(*virtualError); ok {
		return e.ErrorCode, e.error
	}
	return ErrorUnknown, err
}
