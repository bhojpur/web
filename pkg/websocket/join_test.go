package websocket

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
	"io"
	"strings"
	"testing"
)

func TestJoinMessages(t *testing.T) {
	messages := []string{"a", "bc", "def", "ghij", "klmno", "0", "12", "345", "6789"}
	for _, readChunk := range []int{1, 2, 3, 4, 5, 6, 7} {
		for _, term := range []string{"", ","} {
			var connBuf bytes.Buffer
			wc := newTestConn(nil, &connBuf, true)
			rc := newTestConn(&connBuf, nil, false)
			for _, m := range messages {
				wc.WriteMessage(BinaryMessage, []byte(m))
			}

			var result bytes.Buffer
			_, err := io.CopyBuffer(&result, JoinMessages(rc, term), make([]byte, readChunk))
			if IsUnexpectedCloseError(err, CloseAbnormalClosure) {
				t.Errorf("readChunk=%d, term=%q: unexpected error %v", readChunk, term, err)
			}
			want := strings.Join(messages, term) + term
			if result.String() != want {
				t.Errorf("readChunk=%d, term=%q, got %q, want %q", readChunk, term, result.String(), want)
			}
		}
	}
}
