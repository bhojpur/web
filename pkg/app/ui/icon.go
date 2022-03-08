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
	"fmt"
	"strings"

	"github.com/bhojpur/web/pkg/app"
)

// IIcon is the interface that describes an icon.
type IIcon interface {
	app.UI

	// Sets the ID.
	ID(v string) IIcon

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v string) IIcon

	// Sets the style. Multiple styles can be defined by successive calls.
	Style(k, v string) IIcon

	// Sets the icon horizontal and vertical size in px.
	Size(px int) IIcon

	// Sets the SVG code or the source location.
	Src(v string) IIcon
}

// Icon creates an icon.
func Icon() IIcon {
	return &icon{
		Isize: DefaultIconSize,
	}
}

type icon struct {
	app.Compo

	Iid     string
	Iclass  string
	Istyles []style
	Isize   int
	Isrc    string
}

func (i *icon) ID(v string) IIcon {
	i.Iid = v
	return i
}

func (i *icon) Class(v string) IIcon {
	i.Iclass = app.AppendClass(i.Iclass, v)
	return i
}

func (i *icon) Style(k, v string) IIcon {
	if v == "" {
		return i
	}
	i.Istyles = append(i.Istyles, style{
		key:   k,
		value: v,
	})
	return i
}

func (i *icon) Size(px int) IIcon {
	i.Isize = px
	return i
}

func (i *icon) Src(v string) IIcon {
	i.Isrc = v
	return i
}

func (i *icon) Render() app.UI {
	var content app.UI
	if isSVG(i.Isrc) {
		content = app.Raw(fmt.Sprintf(i.Isrc, i.Isize, i.Isize))
	} else {
		content = app.Img().
			Style("max-width", "100%").
			Style("max-height", "100%").
			Src(i.Isrc)
	}

	icon := app.Div().
		DataSet("bhojpur", "Icon").
		ID(i.Iid).
		Class(i.Iclass).
		Style("width", pxToString(i.Isize)).
		Style("height", pxToString(i.Isize)).
		Style("max-width", pxToString(i.Isize)).
		Style("max-height", pxToString(i.Isize)).
		Style("min-width", pxToString(i.Isize)).
		Style("min-height", pxToString(i.Isize)).
		Body(content)
	for _, s := range i.Istyles {
		icon.Style(s.key, s.value)
	}
	return icon
}

func isSVG(v string) bool {
	return strings.HasPrefix(strings.TrimSpace(v), "<svg")
}
