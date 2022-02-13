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

import "io"

const (
	rstFrameLength = 4
)

// Rst is a frame sent to forcibly close a stream
type Rst struct {
	common
}

func (f *Rst) ErrorCode() ErrorCode {
	return ErrorCode(order.Uint32(f.body()))
}

func (f *Rst) readFrom(rd io.Reader) (err error) {
	if f.length != rstFrameLength {
		return frameSizeError(f.length, "RST")
	}
	if _, err = io.ReadFull(rd, f.body()[:rstFrameLength]); err != nil {
		return err
	}
	if f.StreamId() == 0 {
		return protoError("RST stream id must not be zero")
	}
	return
}

func (f *Rst) writeTo(wr io.Writer) (err error) {
	return f.common.writeTo(wr, rstFrameLength)
}

func (f *Rst) Pack(streamId StreamId, errorCode ErrorCode) (err error) {
	if err = f.common.pack(TypeRst, rstFrameLength, streamId, 0); err != nil {
		return
	}
	order.PutUint32(f.body(), uint32(errorCode))
	return
}
