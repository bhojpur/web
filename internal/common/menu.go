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
	webui "github.com/bhojpur/web/pkg/app/ui"
)

type menu struct {
	webapp.Compo

	Iclass string

	appInstallable bool
}

func newMenu() *menu {
	return &menu{}
}

func (m *menu) Class(v string) *menu {
	m.Iclass = webapp.AppendClass(m.Iclass, v)
	return m
}

func (m *menu) OnNav(ctx webapp.Context) {
	m.appInstallable = ctx.IsAppInstallable()
}

func (m *menu) OnAppInstallChange(ctx webapp.Context) {
	m.appInstallable = ctx.IsAppInstallable()
}

func (m *menu) Render() webapp.UI {
	linkClass := "link heading fit unselectable"

	isFocus := func(path string) string {
		if webapp.Window().URL().Path == path {
			return "focus"
		}
		return ""
	}

	return webui.Scroll().
		Class("menu").
		Class(m.Iclass).
		HeaderHeight(headerHeight).
		Header(
			webui.Stack().
				Class("fill").
				Middle().
				Content(
					webapp.Header().Body(
						webapp.A().
							Class("heading").
							Class("app-title").
							Href("/").
							Text("Bhojpur Web"),
					),
				),
		).
		Content(
			webapp.Nav().Body(
				webapp.Div().Class("separator"),

				webui.Link().
					Class(linkClass).
					Icon(homeSVG).
					Label("Home").
					Href("/").
					Class(isFocus("/")),

				webapp.Div().Class("separator"),

				webui.Link().
					Class(linkClass).
					Icon(gridSVG).
					Label("Components").
					Href("/components").
					Class(isFocus("/components")),

				webapp.Div().Class("separator"),

				webui.Link().
					Class(linkClass).
					Icon(swapSVG).
					Label("Migrate Application").
					Href("/migrate").
					Class(isFocus("/migrate")),

				webapp.Div().Class("separator"),

				webui.Link().
					Class(linkClass).
					Icon(githubSVG).
					Label("GitHub").
					Href(githubURL),

				webapp.Div().Class("separator"),

				webapp.If(m.appInstallable,
					webui.Link().
						Class(linkClass).
						Icon(downloadSVG).
						Label("Install").
						OnClick(m.installApp),
				),
				webui.Link().
					Class(linkClass).
					Icon(userLockSVG).
					Label("Privacy Policy").
					Href("/privacy-policy").
					Class(isFocus("/privacy-policy")),

				webapp.Div().Class("separator"),
			),
		)
}

func (m *menu) installApp(ctx webapp.Context, e webapp.Event) {
	ctx.NewAction(installApp)
}
