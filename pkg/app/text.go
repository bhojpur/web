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
	"html"
	"io"
	"net/url"

	"github.com/bhojpur/web/pkg/app/errors"
)

// Text creates a simple text element.
func Text(v interface{}) UI {
	return &text{value: toString(v)}
}

type text struct {
	disp       Dispatcher
	jsvalue    Value
	parentElem UI
	value      string
}

func (t *text) Kind() Kind {
	return SimpleText
}

func (t *text) JSValue() Value {
	return t.jsvalue
}

func (t *text) Mounted() bool {
	return t.jsvalue != nil && t.dispatcher() != nil
}

func (t *text) name() string {
	return "text"
}

func (t *text) self() UI {
	return t
}

func (t *text) setSelf(n UI) {
}

func (t *text) context() context.Context {
	return context.TODO()
}

func (t *text) dispatcher() Dispatcher {
	return t.disp
}

func (t *text) attributes() map[string]string {
	return nil
}

func (t *text) eventHandlers() map[string]eventHandler {
	return nil
}

func (t *text) parent() UI {
	return t.parentElem
}

func (t *text) setParent(p UI) {
	t.parentElem = p
}

func (t *text) children() []UI {
	return nil
}

func (t *text) mount(d Dispatcher) error {
	if t.Mounted() {
		return errors.New("mounting ui element failed").
			Tag("reason", "already mounted").
			Tag("kind", t.Kind()).
			Tag("name", t.name()).
			Tag("value", t.value)
	}

	t.disp = d
	t.jsvalue = Window().createTextNode(t.value)
	return nil
}

func (t *text) dismount() {
	t.jsvalue = nil
}

func (t *text) update(n UI) error {
	if !t.Mounted() {
		return nil
	}

	o, isText := n.(*text)
	if !isText {
		return errors.New("updating ui element failed").
			Tag("replace", true).
			Tag("reason", "different element types").
			Tag("current-kind", t.Kind()).
			Tag("current-name", t.name()).
			Tag("updated-kind", n.Kind()).
			Tag("updated-name", n.name())
	}

	if t.value != o.value {
		t.value = o.value
		t.JSValue().setNodeValue(o.value)
	}

	return nil
}

func (t *text) onNav(*url.URL) {
}

func (t *text) onAppUpdate() {
}

func (t *text) onAppInstallChange() {
}

func (t *text) onResize() {
}

func (t *text) preRender(Page) {
}

func (t *text) html(w io.Writer) {
	w.Write([]byte(html.EscapeString(t.value)))
}

func (t *text) htmlWithIndent(w io.Writer, indent int) {
	writeIndent(w, indent)
	w.Write([]byte(html.EscapeString(t.value)))
}
