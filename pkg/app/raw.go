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

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/bhojpur/web/pkg/app/errors"
)

// Raw returns a ui element from the given raw value. HTML raw value must have a
// single root.
//
// It is not recommended to use this kind of node since there is no check on the
// raw string content.
func Raw(v string) UI {
	v = strings.TrimSpace(v)

	tag := rawRootTagName(v)
	if tag == "" {
		v = "<div></div>"
	}

	return &raw{
		value: v,
		tag:   tag,
	}
}

type raw struct {
	disp       Dispatcher
	jsvalue    Value
	parentElem UI
	tag        string
	value      string
}

func (r *raw) Kind() Kind {
	return RawHTML
}

func (r *raw) JSValue() Value {
	return r.jsvalue
}

func (r *raw) Mounted() bool {
	return r.jsvalue != nil && r.dispatcher() != nil
}

func (r *raw) name() string {
	return "raw." + r.tag
}

func (r *raw) self() UI {
	return r
}

func (r *raw) setSelf(UI) {
}

func (r *raw) context() context.Context {
	return nil
}

func (r *raw) dispatcher() Dispatcher {
	return r.disp
}

func (r *raw) attributes() map[string]string {
	return nil
}

func (r *raw) eventHandlers() map[string]eventHandler {
	return nil
}

func (r *raw) parent() UI {
	return r.parentElem
}

func (r *raw) setParent(p UI) {
	r.parentElem = p
}

func (r *raw) children() []UI {
	return nil
}

func (r *raw) mount(d Dispatcher) error {
	if r.Mounted() {
		return errors.New("mounting raw html element failed").
			Tag("reason", "already mounted").
			Tag("name", r.name()).
			Tag("kind", r.Kind())
	}

	r.disp = d

	wrapper, err := Window().createElement("div")
	if err != nil {
		return errors.New("creating raw node wrapper failed").Wrap(err)
	}

	if IsServer {
		r.jsvalue = wrapper
		return nil
	}

	wrapper.setInnerHTML(r.value)
	value := wrapper.firstChild()
	if !value.Truthy() {
		return errors.New("mounting raw html element failed").
			Tag("reason", "converting raw html to html elements returned nil").
			Tag("name", r.name()).
			Tag("kind", r.Kind()).
			Tag("raw-html", r.value)
	}
	wrapper.removeChild(value)
	r.jsvalue = value
	return nil
}

func (r *raw) dismount() {
	r.jsvalue = nil
}

func (r *raw) update(n UI) error {
	if !r.Mounted() {
		return nil
	}

	if n.Kind() != r.Kind() || n.name() != r.name() {
		return errors.New("updating raw html element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", r.Kind()).
			Tag("current-name", r.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	if v := n.(*raw).value; r.value != v {
		return errors.New("updating raw html element failed").
			Tag("replace", true).
			Tag("reason", "different raw values").
			Tag("current-value", r.value).
			Tag("new-value", v)
	}

	return nil
}

func (r *raw) onNav(*url.URL) {
}

func (r *raw) onAppUpdate() {
}

func (r *raw) onAppInstallChange() {
}

func (r *raw) onResize() {
}

func (r *raw) preRender(Page) {
}

func (r *raw) html(w io.Writer) {
	w.Write([]byte(r.value))
}

func (r *raw) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte(r.value))
}

func rawRootTagName(raw string) string {
	raw = strings.TrimSpace(raw)

	if strings.HasPrefix(raw, "</") || !strings.HasPrefix(raw, "<") {
		return ""
	}

	end := -1
	for i := 1; i < len(raw); i++ {
		if raw[i] == ' ' ||
			raw[i] == '\t' ||
			raw[i] == '\n' ||
			raw[i] == '>' {
			end = i
			break
		}
	}

	if end <= 0 {
		return ""
	}

	return raw[1:end]
}
