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
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

var (
	bufferFull   = errors.New("buffer is full")
	bufferClosed = errors.New("buffer closed previously")
)

type buffer interface {
	Read([]byte) (int, error)
	ReadFrom(io.Reader) (int, error)
	SetError(error)
	SetDeadline(time.Time)
}

type inboundBuffer struct {
	cond sync.Cond
	mu   sync.Mutex
	bytes.Buffer
	err     error
	maxSize int
}

func (b *inboundBuffer) Init(maxSize int) {
	b.cond.L = &b.mu
	b.maxSize = maxSize
}

func (b *inboundBuffer) ReadFrom(rd io.Reader) (n int, err error) {
	var n64 int64
	b.mu.Lock()
	if b.err != nil {
		if _, err = ioutil.ReadAll(rd); err == nil {
			err = bufferClosed
		}
		goto DONE
	}

	n64, err = b.Buffer.ReadFrom(rd)
	n += int(n64)
	if b.Buffer.Len() > b.maxSize {
		err = bufferFull
		b.err = bufferFull
	}

	b.cond.Broadcast()
DONE:
	b.mu.Unlock()
	return int(n), err
}

func (b *inboundBuffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	for {
		if b.Len() != 0 {
			n, err = b.Buffer.Read(p)
			break
		}
		if b.err != nil {
			err = b.err
			break
		}
		b.cond.Wait()
	}
	b.mu.Unlock()
	return
}

func (b *inboundBuffer) SetError(err error) {
	b.mu.Lock()
	b.err = err
	b.mu.Unlock()
	b.cond.Broadcast()
}

func (b *inboundBuffer) SetDeadline(t time.Time) {
	// XXX: implement
	/*
		b.L.Lock()

		// set the deadline
		b.deadline = t

		// how long until the deadline
		delay := t.Sub(time.Now())

		if b.timer != nil {
			b.timer.Stop()
		}

		// after the delay, wake up waiters
		b.timer = time.AfterFunc(delay, func() {
			b.Broadcast()
		})

		b.L.Unlock()
	*/
}
