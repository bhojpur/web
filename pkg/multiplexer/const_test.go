package multiplexer

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
	"testing"
)

func TestConst(t *testing.T) {
	if protoVersion != 0 {
		t.Fatalf("bad: %v", protoVersion)
	}

	if typeData != 0 {
		t.Fatalf("bad: %v", typeData)
	}
	if typeWindowUpdate != 1 {
		t.Fatalf("bad: %v", typeWindowUpdate)
	}
	if typePing != 2 {
		t.Fatalf("bad: %v", typePing)
	}
	if typeGoAway != 3 {
		t.Fatalf("bad: %v", typeGoAway)
	}

	if flagSYN != 1 {
		t.Fatalf("bad: %v", flagSYN)
	}
	if flagACK != 2 {
		t.Fatalf("bad: %v", flagACK)
	}
	if flagFIN != 4 {
		t.Fatalf("bad: %v", flagFIN)
	}
	if flagRST != 8 {
		t.Fatalf("bad: %v", flagRST)
	}

	if goAwayNormal != 0 {
		t.Fatalf("bad: %v", goAwayNormal)
	}
	if goAwayProtoErr != 1 {
		t.Fatalf("bad: %v", goAwayProtoErr)
	}
	if goAwayInternalErr != 2 {
		t.Fatalf("bad: %v", goAwayInternalErr)
	}

	if headerSize != 12 {
		t.Fatalf("bad header size")
	}
}

func TestEncodeDecode(t *testing.T) {
	hdr := header(make([]byte, headerSize))
	hdr.encode(typeWindowUpdate, flagACK|flagRST, 1234, 4321)

	if hdr.Version() != protoVersion {
		t.Fatalf("bad: %v", hdr)
	}
	if hdr.MsgType() != typeWindowUpdate {
		t.Fatalf("bad: %v", hdr)
	}
	if hdr.Flags() != flagACK|flagRST {
		t.Fatalf("bad: %v", hdr)
	}
	if hdr.StreamID() != 1234 {
		t.Fatalf("bad: %v", hdr)
	}
	if hdr.Length() != 4321 {
		t.Fatalf("bad: %v", hdr)
	}
}
