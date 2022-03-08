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

var (
	// NotFound is the ui element that is displayed when a request is not
	// routed.
	NotFound UI = &notFound{}
)

type notFound struct {
	Compo
	Icon string
}

func (n *notFound) OnMount(Context) {
	links := Window().Get("document").Call("getElementsByTagName", "link")

	for i := 0; i < links.Length(); i++ {
		link := links.Index(i)
		rel := link.Call("getAttribute", "rel")

		if rel.String() == "icon" {
			favicon := link.Call("getAttribute", "href")
			n.Icon = favicon.String()
			return
		}
	}
}

func (n *notFound) Render() UI {
	return Div().
		Class("bhojpur-app-info").
		Body(
			Div().
				Class("bhojpur-notfound-title").
				Body(
					Text("4"),
					Img().
						Class("bhojpur-logo").
						Alt("0").
						Src(n.Icon),
					Text("4"),
				),
			P().
				Class("bhojpur-label").
				Text("Not Found"),
		)
}
