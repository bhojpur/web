package frame

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
)

const (
	wndIncFrameLength = 4
)

// Increase a stream's flow control window size
type WndInc struct {
	common
}

func (f *WndInc) WindowIncrement() uint32 {
	return order.Uint32(f.body()) & wndIncMask
}

func (f *WndInc) readFrom(rd io.Reader) error {
	if f.length != wndIncFrameLength {
		return frameSizeError(f.length, "WNDINC")
	}
	if _, err := io.ReadFull(rd, f.body()[:wndIncFrameLength]); err != nil {
		return err
	}
	if f.StreamId() == 0 {
		return protoError("WndInc stream id must not be zero, got: %d", f.StreamId())
	}
	if f.WindowIncrement() == 0 {
		return protoStreamError("WndInc increment must not be zero, got: %d", f.WindowIncrement())
	}
	return nil
}

func (f *WndInc) writeTo(wr io.Writer) error {
	return f.common.writeTo(wr, wndIncFrameLength)
}

func (f *WndInc) Pack(streamId StreamId, inc uint32) (err error) {
	if inc > wndIncMask || inc == 0 {
		return fmt.Errorf("invalid window increment: %d", inc)
	}
	if err = f.common.pack(TypeWndInc, wndIncFrameLength, streamId, 0); err != nil {
		return
	}
	order.PutUint32(f.body(), inc)
	return
}
