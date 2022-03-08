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

func TestRawRootTagName(t *testing.T) {
	tests := []struct {
		scenario string
		raw      string
		expected string
	}{
		{
			scenario: "tag set",
			raw: `
			<div>
				<span></span>
			</div>`,
			expected: "div",
		},
		{
			scenario: "tag is empty",
		},
		{
			scenario: "opening tag missing",
			raw:      "</div>",
		},
		{
			scenario: "tag is not set",
			raw:      "div",
		},
		{
			scenario: "tag is not closing",
			raw:      "<div",
		},
		{
			scenario: "tag is not closing",
			raw:      "<div",
		},
		{
			scenario: "tag without value",
			raw:      "<>",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			tag := rawRootTagName(test.raw)
			require.Equal(t, test.expected, tag)
		})
	}
}

func TestRawMountDismount(t *testing.T) {
	testMountDismount(t, []mountTest{
		{
			scenario: "raw html element",
			node:     Raw(`<h1>Hello</h1>`),
		},
		{
			scenario: "raw svg element",
			node:     Raw(`<svg></svg>`),
		},
	})
}

func TestRawUpdate(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario:   "raw html element returns replace error when updated with a non text-element",
			a:          Raw("<svg></svg>"),
			b:          Div(),
			replaceErr: true,
		},
		{
			scenario: "raw html element is replace by another raw html element",
			a: Div().Body(
				Raw("<div></div>"),
			),
			b: Div().Body(
				Raw("<svg></svg>"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Raw("<svg></svg>"),
				},
			},
		},
		{
			scenario: "raw html element is replace by non-raw html element",
			a: Div().Body(
				Raw("<div></div>"),
			),
			b: Div().Body(
				Text("hello"),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("hello"),
				},
			},
		},
	})
}
