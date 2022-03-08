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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	Handle("/test", func(Context, Action) {})
	require.Len(t, actionHandlers, 1)
}

func TestActionManagerHandle(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	m := actionManager{}

	h := &hello{}
	e.Mount(h)
	e.Consume()

	isHandleACalled := false
	isHandleBCalled := false
	isHandleCCalled := false
	isHandleDCalled := false

	m.handle("/test", false, h, func(ctx Context, a Action) {
		isHandleACalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	m.handle("/test", false, h, func(ctx Context, a Action) {
		isHandleBCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 2)

	f := &foo{}
	m.handle("/test", false, f, func(ctx Context, a Action) {
		isHandleCCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 3)

	m.handle("/test", true, e.Body, func(ctx Context, a Action) {
		isHandleDCalled = true
	})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 4)

	m.post(Action{Name: "/test"})
	e.Consume()
	require.True(t, isHandleACalled)
	require.True(t, isHandleBCalled)
	require.False(t, isHandleCCalled)
	require.True(t, isHandleDCalled)
	require.Len(t, m.handlers["/test"], 3)
}

func TestActionManagerCloseUnusedHandlers(t *testing.T) {
	e := engine{}
	e.init()
	defer e.Close()

	m := actionManager{}

	h := &hello{}
	e.Mount(h)
	e.Consume()

	m.handle("/test", false, h, func(ctx Context, a Action) {})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	f := &foo{}
	m.handle("/test", false, f, func(ctx Context, a Action) {})
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 2)

	m.closeUnusedHandlers()
	require.Len(t, m.handlers, 1)
	require.Len(t, m.handlers["/test"], 1)

	e.Mount(Div())
	e.Consume()
	m.closeUnusedHandlers()
	require.Empty(t, m.handlers)
}
