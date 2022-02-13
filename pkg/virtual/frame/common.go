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
	"encoding/binary"
	"fmt"
	"io"
)

var (
	// the byte order of all serialized integers
	order = binary.BigEndian
)

const (
	// masks for packing/unpacking frames
	streamMask = 0x7FFFFFFF
	typeMask   = 0xF0
	flagsMask  = 0x0F
	wndIncMask = 0x7FFFFFFF
	lengthMask = 0x00FFFFFF
)

// StreamId is 31-bit integer uniquely identifying a stream within a session
type StreamId uint32

func (id StreamId) valid() error {
	if id > streamMask {
		return fmt.Errorf("invalid stream id: %d", id)
	}
	return nil
}

// ErrorCode is a 32-bit integer indicating an error condition on a stream or session
type ErrorCode uint32

// Type is a 4-bit integer in the frame header that identifies the type of frame
type Type uint8

const (
	TypeRst    Type = 0x0
	TypeData   Type = 0x1
	TypeWndInc Type = 0x2
	TypeGoAway Type = 0x3
)

func (t Type) String() string {
	switch t {
	case TypeRst:
		return "RST"
	case TypeData:
		return "DATA"
	case TypeWndInc:
		return "WNDINC"
	case TypeGoAway:
		return "GOAWAY"
	}
	return "UNKNOWN"
}

// Flags is a 4-bit integer containing frame-specific flag bits in the frame header
type Flags uint8

const (
	FlagDataFin = 0x1
	FlagDataSyn = 0x2
)

func (f Flags) IsSet(g Flags) bool {
	return (f & g) != 0
}

func (f *Flags) Set(g Flags) {
	*f |= g
}

func (f *Flags) Unset(g Flags) {
	*f = *f &^ g
}

const (
	headerSize       = 8
	maxFixedBodySize = 8 // goaway frame has streamid + errorcode
	maxBufferSize    = headerSize + maxFixedBodySize
)

type common struct {
	streamId StreamId
	length   uint32
	ftype    Type
	flags    Flags
	b        [maxBufferSize]byte
}

func (f *common) StreamId() StreamId {
	return f.streamId
}

func (f *common) Length() uint32 {
	return f.length
}

func (f *common) Type() Type {
	return f.ftype
}

func (f *common) Flags() Flags {
	return f.flags
}

func (f *common) readFrom(r io.Reader) error {
	b := f.b[:headerSize]
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	f.length = (uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
	f.ftype = Type(b[3] >> 4)
	f.flags = Flags(b[3] & flagsMask)
	f.streamId = StreamId(order.Uint32(b[4:]))
	return nil
}

func (f *common) writeTo(w io.Writer, fixedSize int) error {
	_, err := w.Write(f.b[:headerSize+fixedSize])
	return err
}

func (f *common) pack(ftype Type, length int, streamId StreamId, flags Flags) error {
	if err := streamId.valid(); err != nil {
		return err
	}
	if !isValidLength(length) {
		return fmt.Errorf("invalid length: %d", length)
	}
	f.ftype = ftype
	f.streamId = streamId
	f.length = uint32(length)
	f.flags = flags
	_ = append(f.b[:0],
		byte(f.length>>16),
		byte(f.length>>8),
		byte(f.length),
		byte(uint8(f.ftype<<4)|uint8(f.flags&flagsMask)),
		byte(f.streamId>>24),
		byte(f.streamId>>16),
		byte(f.streamId>>8),
		byte(f.streamId),
	)
	return nil
}

func (f *common) body() []byte {
	return f.b[headerSize:]
}

func (f *common) String() string {
	s := fmt.Sprintf(
		"FRAME [TYPE: %s | LENGTH: %d | STREAMID: %x | FLAGS: %d",
		f.Type(), f.Length(), f.StreamId(), f.Flags())
	if f.Type() != TypeData && f.Type() != TypeGoAway {
		s += fmt.Sprintf(" | BODY: %x", f.body()[:f.Length()])
	}
	s += "]"
	return s
}

func isValidLength(length int) bool {
	return length >= 0 && length <= lengthMask
}
