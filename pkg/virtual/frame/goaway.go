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

const goAwayFrameLength = 8

type GoAway struct {
	common
	debugToWrite []byte
	debugToRead  io.LimitedReader
}

func (f *GoAway) LastStreamId() StreamId {
	return StreamId(order.Uint32(f.body()))
}

func (f *GoAway) ErrorCode() ErrorCode {
	return ErrorCode(order.Uint32(f.body()[4:]))
}

func (f *GoAway) Debug() io.Reader {
	return &f.debugToRead
}

func (f *GoAway) readFrom(rd io.Reader) error {
	if f.length < goAwayFrameLength {
		return frameSizeError(f.length, "GOAWAY")
	}
	if _, err := io.ReadFull(rd, f.body()[:goAwayFrameLength]); err != nil {
		return err
	}
	if f.StreamId() != 0 {
		return protoError("GOAWAY stream id must be zero, not: %d", f.StreamId())
	}
	f.debugToRead.R = rd
	f.debugToRead.N = int64(f.Length())
	return nil
}

func (f *GoAway) writeTo(wr io.Writer) (err error) {
	if err = f.common.writeTo(wr, goAwayFrameLength); err != nil {
		return
	}
	if _, err = wr.Write(f.debugToWrite); err != nil {
		return err
	}
	return
}

func (f *GoAway) Pack(lastStreamId StreamId, errCode ErrorCode, debug []byte) (err error) {
	if err = lastStreamId.valid(); err != nil {
		return
	}
	if err = f.common.pack(TypeGoAway, goAwayFrameLength+len(debug), 0, 0); err != nil {
		return
	}
	order.PutUint32(f.body(), uint32(lastStreamId))
	order.PutUint32(f.body()[4:], uint32(errCode))
	f.debugToWrite = debug
	return nil
}
