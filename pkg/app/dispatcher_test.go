package app

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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDispatcherMultipleMount(t *testing.T) {
	d := NewClientTester(Div())
	defer d.Close()
	d.Mount(A())
	d.Mount(Text("hello"))
	d.Mount(&hello{})
	d.Mount(&hello{})
	d.Consume()
}

func TestDispatcherAsyncWaitClient(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testDispatcherAsyncWait(t, d)
}

func TestDispatcherAsyncWaitServer(t *testing.T) {
	d := NewServerTester(&hello{})
	defer d.Close()
	testDispatcherAsyncWait(t, d)
}

func testDispatcherAsyncWait(t *testing.T, d Dispatcher) {
	var mu sync.Mutex
	var counts int

	inc := func() {
		mu.Lock()
		counts++
		mu.Unlock()
	}

	d.Async(inc)
	d.Async(inc)
	d.Async(inc)
	d.Async(inc)
	d.Async(inc)

	d.Wait()
	require.Equal(t, 5, counts)
}

func TestDispatcherLocalStorage(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testBrowserStorage(t, d.localStorage())
}

func TestDispatcherSessionStorage(t *testing.T) {
	d := NewClientTester(&hello{})
	defer d.Close()
	testBrowserStorage(t, d.sessionStorage())
}
