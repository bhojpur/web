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
	"sync"

	"github.com/bhojpur/web/pkg/virtual/frame"
)

const (
	initMapCapacity = 128 // not too much extra memory wasted to avoid allocations
)

// streamMap is a map of stream ids -> streams guarded by a read/write lock
type streamMap struct {
	sync.RWMutex
	table map[frame.StreamId]streamPrivate
}

func (m *streamMap) Get(id frame.StreamId) (s streamPrivate, ok bool) {
	m.RLock()
	s, ok = m.table[id]
	m.RUnlock()
	return
}

func (m *streamMap) Set(id frame.StreamId, str streamPrivate) {
	m.Lock()
	m.table[id] = str
	m.Unlock()
}

func (m *streamMap) Delete(id frame.StreamId) {
	m.Lock()
	delete(m.table, id)
	m.Unlock()
}

func (m *streamMap) Each(fn func(frame.StreamId, streamPrivate)) {
	m.RLock()
	streams := make(map[frame.StreamId]streamPrivate, len(m.table))
	for k, v := range m.table {
		streams[k] = v
	}
	m.RUnlock()

	for id, str := range streams {
		fn(id, str)
	}
}

func newStreamMap() *streamMap {
	return &streamMap{table: make(map[frame.StreamId]streamPrivate, initMapCapacity)}
}
