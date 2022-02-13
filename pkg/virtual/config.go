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
	"io"
	"sync"

	"github.com/bhojpur/web/pkg/virtual/frame"
)

var zeroConfig Config

type Config struct {
	// Maximum size of unread data to receive and buffer (per-stream). Default 256KB.
	MaxWindowSize uint32
	// Maximum number of inbound streams to queue for Accept(). Default 128.
	AcceptBacklog uint32
	// Function creating the Session's framer. Deafult frame.NewFramer()
	NewFramer func(io.Reader, io.Writer) frame.Framer

	// allow safe concurrent initialization
	initOnce sync.Once

	// Function to create new streams
	newStream streamFactory

	// Size of writeFrames channel
	writeFrameQueueDepth int
}

func (c *Config) initDefaults() {
	c.initOnce.Do(func() {
		if c.MaxWindowSize == 0 {
			c.MaxWindowSize = 0x40000 // 256KB
		}
		if c.AcceptBacklog == 0 {
			c.AcceptBacklog = 128
		}
		if c.NewFramer == nil {
			c.NewFramer = frame.NewFramer
		}
		if c.newStream == nil {
			c.newStream = newStream
		}
		if c.writeFrameQueueDepth == 0 {
			c.writeFrameQueueDepth = 64
		}
	})
}
