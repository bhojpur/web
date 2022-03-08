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
	"path/filepath"

	webapp "github.com/bhojpur/web/pkg/app"
	webui "github.com/bhojpur/web/pkg/app/ui"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type markdownDoc struct {
	webapp.Compo

	Iid    string
	Iclass string
	Imd    string
}

func newMarkdownDoc() *markdownDoc {
	return &markdownDoc{}
}

func (d *markdownDoc) ID(v string) *markdownDoc {
	d.Iid = v
	return d
}

func (d *markdownDoc) Class(v string) *markdownDoc {
	d.Iclass = webapp.AppendClass(d.Iclass, v)
	return d
}

func (d *markdownDoc) MD(v string) *markdownDoc {
	d.Imd = fmt.Sprintf(`<div class="markdown">%s</div>`, parseMarkdown([]byte(v)))
	return d
}

func (d *markdownDoc) OnMount(ctx webapp.Context) {
	ctx.Defer(d.highlightCode)
}

func (d *markdownDoc) OnUpdate(ctx webapp.Context) {
	ctx.Defer(d.highlightCode)
}

func (d *markdownDoc) Render() webapp.UI {
	return webapp.Div().
		ID(d.Iid).
		Class(d.Iclass).
		Body(
			webapp.Raw(d.Imd),
		)
}

func (d *markdownDoc) highlightCode(ctx webapp.Context) {
	webapp.Window().Get("Prism").Call("highlightAll")
}

func parseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	return markdown.ToHTML(md, parser, nil)
}

type remoteMarkdownDoc struct {
	webapp.Compo

	Iid    string
	Iclass string
	Isrc   string

	md markdownContent
}

func newRemoteMarkdownDoc() *remoteMarkdownDoc {
	return &remoteMarkdownDoc{}
}

func (d *remoteMarkdownDoc) ID(v string) *remoteMarkdownDoc {
	d.Iid = v
	return d
}

func (d *remoteMarkdownDoc) Class(v string) *remoteMarkdownDoc {
	d.Iclass = webapp.AppendClass(d.Iclass, v)
	return d
}

func (d *remoteMarkdownDoc) Src(v string) *remoteMarkdownDoc {
	d.Isrc = v
	return d
}

func (d *remoteMarkdownDoc) OnMount(ctx webapp.Context) {
	d.load(ctx)
}

func (d *remoteMarkdownDoc) OnUpdate(ctx webapp.Context) {
	d.load(ctx)
}

func (d *remoteMarkdownDoc) load(ctx webapp.Context) {
	src := d.Isrc
	ctx.ObserveState(markdownState(src)).
		While(func() bool {
			return src == d.Isrc
		}).
		OnChange(func() {
			ctx.Defer(scrollTo)
		}).
		Value(&d.md)

	ctx.NewAction(getMarkdown, webapp.T("path", d.Isrc))
}

func (d *remoteMarkdownDoc) Render() webapp.UI {
	return webapp.Div().
		ID(d.Iid).
		Class(d.Iclass).
		Body(
			webui.Loader().
				Class("heading").
				Class("fill").
				Loading(d.md.Status == loading).
				Err(d.md.Err).
				Label(fmt.Sprintf("Loading %s...", filepath.Base(d.Isrc))),
			webapp.If(d.md.Status == loaded,
				newMarkdownDoc().
					Class("fill").
					MD(d.md.Data),
			).Else(),
		)
}
