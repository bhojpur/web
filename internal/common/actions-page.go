package common

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
	webapp "github.com/bhojpur/web/pkg/app"
	analytics "github.com/bhojpur/web/pkg/app/analytics"
)

type actionPage struct {
	webapp.Compo
}

func newActionPage() *actionPage {
	return &actionPage{}
}

func (p *actionPage) OnPreRender(ctx webapp.Context) {
	p.initPage(ctx)
}

func (p *actionPage) OnNav(ctx webapp.Context) {
	p.initPage(ctx)
}

func (p *actionPage) initPage(ctx webapp.Context) {
	ctx.Page().SetTitle("Creating and Listening to Actions")
	ctx.Page().SetDescription("Documentation about how to create and listen to actions.")
	analytics.Page("actions", nil)
}

func (p *actionPage) Render() webapp.UI {
	return newPage().
		Title("Actions").
		Icon(actionSVG).
		Index(
			newIndexLink().Title("What is an Action?"),
			newIndexLink().Title("Create"),
			newIndexLink().Title("Handling"),
			newIndexLink().Title("    Global Level"),
			newIndexLink().Title("    Component Level"),

			webapp.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			newRemoteMarkdownDoc().Src("/web/documents/actions.md"),
		)
}
