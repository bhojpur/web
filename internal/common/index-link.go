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
	"strings"

	webapp "github.com/bhojpur/web/pkg/app"
)

type indexLink struct {
	webapp.Compo

	Iclass string
	Ititle string
	Ihref  string
}

func newIndexLink() *indexLink {
	return &indexLink{}
}

func (l *indexLink) Class(v string) *indexLink {
	l.Iclass = webapp.AppendClass(l.Iclass, v)
	return l
}

func (l *indexLink) Title(v string) *indexLink {
	l.Ititle = v
	return l
}

func (l *indexLink) Href(v string) *indexLink {
	l.Ihref = v
	return l
}

func (l *indexLink) OnNav(ctx webapp.Context) {}

func (l *indexLink) Render() webapp.UI {
	fragment := titleToFragment(l.Ititle)

	href := l.Ihref
	if href == "" {
		href = "#" + fragment
	}

	return webapp.A().
		Class("index-link").
		Class(l.Iclass).
		Class(fragmentFocus(fragment)).
		Href(href).
		Text(l.Ititle).
		Title(l.Ititle)
}

func titleToFragment(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	v = strings.ReplaceAll(v, " ", "-")
	v = strings.ReplaceAll(v, ".", "-")
	v = strings.ReplaceAll(v, "?", "")
	return v
}
