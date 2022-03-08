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

func TestCondition(t *testing.T) {
	testUpdate(t, []updateTest{
		{
			scenario: "if is interpreted",
			a: Div().Body(
				If(false,
					H1(),
				),
			),
			b: Div().Body(
				If(true,
					H1(),
				),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H1(),
				},
			},
		},
		{
			scenario: "if is not interpreted",
			a: Div().Body(
				If(true,
					H1(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				),
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
			},
		},
		{
			scenario: "else if is interpreted",
			a: Div().Body(
				If(true,
					H1(),
				).ElseIf(false,
					H2(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H2(),
				},
			},
		},
		{
			scenario: "else if is not interpreted",
			a: Div().Body(
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(false,
					H2(),
				),
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
			},
		},
		{
			scenario: "else is interpreted",
			a: Div().Body(
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				).Else(
					H3(),
				),
			),
			b: Div().Body(
				If(false,
					H1(),
				).ElseIf(false,
					H2(),
				).Else(
					H3(),
				),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H3(),
				},
			},
		},
		{
			scenario: "else is not interpreted",
			a: Div().Body(
				If(false,
					H1(),
				).ElseIf(true,
					H2(),
				).Else(
					H3(),
				),
			),
			b: Div().Body(
				If(true,
					H1(),
				).ElseIf(false,
					H2(),
				).Else(
					H3(),
				),
			),
			matches: []TestUIDescriptor{
				{
					Path:     TestPath(),
					Expected: Div(),
				},

				{
					Path:     TestPath(0),
					Expected: H1(),
				},
			},
		},
	})
}
