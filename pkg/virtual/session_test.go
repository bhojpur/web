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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/bhojpur/web/pkg/virtual/frame"
)

func newFakeVirtualStream(sess sessionPrivate, id frame.StreamId, windowSize uint32, fin bool, init bool) streamPrivate {
	return &fakeVirtualStream{sess, id}
}

type fakeVirtualStream struct {
	sess     sessionPrivate
	streamId frame.StreamId
}

func (s *fakeVirtualStream) Write([]byte) (int, error)              { return 0, nil }
func (s *fakeVirtualStream) Read([]byte) (int, error)               { return 0, nil }
func (s *fakeVirtualStream) Close() error                           { return nil }
func (s *fakeVirtualStream) SetDeadline(time.Time) error            { return nil }
func (s *fakeVirtualStream) SetReadDeadline(time.Time) error        { return nil }
func (s *fakeVirtualStream) SetWriteDeadline(time.Time) error       { return nil }
func (s *fakeVirtualStream) CloseWrite() error                      { return nil }
func (s *fakeVirtualStream) Id() uint32                             { return uint32(s.streamId) }
func (s *fakeVirtualStream) Session() Session                       { return s.sess }
func (s *fakeVirtualStream) RemoteAddr() net.Addr                   { return nil }
func (s *fakeVirtualStream) LocalAddr() net.Addr                    { return nil }
func (s *fakeVirtualStream) handleStreamData(*frame.Data) error     { return nil }
func (s *fakeVirtualStream) handleStreamWndInc(*frame.WndInc) error { return nil }
func (s *fakeVirtualStream) handleStreamRst(*frame.Rst) error       { return nil }
func (s *fakeVirtualStream) closeWith(error)                        {}

type fakeVirtualConn struct {
	in     *io.PipeReader
	out    *io.PipeWriter
	closed bool
}

func (c *fakeVirtualConn) SetDeadline(time.Time) error      { return nil }
func (c *fakeVirtualConn) SetReadDeadline(time.Time) error  { return nil }
func (c *fakeVirtualConn) SetWriteDeadline(time.Time) error { return nil }
func (c *fakeVirtualConn) LocalAddr() net.Addr              { return nil }
func (c *fakeVirtualConn) RemoteAddr() net.Addr             { return nil }
func (c *fakeVirtualConn) Close() error                     { c.closed = true; c.in.Close(); return c.out.Close() }
func (c *fakeVirtualConn) Read(p []byte) (int, error)       { return c.in.Read(p) }
func (c *fakeVirtualConn) Write(p []byte) (int, error)      { return c.out.Write(p) }
func (c *fakeVirtualConn) Discard()                         { go io.Copy(ioutil.Discard, c.in) }

func newFakeVirtualConnPair() (local *fakeVirtualConn, remote *fakeVirtualConn) {
	local, remote = new(fakeVirtualConn), new(fakeVirtualConn)
	local.in, remote.out = io.Pipe()
	remote.in, local.out = io.Pipe()
	return
}

var debugFramer = func(name string) func(io.Reader, io.Writer) frame.Framer {
	return func(rd io.Reader, wr io.Writer) frame.Framer {
		return frame.NewNamedDebugFramer(name, os.Stdout, frame.NewFramer(rd, wr))
	}
}

func TestWrongClientParity(t *testing.T) {
	t.Parallel()
	local, remote := newFakeVirtualConnPair()
	// don't need the remote output
	remote.Discard()
	s := StreamServer(local, &Config{newStream: newFakeVirtualStream})

	// 300 is even, and only servers send even stream ids
	f := new(frame.Data)
	f.Pack(300, []byte{}, false, true)

	// send the frame into the session
	fr := frame.NewFramer(remote, remote)
	fr.WriteFrame(f)

	// wait for failure
	err, _, _ := s.Wait()

	if code, _ := GetError(err); code != ProtocolError {
		t.Errorf("Session not terminated with protocol error. Got %d, expected %d. Session error: %v", code, ProtocolError, err)
	}

	if !local.closed {
		t.Errorf("Session transport not closed after protocol failure.")
	}
}

func TestWrongServerParity(t *testing.T) {
	t.Parallel()

	local, remote := newFakeVirtualConnPair()
	s := StreamClient(local, &Config{newStream: newFakeVirtualStream})

	// don't need the remote output
	remote.Discard()

	// 301 is odd, and only clients send even stream ids
	f := new(frame.Data)
	f.Pack(301, []byte{}, false, true)

	// send the frame into the session
	fr := frame.NewFramer(remote, remote)
	fr.WriteFrame(f)

	// wait for failure
	err, _, _ := s.Wait()

	if code, _ := GetError(err); code != ProtocolError {
		t.Errorf("Session not terminated with protocol error. Got %d, expected %d. Session error: %v", code, ProtocolError, err)
	}

	if !local.closed {
		t.Errorf("Session transport not closed after protocol failure.")
	}
}

func TestAcceptStream(t *testing.T) {
	t.Parallel()

	local, remote := newFakeVirtualConnPair()

	// don't need the remote output
	remote.Discard()

	s := StreamClient(local, &Config{newStream: newFakeVirtualStream})
	defer s.Close()

	f := new(frame.Data)
	f.Pack(300, []byte{}, false, true)

	// send the frame into the session
	fr := frame.NewFramer(remote, remote)
	fr.WriteFrame(f)

	done := make(chan int)
	go func() {
		defer func() { done <- 1 }()

		// wait for accept
		str, err := s.AcceptStream()

		if err != nil {
			t.Errorf("Error accepting stream: %v", err)
			return
		}

		if str.Id() != 300 {
			t.Errorf("Stream has wrong id. Expected %d, got %d", str.Id(), 300)
		}
	}()

	select {
	case <-time.After(time.Second):
		t.Fatalf("Timed out!")
	case <-done:
	}
}

// validate that a session fulfills the net.Listener interface
// compile-only check
func TestNetListener(t *testing.T) {
	if false {
		var _ net.Listener = StreamServer(nil, nil)
	}
}

// Test for the Close() behavior
// Close() issues a data frame with the fin flag
// if any further data is received from the remote side, then RST is sent
func TestWriteAfterClose(t *testing.T) {
	t.Parallel()
	local, remote := newFakeVirtualConnPair()
	sLocal := StreamServer(local, &Config{NewFramer: debugFramer("SERVER")})
	sRemote := StreamClient(remote, &Config{NewFramer: debugFramer("CLIENT")})

	closed := make(chan int)
	go func() {
		stream, err := sRemote.Open()
		if err != nil {
			t.Errorf("Failed to open stream: %v", err)
			return
		}
		stream.Write([]byte("hello local"))
		defer sRemote.Close()

		<-closed
		// this write should succeed
		if _, err = stream.Write([]byte("test!")); err != nil {
			t.Errorf("Failed to write test data: %v", err)
			return
		}

		// give the remote end some time to send us an RST
		time.Sleep(200 * time.Millisecond)

		// this write should fail
		if _, err = stream.Write([]byte("test!")); err == nil {
			fmt.Println("WROTE FRAME WITHOUT ERROR")
			t.Errorf("expected error, but not did not receive one")
			return
		}
	}()

	stream, err := sLocal.Accept()
	if err != nil {
		t.Fatalf("Failed to accept stream!")
	}

	// tell the other side that we closed so they can write late
	stream.Close()
	closed <- 1

	err, remoteErr, debug := sLocal.Wait()
	if code, _ := GetError(err); code != PeerEOF {
		t.Fatalf("session closed with error: %v, expected PeerEOF", err)
	}
	remoteCode, _ := GetError(remoteErr)
	if remoteCode != NoError {
		t.Fatalf("remote session closed with error code: %v, expected NoError (debug: %s)", remoteCode, debug)
	}
}
