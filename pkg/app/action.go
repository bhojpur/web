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
	"fmt"
	"sync"
)

// Action represents a custom event that can be propagated across the app. It
// can contain a payload and be given additional context with tags.
type Action struct {
	// The name that identifies the action..
	Name string

	// The value passed along with the action. Can be nil.
	Value interface{}

	// Tags that provide some context to the action.
	Tags Tags
}

// ActionHandler represents a handler that is executed when an action is created
// with Context.NewAction().
type ActionHandler func(Context, Action)

// Handle registers the handler for the given action name. When an action
// occurs, the handler is executed on its own goroutine.
func Handle(actionName string, h ActionHandler) {
	actionHandlers[actionName] = h
}

var actionHandlers = make(map[string]ActionHandler)

type actionHandler struct {
	async    bool
	source   UI
	function ActionHandler
}

type actionManager struct {
	once     sync.Once
	mutex    sync.Mutex
	handlers map[string]map[string]actionHandler
}

func (m *actionManager) init() {
	m.handlers = make(map[string]map[string]actionHandler)
}

func (m *actionManager) post(a Action) {
	m.once.Do(m.init)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	handlers := m.handlers[a.Name]
	for key, h := range handlers {
		source := h.source
		if !source.Mounted() {
			delete(handlers, key)
			continue
		}

		ctx := makeContext(source)
		function := h.function

		if h.async {
			ctx.Async(func() { function(ctx, a) })
		} else {
			ctx.Dispatch(func(ctx Context) { function(ctx, a) })
		}
	}
}

func (m *actionManager) handle(actionName string, async bool, source UI, h ActionHandler) {
	m.once.Do(m.init)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	handlers, isRegistered := m.handlers[actionName]
	if !isRegistered {
		handlers = make(map[string]actionHandler)
		m.handlers[actionName] = handlers
	}

	key := fmt.Sprintf("/%T:%p/%p", source, source, h)
	handlers[key] = actionHandler{
		async:    async,
		source:   source,
		function: h,
	}
}

func (m *actionManager) closeUnusedHandlers() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for actionName, handlers := range m.handlers {
		for key, h := range handlers {
			if !h.source.Mounted() {
				delete(handlers, key)
			}
		}

		if len(handlers) == 0 {
			delete(m.handlers, actionName)
		}
	}
}
