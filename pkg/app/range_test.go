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

import "testing"

func TestRange(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "range slice is updated",
			a: Div().Body(
				Range([]string{"hello", "world"}).Slice(func(i int) UI {
					src := []string{"hello", "world"}
					return Text(src[i])
				}),
			),
			b: Div().Body(
				Range([]string{"hello", "maxoo"}).Slice(func(i int) UI {
					src := []string{"hello", "maxoo"}
					return Text(src[i])
				}),
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
				{
					Path:     TestPath(1),
					Expected: Text("maxoo"),
				},
			},
		},
		{
			scenario: "range slice is updated to be empty",
			a: Div().Body(
				Range([]string{"hello", "world"}).Slice(func(i int) UI {
					src := []string{"hello", "world"}
					return Text(src[i])
				}),
			),
			b: Div().Body(
				Range([]string{}).Slice(func(i int) UI {
					src := []string{"hello", "maxoo"}
					return Text(src[i])
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: nil,
				},
				{
					Path:     TestPath(1),
					Expected: nil,
				},
			},
		},
		{
			scenario: "range map is updated",
			a: Div().Body(
				Range(map[string]string{"key": "value"}).Map(func(k string) UI {
					src := map[string]string{"key": "value"}
					return Text(src[k])
				}),
			),
			b: Div().Body(
				Range(map[string]string{"key": "value"}).Map(func(k string) UI {
					src := map[string]string{"key": "maxoo"}
					return Text(src[k])
				}),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},
				{
					Path:     TestPath(0),
					Expected: Text("maxoo"),
				},
			},
		},
	})
}
