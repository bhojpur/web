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
	"io"
	"io/ioutil"
	"testing"

	"github.com/bhojpur/web/pkg/virtual/frame"
)

func TestCloseWrite(t *testing.T) {
	t.Parallel()

	local, remote := newFakeVirtualConnPair()

	s := StreamClient(local, &Config{NewFramer: debugFramer("CLIENT")})

	done := make(chan int)
	go func() {
		fr := debugFramer("SERVER")(remote, remote)
		readFrame := func(ftype frame.Type) frame.Frame {
			f, err := fr.ReadFrame()
			if err != nil {
				t.Errorf("Failed to read next frame: %v", err)
				return nil
			}
			if f.Type() != ftype {
				t.Errorf("Wrong frame type. Got %v, expected %v", f.Type(), ftype)
				return nil
			}
			return f
		}

		f := readFrame(frame.TypeData)
		if f.Length() != 10 {
			t.Errorf("Wrong data length. Got %d, expected %d", f.Length(), 10)
			return
		}
		n, _ := io.Copy(ioutil.Discard, f.(*frame.Data).Reader())
		if n != 10 {
			t.Errorf("Wrong read size. Got %d, expected %d", n, 10)
			return
		}

		f = readFrame(frame.TypeData)
		if f.Length() != 0 {
			t.Errorf("Wrong data length. Got %d, expected %d", f.Length(), 0)
			return
		}
		if !f.(*frame.Data).Fin() {
			t.Errorf("Expected fin flag on final data frame.")
			return
		}

		// io.Pipe is weird and apparently sometimes doesn't acknowledge a write completion
		// until the next read
		remote.Read([]byte{})

		close(done)
	}()

	str, err := s.OpenStream()
	if err != nil {
		t.Fatalf("Failed to open stream: %v", err)
	}

	_, err = str.Write(make([]byte, 10))
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	err = str.CloseWrite()
	if err != nil {
		t.Fatalf("Failed to close write")
	}

	<-done
}

/*
func TestDataAfterRst(t *testing.T) {
	local, remote := newFakeVirtualConnPair()

	_ = NewSession(local, NewStream, false, []Extension{})
	trans := frame.NewBasicTransport(remote)

	// make sure that we get an RST STREAM_CLOSED
	done := make(chan int)
	go func() {
		defer func() { done <- 1 }()

		f, err := trans.ReadFrame()
		if err != nil {
			t.Errorf("Failed to read frame sent from session: %v", err)
			return
		}

		fr, ok := f.(*frame.RStreamRst)
		if !ok {
			t.Errorf("Frame is not STREAM_RST: %v", f)
			return
		}

		if fr.ErrorCode() != frame.StreamClosed {
			t.Errorf("Error code on STREAM_RST is not STREAM_CLOSED. Got %d, expected %d", fr.ErrorCode(), frame.StreamClosed)
			return
		}
	}()

	fSyn := frame.NewWStreamSyn()
	if err := fSyn.Set(301, 0, 0, false); err != nil {
		t.Fatalf("Failed to make syn frame: %v", err)
	}

	if err := trans.WriteFrame(fSyn); err != nil {
		t.Fatalf("Failed to send syn: %v", err)
	}

	fRst := frame.NewWStreamRst()
	if err := fRst.Set(301, frame.Cancel); err != nil {
		t.Fatal("Failed to make rst frame: %v", err)
	}

	if err := trans.WriteFrame(fRst); err != nil {
		t.Fatalf("Failed to write rst frame: %v", err)
	}

	fData := frame.NewWStreamData()
	if err := fData.Set(301, []byte{0xa, 0xFF}, false); err != nil {
		t.Fatalf("Failed to set data frame")
	}

	trans.WriteFrame(fData)

	<-done
}

func TestFlowControlError(t *testing.T) {
	local, remote := newFakeVirtualConnPair()

	s := NewSession(local, NewStream, false, []Extension{})
	s.(*Session).defaultWindowSize = 10
	trans := frame.NewBasicTransport(remote)

	// make sure that we get an RST FLOW_CONTROL_ERROR
	done := make(chan int)
	go func() {
		defer func() { done <- 1 }()

		f, err := trans.ReadFrame()
		if err != nil {
			t.Errorf("Failed to read frame sent from session: %v", err)
			return
		}

		fr, ok := f.(*frame.RStreamRst)
		if !ok {
			t.Errorf("Frame is not STREAM_RST: %v", f)
			return
		}

		if fr.ErrorCode() != frame.FlowControlError {
			t.Errorf("Error code on STREAM_RST is not FLOW_CONTROL_ERROR. Got %d, expected %d", fr.ErrorCode(), frame.FlowControlError)
			return
		}
	}()

	fSyn := frame.NewWStreamSyn()
	if err := fSyn.Set(301, 0, 0, false); err != nil {
		t.Fatalf("Failed to make syn frame: %v", err)
	}

	if err := trans.WriteFrame(fSyn); err != nil {
		t.Fatalf("Failed to send syn: %v", err)
	}

	fData := frame.NewWStreamData()
	if err := fData.Set(301, make([]byte, 11), false); err != nil {
		t.Fatalf("Failed to set data frame")
	}

	trans.WriteFrame(fData)

	<-done
}

func TestTolerateLateFrameAfterRst(t *testing.T) {
	local, remote := newFakeVirtualConnPair()

	s := NewSession(local, NewStream, false, []Extension{})
	trans := frame.NewBasicTransport(remote)

	// make sure that we don't get any error on a late frame
	done := make(chan int)
	go func() {
		defer func() { done <- 1 }()
		// read syn
		trans.ReadFrame()
		// read rst
		trans.ReadFrame()

		// should block
		if f, err := trans.ReadFrame(); err != nil {
			t.Errorf("Error reading frame: %v", err)
		} else {
			t.Errorf("Got frame that we shouldn't have read: %v. Type: %v", f, f.Type())
		}
	}()

	str, err := s.Open()
	if err != nil {
		t.Fatalf("failed to open stream")
	}

	str.(*Stream).resetWith(frame.Cancel, fmt.Errorf("cancel"))

	fData := frame.NewWStreamData()
	if err := fData.Set(str.Id(), []byte{0x1, 0x2, 0x3}, false); err != nil {
		t.Fatalf("Failed to set data frame")
	}
	trans.WriteFrame(fData)

	select {
	case <-done:
		t.Fatalf("Stream sent response to late DATA frame")

	case <-time.After(1 * time.Second):
		// ok
	}
}

// Test that we remove a stream from the session if both sides half-close
func TestRemoveAfterHalfClose(t *testing.T) {
	local, remote := newFakeVirtualConnPair()
	remote.Discard()

	s := NewSession(local, NewStream, false, []Extension{})
	trans := frame.NewBasicTransport(remote)

	// open stream
	str, err := s.Open()
	if err != nil {
		t.Fatalf("failed to open stream")
	}

	// half close remote side (true means half-close)
	fData := frame.NewWStreamData()
	if err := fData.Set(str.Id(), []byte{0x1, 0x2, 0x3}, true); err != nil {
		t.Fatalf("Failed to set data frame")
	}

	if err := trans.WriteFrame(fData); err != nil {
		t.Fatalf("Failed to write data frame")
	}

	// half-close local side
	str.HalfClose([]byte{0xFF, 0xFE, 0xFD, 0xFC})

	// yield so the stream can process
	time.Sleep(0)

	// verify stream is removed
	if stream, ok := s.(*Session).streams.Get(str.Id()); ok {
		t.Fatalf("Expected stream %d to be removed after both sides half-closed, but found: %v!", str.Id(), stream)

	}
}

// Test that we get a RST if we send a DATA frame after we send a DATA frame with a FIN
func TestDataAfterFin(t *testing.T) {
}
*/
