package adapter

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
	"github.com/bhojpur/web/pkg/adapter/context"
	ctxsvr "github.com/bhojpur/web/pkg/context"

	web "github.com/bhojpur/web/pkg/engine"
)

// Tree has three elements: FixRouter/wildcard/leaves
// fixRouter stores Fixed Router
// wildcard stores params
// leaves store the endpoint information
type Tree web.Tree

// NewTree return a new Tree
func NewTree() *Tree {
	return (*Tree)(web.NewTree())
}

// AddTree will add tree to the exist Tree
// prefix should has no params
func (t *Tree) AddTree(prefix string, tree *Tree) {
	(*web.Tree)(t).AddTree(prefix, (*web.Tree)(tree))
}

// AddRouter call addseg function
func (t *Tree) AddRouter(pattern string, runObject interface{}) {
	(*web.Tree)(t).AddRouter(pattern, runObject)
}

// Match router to runObject & params
func (t *Tree) Match(pattern string, ctx *context.Context) (runObject interface{}) {
	return (*web.Tree)(t).Match(pattern, (*ctxsvr.Context)(ctx))
}
