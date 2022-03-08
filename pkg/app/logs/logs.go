package logs

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

// It implements functions to manipulate logs.
//
// Logs created are taggable.
//
//   logWithTags := logs.New("a log with tags").
//       Tag("a", 42).
// 	     Tag("b", 21)

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unsafe"
)

// New returns a log with the given description that can be tagged.
func New(v string) Log {
	return Log{
		description: v,
	}
}

// Newf returns a log with the given formatted description that can be tagged.
func Newf(format string, v ...interface{}) Log {
	return New(fmt.Sprintf(format, v...))
}

// Log is a implementation that supports tagging.
type Log struct {
	description string
	tags        []tag
	maxKeyLen   int
}

// Tag sets the named tag with the given value.
func (l Log) Tag(k string, v interface{}) Log {
	if l.tags == nil {
		l.tags = make([]tag, 0, 8)
	}

	if length := len(k); length > l.maxKeyLen {
		l.maxKeyLen = length
	}

	switch v := v.(type) {
	case string:
		l.tags = append(l.tags, tag{key: k, value: v})

	default:
		l.tags = append(l.tags, tag{key: k, value: fmt.Sprintf("%+v", v)})
	}

	return l
}

func (l Log) String() string {
	w := bytes.NewBuffer(make([]byte, 0, len(l.description)+len(l.tags)*(l.maxKeyLen+11)))
	l.format(w, 0)
	return bytesToString(w.Bytes())
}

func (l Log) format(w *bytes.Buffer, indent int) {
	w.WriteString(l.description)
	if len(l.tags) != 0 {
		w.WriteByte(':')
	}

	tags := l.tags
	sort.Slice(tags, func(a, b int) bool {
		return strings.Compare(tags[a].key, tags[b].key) < 0
	})

	for _, t := range l.tags {
		k := t.key
		v := t.value

		w.WriteByte('\n')
		l.indent(w, indent+4)
		w.WriteString(k)
		w.WriteByte(':')
		l.indent(w, l.maxKeyLen-len(k)+1)
		w.WriteString(v)
	}
}

func (l Log) indent(w *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		w.WriteByte(' ')
	}
}

type tag struct {
	key   string
	value string
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
