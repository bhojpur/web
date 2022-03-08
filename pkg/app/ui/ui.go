package ui

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

// It provides a set of components to organize a web application's layout.

import (
	"strconv"
)

var (
	// The padding of block-like components in px.
	BlockPadding = 30

	// The padding of block-like components in px when app width is <= 480px.
	BlockMobilePadding = 18

	// The content width of block-like components in px.
	BlockContentWidth = 580

	// The horizontal padding of base-like components in px.
	BaseHPadding = 36

	// The horizontal padding of base-like components in px when app width is <= 480px.
	BaseMobileHPadding = 12

	// The horizontal padding of base-like ad components in px.
	BaseAdHPadding = BaseHPadding / 2

	// The vertical padding of base-like components in px.
	BaseVPadding = 12

	// The default icon size in px.
	DefaultIconSize = 24

	// The default icon space.
	DefaultIconSpace = 6

	// The default width for flow items in px.
	DefaultFlowItemWidth = 372
)

const (
	defaultHeaderHeight = 90
)

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}

type alignment int

const (
	stretch alignment = iota
	top
	right
	bottom
	left
	middle
)

type style struct {
	key   string
	value string
}
