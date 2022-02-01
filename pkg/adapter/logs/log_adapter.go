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

import (
	"time"

	logs "github.com/bhojpur/logger/pkg/engine"
)

type oldToNewAdapter struct {
	old Logger
}

func (o *oldToNewAdapter) Init(config string) error {
	return o.old.Init(config)
}

func (o *oldToNewAdapter) WriteMsg(lm *logs.LogMsg) error {
	return o.old.WriteMsg(lm.When, lm.OldStyleFormat(), lm.Level)
}

func (o *oldToNewAdapter) Destroy() {
	o.old.Destroy()
}

func (o *oldToNewAdapter) Flush() {
	o.old.Flush()
}

func (o *oldToNewAdapter) SetFormatter(f logs.LogFormatter) {
	panic("unsupported operation, you should not invoke this method")
}

type newToOldAdapter struct {
	n logs.Logger
}

func (n *newToOldAdapter) Init(config string) error {
	return n.n.Init(config)
}

func (n *newToOldAdapter) WriteMsg(when time.Time, msg string, level int) error {
	return n.n.WriteMsg(&logs.LogMsg{
		When:  when,
		Msg:   msg,
		Level: level,
	})
}

func (n *newToOldAdapter) Destroy() {
	panic("implement me")
}

func (n *newToOldAdapter) Flush() {
	panic("implement me")
}
