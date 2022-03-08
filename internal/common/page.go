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
	"fmt"

	webapp "github.com/bhojpur/web/pkg/app"
	webui "github.com/bhojpur/web/pkg/app/ui"
)

const (
	headerHeight  = 72
	adsenseClient = "ca-pub-1234"
	adsenseSlot   = "1234"
)

type page struct {
	webapp.Compo

	Iclass   string
	Iindex   []webapp.UI
	Iicon    string
	Ititle   string
	Icontent []webapp.UI

	updateAvailable bool
}

func newPage() *page {
	return &page{}
}

func (p *page) Index(v ...webapp.UI) *page {
	p.Iindex = webapp.FilterUIElems(v...)
	return p
}

func (p *page) Icon(v string) *page {
	p.Iicon = v
	return p
}

func (p *page) Title(v string) *page {
	p.Ititle = v
	return p
}

func (p *page) Content(v ...webapp.UI) *page {
	p.Icontent = webapp.FilterUIElems(v...)
	return p
}

func (p *page) OnNav(ctx webapp.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
	ctx.Defer(scrollTo)
}

func (p *page) OnAppUpdate(ctx webapp.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *page) Render() webapp.UI {
	return webui.Shell().
		Class("fill").
		Class("background").
		HamburgerMenu(
			newMenu().
				Class("fill").
				Class("menu-hamburger-background"),
		).
		Menu(
			newMenu().Class("fill"),
		).
		Index(
			webapp.If(len(p.Iindex) != 0,
				webui.Scroll().
					Class("fill").
					HeaderHeight(headerHeight).
					Content(
						webapp.Nav().
							Class("index").
							Body(
								webapp.Div().Class("separator"),
								webapp.Header().
									Class("h2").
									Text("Index"),
								webapp.Div().Class("separator"),
								webapp.Range(p.Iindex).Slice(func(i int) webapp.UI {
									return p.Iindex[i]
								}),
								newIndexLink().Title("Report an Issue"),
								webapp.Div().Class("separator"),
							),
					),
			),
		).
		Content(
			webui.Scroll().
				Class("fill").
				Header(
					webapp.Nav().
						Class("fill").
						Body(
							webui.Stack().
								Class("fill").
								Right().
								Middle().
								Content(
									webapp.If(p.updateAvailable,
										webapp.Div().
											Class("link-update").
											Body(
												webui.Link().
													Class("link").
													Class("heading").
													Class("fit").
													Class("unselectable").
													Icon(downloadSVG).
													Label("Update").
													OnClick(p.updateApp),
											),
									),
								),
						),
				).
				HeaderHeight(headerHeight).
				Content(
					webapp.Main().Body(
						webapp.Article().Body(
							webapp.Header().
								ID("page-top").
								Class("page-title").
								Class("center").
								Body(
									webui.Stack().
										Center().
										Middle().
										Content(
											webui.Icon().
												Class("icon-left").
												Class("unselectable").
												Size(90).
												Src(p.Iicon),
											webapp.H1().Text(p.Ititle),
										),
								),
							webapp.Div().Class("separator"),
							webapp.Range(p.Icontent).Slice(func(i int) webapp.UI {
								return p.Icontent[i]
							}),

							webapp.Div().Class("separator"),
							webapp.Aside().Body(
								webapp.Header().
									ID("repport-an-issue").
									Class("h2").
									Text("Report an issue"),
								webapp.P().Body(
									webapp.Text("Found something incorrect, a typo or have suggestions to improve this page? "),
									webapp.A().
										Href(fmt.Sprintf(
											"%s/issues/new?title=Documentation issue in %s page",
											githubURL,
											p.Ititle,
										)).
										Text("🚀 Submit a GitHub issue!"),
								),
							),
							webapp.Div().Class("separator"),
						),
					),
				),
		).
		Ads(
			webui.Flyer().
				Class("fill").
				HeaderHeight(headerHeight).
				Banner(
					webapp.Aside().
						Class("fill").
						Body(
							webui.AdsenseDisplay().
								Class("fill").
								Class("no-scroll").
								Client(adsenseClient).
								Slot(adsenseSlot),
						),
				).
				PremiumHeight(200).
				Premium(
					newGithubSponsor().Class("fill"),
				),
		)
}

func (p *page) updateApp(ctx webapp.Context, e webapp.Event) {
	ctx.NewAction(updateApp)
}

func scrollTo(ctx webapp.Context) {
	id := ctx.Page().URL().Fragment
	if id == "" {
		id = "page-top"
	}
	ctx.ScrollTo(id)
}

func fragmentFocus(fragment string) string {
	if fragment == webapp.Window().URL().Fragment {
		return "focus"
	}
	return ""
}
