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
	"sync"
	"text/tabwriter"
)

type Frame interface {
	StreamId() StreamId
	Type() Type
	Flags() Flags
	Length() uint32
	readFrom(io.Reader) error
	writeTo(io.Writer) error
}

// A Framer serializes/deserializer frames to/from an io.ReadWriter
type Framer interface {
	// WriteFrame writes the given frame to the underlying transport
	WriteFrame(Frame) error

	// ReadFrame reads the next frame from the underlying transport
	ReadFrame() (Frame, error)
}

type framer struct {
	io.Reader
	io.Writer
	common

	// frames
	Rst
	Data
	WndInc
	GoAway
	Unknown
}

func (fr *framer) WriteFrame(f Frame) error {
	return f.writeTo(fr.Writer)
}

func (fr *framer) ReadFrame() (f Frame, err error) {
	if err := fr.common.readFrom(fr.Reader); err != nil {
		return nil, err
	}
	switch fr.common.ftype {
	case TypeRst:
		f = &fr.Rst
		fr.Rst.common = fr.common
	case TypeData:
		f = &fr.Data
		fr.Data.common = fr.common
	case TypeWndInc:
		f = &fr.WndInc
		fr.WndInc.common = fr.common
	case TypeGoAway:
		f = &fr.GoAway
		fr.GoAway.common = fr.common
	default:
		f = &fr.Unknown
		fr.Unknown.common = fr.common
	}
	return f, f.readFrom(fr)
}

func NewFramer(r io.Reader, w io.Writer) Framer {
	fr := &framer{
		Reader: r,
		Writer: w,
	}
	return fr
}

type debugFramer struct {
	debugWr *tabwriter.Writer
	once    sync.Once
	name    string
	Framer
}

func (fr *debugFramer) WriteFrame(f Frame) error {
	defer fr.debugWr.Flush()
	fr.printHeader()

	// actually write the frame to the real framer
	err := fr.Framer.WriteFrame(f)

	// each frame knows how to write iteself to the framer
	fmt.Fprintf(fr.debugWr, "%s\t%s\t%s\t0x%x\t%d\t0x%x\t%v\n", fr.name, "WRITE", f.Type(), f.StreamId(), f.Length(), f.Flags(), err)
	return err
}

func (fr *debugFramer) ReadFrame() (Frame, error) {
	defer fr.debugWr.Flush()
	fr.printHeader()
	f, err := fr.Framer.ReadFrame()
	if err == nil {
		fmt.Fprintf(fr.debugWr, "%s\t%s\t%s\t0x%x\t%d\t0x%x\t%v\n", fr.name, "READ", f.Type(), f.StreamId(), f.Length(), f.Flags(), nil)
	} else {
		fmt.Fprintf(fr.debugWr, "%s\t%s\t\t\t\t\t%v\n", fr.name, "READ", err)
	}
	return f, err
}

func (fr *debugFramer) printHeader() {
	fr.once.Do(func() {
		fmt.Fprintf(fr.debugWr, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "NAME", "OP", "TYPE", "STREAMID", "LENGTH", "FLAGS", "ERROR")
		fmt.Fprintf(fr.debugWr, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "----", "--", "----", "--------", "------", "-----", "-----")
	})
}

func NewDebugFramer(wr io.Writer, fr Framer) Framer {
	return NewNamedDebugFramer("", wr, fr)
}

func NewNamedDebugFramer(name string, wr io.Writer, fr Framer) Framer {
	return &debugFramer{
		Framer:  fr,
		debugWr: tabwriter.NewWriter(wr, 12, 2, 2, ' ', 0),
		name:    name,
	}
}
