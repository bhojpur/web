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
	"context"
	"net/url"
)

const (
	dispatcherSize = 4096
)

// Dispatcher is the interface that describes an environment that synchronizes
// UI instructions and UI elements lifecycle.
type Dispatcher interface {
	// Context returns the context associated with the root element.
	Context() Context

	// Executes the given dispatch operation on the UI goroutine.
	Dispatch(d Dispatch)

	// Emit executes the given function and notifies the source's parent
	// components to update their state.
	Emit(src UI, fn func())

	// Handle registers the handler for the given action name. When an action
	// occurs, the handler is executed on the UI goroutine.
	Handle(actionName string, src UI, h ActionHandler)

	// Post posts the given action. The action is then handled by handlers
	// registered with Handle() and Context.Handle().
	Post(a Action)

	// Sets the state with the given value.
	SetState(state string, v interface{}, opts ...StateOption)

	// Stores the specified state value into the given receiver. Panics when the
	// receiver is not a pointer or nil.
	GetState(state string, recv interface{})

	// Deletes the given state.
	DelState(state string)

	// Creates an observer that observes changes for the specified state while
	// the given element is mounted.
	ObserveState(state string, elem UI) Observer

	// 	Async launches the given function on a new goroutine.
	//
	// The difference versus just launching a goroutine is that it ensures that
	// the asynchronous instructions are called before the dispatcher is closed.
	//
	// This is important during component prerendering since asynchronous
	// operations need to complete before sending a pre-rendered page over HTTP.
	Async(fn func())

	// Wait waits for the asynchronous operations launched with Async() to
	// complete.
	Wait()

	start(context.Context)
	currentPage() Page
	localStorage() BrowserStorage
	sessionStorage() BrowserStorage
	runsInServer() bool
	resolveStaticResource(string) string
	removeFromUpdates(Composer)
}

// ClientDispatcher is the interface that describes a dispatcher that emulates a
// client environment.
type ClientDispatcher interface {
	Dispatcher

	// Consume executes all the remaining UI instructions.
	Consume()

	// ConsumeNext executes the next UI instructions.
	ConsumeNext()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()

	// Mounts the given component as root element.
	Mount(UI)

	// Triggers OnNav from the root component.
	Nav(*url.URL)

	// Triggers OnAppUpdate from the root component.
	AppUpdate()

	// Triggers OnAppInstallChange from the root component.
	AppInstallChange()

	// Triggers OnAppResize from the root component.
	AppResize()
}

// NewClientTester creates a testing dispatcher that simulates a
// client environment. The given UI element is mounted upon creation.
func NewClientTester(n UI) ClientDispatcher {
	e := &engine{ActionHandlers: actionHandlers}
	e.init()
	e.Mount(n)
	e.Consume()
	return e
}

// ServerDispatcher is the interface that describes a dispatcher that emulates a server environment.
type ServerDispatcher interface {
	Dispatcher

	// Consume executes all the remaining UI instructions.
	Consume()

	// ConsumeNext executes the next UI instructions.
	ConsumeNext()

	// Close consumes all the remaining UI instruction and releases allocated
	// resources.
	Close()

	// Pre-renders the given component.
	PreRender()
}

// NewServerTester creates a testing dispatcher that simulates a
// client environment.
func NewServerTester(n UI) ServerDispatcher {
	e := &engine{
		RunsInServer:   true,
		ActionHandlers: actionHandlers,
	}
	e.init()
	e.Mount(n)
	e.Consume()
	return e
}

// Dispatch represents an operation executed on the UI goroutine.
type Dispatch struct {
	Mode     DispatchMode
	Source   UI
	Function func(Context)
}

// DispatchMode represents how a dispatch is processed.
type DispatchMode int

const (
	// A dispatch mode where the dispatched operation is enqueued to be executed
	// as soon as possible and its associated UI element is updated at the end
	// of the current update cycle.
	Update DispatchMode = iota

	// A dispatch mode that schedules the dispatched operation to be executed
	// after the current update frame.
	Defer

	// A dispatch mode where the dispatched operation is enqueued to be executed
	// as soon as possible.
	Next
)

// MsgHandler represents a handler to listen to messages sent with Context.Post.
type MsgHandler func(Context, interface{})
