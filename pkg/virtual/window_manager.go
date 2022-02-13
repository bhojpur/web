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
)

type windowManager interface {
	Increment(int)
	Decrement(int) (int, error)
	SetError(error)
}

type condWindow struct {
	val     int
	maxSize int
	err     error
	sync.Cond
	sync.Mutex
}

func newCondWindow(initialSize int) *condWindow {
	w := new(condWindow)
	w.Init(initialSize)
	return w
}

func (w *condWindow) Init(initialSize int) {
	w.val = initialSize
	w.maxSize = initialSize
	w.Cond.L = &w.Mutex
}

func (w *condWindow) Increment(inc int) {
	w.L.Lock()
	w.val += inc
	w.Broadcast()
	w.L.Unlock()
}

func (w *condWindow) SetError(err error) {
	w.L.Lock()
	w.err = err
	w.Broadcast()
	w.L.Unlock()
}

func (w *condWindow) Decrement(dec int) (ret int, err error) {
	if dec == 0 {
		return
	}

	w.L.Lock()
	for {
		if w.err != nil {
			err = w.err
			break
		}

		if w.val > 0 {
			if dec > w.val {
				ret = w.val
				w.val = 0
				break
			} else {
				ret = dec
				w.val -= dec
				break
			}
		} else {
			w.Wait()
		}
	}
	w.L.Unlock()
	return
}
