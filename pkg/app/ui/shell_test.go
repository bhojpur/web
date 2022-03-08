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

import (
	"testing"

	"github.com/bhojpur/web/pkg/app"
)

func TestShellPreRender(t *testing.T) {
	utests := []struct {
		scenario string
		shell    app.UI
	}{
		{
			scenario: "empty shell",
			shell:    Shell(),
		},
		{
			scenario: "shell with content",
			shell:    Shell().Content(app.Div()),
		},
		{
			scenario: "shell with menu",
			shell:    Shell().Menu(app.Div()),
		},
		{
			scenario: "shell with submenu",
			shell:    Shell().Index(app.Div()),
		},
		{
			scenario: "shell with overlay",
			shell:    Shell().HamburgerMenu(app.Div()),
		},
		{
			scenario: "shell with menu and content",
			shell: Shell().
				Menu(app.Div()).
				Content(app.Div()),
		},
		{
			scenario: "shell with menu, submenu and content",
			shell: Shell().
				Menu(app.Div()).
				Index(app.Div()).
				Content(app.Div()),
		},
		{
			scenario: "shell with menu, submenu, overlay menu and content",
			shell: Shell().
				Menu(app.Div()).
				Index(app.Div()).
				HamburgerMenu(app.Div()).
				Content(app.Div()),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			d := app.NewServerTester(u.shell)
			defer d.Close()
			d.PreRender()
		})
	}
}
