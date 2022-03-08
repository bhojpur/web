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
	webui "github.com/bhojpur/web/pkg/app/ui"
)

const (
	DefaultTitle       = "Bhojpur Web - Developer's Sandbox"
	DefaultDescription = "A landing page of Bhojpur Web application used by software developers"
	BackgroundColor    = "#2e343a"
)

type homePage struct {
	webapp.Compo
}

func NewHomePage() *homePage {
	return &homePage{}
}

func (p *homePage) OnPreRender(ctx webapp.Context) {
	p.initPage(ctx)
}

func (p *homePage) OnNav(ctx webapp.Context) {
	p.initPage(ctx)
}

func (p *homePage) initPage(ctx webapp.Context) {
	ctx.Page().SetTitle(DefaultTitle)
	ctx.Page().SetDescription(DefaultDescription)
	analytics.Page("home", nil)
}

func (p *homePage) Render() webapp.UI {
	return newPage().
		Title("Bhojpur Web").
		Icon("https://static.bhojpur.net/favicon.ico").
		Index(
			newIndexLink().Title("What is Bhojpur Web?"),
			newIndexLink().Title("Updates"),
			newIndexLink().Title("Declarative Syntax"),
			newIndexLink().Title("Standard HTTP Server"),
			newIndexLink().Title("Other features"),
			newIndexLink().Title("Built with Bhojpur Web"),

			webapp.Div().Class("separator"),

			newIndexLink().Title("Next"),
		).
		Content(
			webui.Flow().
				StretchItems().
				Spacing(84).
				Content(
					newRemoteMarkdownDoc().
						Class("fill").
						Src("/web/documents/what-is-bhojpur-web.md"),
					newRemoteMarkdownDoc().
						Class("fill").
						Class("updates").
						Src("/web/documents/updates.md"),
				),

			webapp.Div().Class("separator"),

			newRemoteMarkdownDoc().Src("/web/documents/home.md"),

			webapp.Div().Class("separator"),

			newBuiltWithBhojpur().ID("built-with-bhojpur"),

			webapp.Div().Class("separator"),

			newRemoteMarkdownDoc().Src("/web/documents/home-next.md"),
		)
}
