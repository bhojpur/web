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
	"bytes"
	"reflect"
	"testing"
)

type FrameTest interface {
	FrameName() string
	SerializeError() bool
	DeserializeError() bool
	Serialized() []byte
	Pack() (Frame, error)
	WithHeader(common) Frame
	Eq(Frame) error
}

func RunFrameTest(t *testing.T, ft FrameTest) {
	if !ft.DeserializeError() {
		runSerializeTest(t, ft, ft.SerializeError())
	}
	if !ft.SerializeError() {
		runDeserializeTest(t, ft, ft.DeserializeError())
	}
	if !ft.DeserializeError() && !ft.SerializeError() {
		runFramerTest(t, ft)
	}
}

func runSerializeTest(t *testing.T, ft FrameTest, expectError bool) {
	buf := new(bytes.Buffer)
	f, err := ft.Pack()
	switch {
	case err != nil && !expectError:
		t.Errorf("failed to pack %s frame: %v, %+v!", ft.FrameName(), err, ft)
		return
	case err == nil && expectError:
		t.Errorf("expected packing %s frame to error but it succeeded: %+v", ft.FrameName(), ft)
		return
	case expectError:
		return
	}
	if err := f.writeTo(buf); err != nil {
		t.Errorf("failed to write %s frame: %v, %+v!", ft.FrameName(), err, ft)
		return
	}
	if !reflect.DeepEqual(ft.Serialized(), buf.Bytes()) {
		t.Errorf("failed %s frame serialization, expected: %v got %v", ft.FrameName(), ft.Serialized(), buf.Bytes())
		return
	}
}

// test deserialization
func runDeserializeTest(t *testing.T, ft FrameTest, expectError bool) {
	buf := bytes.NewReader(ft.Serialized())
	var c common
	if err := c.readFrom(buf); err != nil {
		t.Errorf("failed read %s frame header: %v, %+v", ft.FrameName(), err, ft)
		return
	}
	f := ft.WithHeader(c)
	err := f.readFrom(buf)
	switch {
	case err != nil && !expectError:
		t.Errorf("failed to read %s frame: %v, %+v!", ft.FrameName(), err, ft)
		return
	case err == nil && expectError:
		t.Errorf("expected error while reading %s frame but got none: %+v!", ft.FrameName(), ft)
		return
	case expectError:
		return
	}

	// test for correctness
	if err := ft.Eq(f); err != nil {
		t.Errorf(err.Error())
	}
}

func runFramerTest(t *testing.T, ft FrameTest) {
	buf := new(bytes.Buffer)
	fr := NewFramer(buf, buf)

	f, err := ft.Pack()
	if err != nil {
		t.Errorf("failed to pack %s frame: %v, %+v!", ft.FrameName(), err, ft)
		return
	}
	err = fr.WriteFrame(f)
	if err != nil {
		t.Errorf("framer failed to write: %v", err)
		return
	}
	rf, err := fr.ReadFrame()
	if err != nil {
		t.Errorf("framer failed to read: %v", err)
		return
	}
	if err := ft.Eq(rf); err != nil {
		t.Errorf(err.Error())
	}
}
