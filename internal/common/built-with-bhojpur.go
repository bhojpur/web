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

type builtWithBhojpur struct {
	webapp.Compo

	Iid    string
	Iclass string
}

func newBuiltWithBhojpur() *builtWithBhojpur {
	return &builtWithBhojpur{}
}

func (b *builtWithBhojpur) ID(v string) *builtWithBhojpur {
	b.Iid = v
	return b
}

func (b *builtWithBhojpur) Class(v string) *builtWithBhojpur {
	b.Iclass = webapp.AppendClass(b.Iclass, v)
	return b
}

func (b *builtWithBhojpur) Render() webapp.UI {
	return webapp.Div().
		Class(b.Iclass).
		Body(
			webapp.H2().
				ID(b.Iid).
				Text("Built with Bhojpur Web"),
			webui.Flow().
				Class("p").
				StretchItems().
				Spacing(18).
				ItemWidth(360).
				Content(
					newBuiltWithBhojpurItem().
						Class("fill").
						Image("https://static.bhojpur.net/image/logo.png").
						Name("iam.bhojpur.net").
						Description("An identity & access control solution").
						Href("https://iam.bhojpur.net"),
					newBuiltWithBhojpurItem().
						Class("fill").
						Image("https://static.bhojpur.net/image/logo.png").
						Name("congress.bhojpur.net").
						Description("A web conference management system").
						Href("https://congress.bhojpur.net"),
					newBuiltWithBhojpurItem().
						Class("fill").
						Image("https://static.bhojpur.net/image/logo.png").
						Name("ode.bhojpur.net").
						Description("An optical data engine for web-based microscopy").
						Href("https://ode.bhojpur.net"),
					newBuiltWithBhojpurItem().
						Class("fill").
						Image("https://static.bhojpur.net/image/logo.png").
						Name("nsm.bhojpur.net").
						Description("A network service mesh").
						Href("https://nsm.bhojpur.net"),
				),
		)
}

type builtWithBhojpurItem struct {
	webapp.Compo

	Iclass       string
	Iimage       string
	Iname        string
	Idescription string
	Ihref        string
}

func newBuiltWithBhojpurItem() *builtWithBhojpurItem {
	return &builtWithBhojpurItem{}
}

func (i *builtWithBhojpurItem) Class(v string) *builtWithBhojpurItem {
	i.Iclass = webapp.AppendClass(i.Iclass, v)
	return i
}

func (i *builtWithBhojpurItem) Image(v string) *builtWithBhojpurItem {
	i.Iimage = v
	return i
}

func (i *builtWithBhojpurItem) Name(v string) *builtWithBhojpurItem {
	i.Iname = v
	return i
}

func (i *builtWithBhojpurItem) Description(v string) *builtWithBhojpurItem {
	i.Idescription = v
	return i
}

func (i *builtWithBhojpurItem) Href(v string) *builtWithBhojpurItem {
	i.Ihref = v
	return i
}

func (i *builtWithBhojpurItem) Render() webapp.UI {
	return webapp.A().
		Class(i.Iclass).
		Class("block").
		Class("rounded").
		Class("text-center").
		Class("magnify").
		Class("default").
		Href(i.Ihref).
		Body(
			webui.Block().
				Class("fill").
				Middle().
				Content(
					webapp.Img().
						Class("hstretch").
						Alt(i.Iname+" tumbnail.").
						Src(i.Iimage),
					webapp.H3().Text(i.Iname),
					webapp.Div().
						Class("text-tiny-top").
						Text(i.Idescription),
				),
		)
}
