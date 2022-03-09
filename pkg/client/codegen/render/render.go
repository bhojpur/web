package render

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
	"net/http"
	"path"

	tmplsvr "github.com/bhojpur/web/pkg/template"
)

//	var render = templaterender.NewRender("./templates")
//
//	http.HandleFunc("/m", func(w http.ResponseWriter, req *http.Request) {
//		render.HTML(w, 200, "index.html", template.Context{"aa": "eeeeeee"})
//	})
//	http.ListenAndServe(":9005", nil)

// --------------------------------------------------------------------------------
var htmlContentType = []string{"text/html; charset=utf-8"}

type Render struct {
	TemplateDir string
	Cache       bool
}

func NewRender(templateDir string) *Render {
	var r = &Render{}
	r.TemplateDir = templateDir
	return r
}

func (this *Render) Template(name string) *Template {
	var template *tmplsvr.Template
	var filename string
	if len(this.TemplateDir) > 0 {
		filename = path.Join(this.TemplateDir, name)
	} else {
		filename = name
	}

	if this.Cache {
		template = tmplsvr.Must(tmplsvr.DefaultSet.FromCache(filename))
	} else {
		template = tmplsvr.Must(tmplsvr.DefaultSet.FromFile(filename))
	}

	if template == nil {
		panic("template " + name + " not exists")
		return nil
	}

	var r = &Template{}
	r.template = template
	return r
}

func (this *Render) TemplateFromString(tpl string) *Template {
	var template = tmplsvr.Must(tmplsvr.DefaultSet.FromString(tpl))
	var r = &Template{}
	r.template = template
	return r
}

func (this *Render) HTML(w http.ResponseWriter, status int, name string, data interface{}) {
	w.WriteHeader(status)
	this.Template(name).ExecuteWriter(w, data)
}

// --------------------------------------------------------------------------------
type Template struct {
	template *tmplsvr.Template
	context  tmplsvr.Context
}

func (this *Template) ExecuteWriter(w http.ResponseWriter, data interface{}) (err error) {
	WriteContentType(w, htmlContentType)
	this.context = DataToContext(data)
	err = this.template.ExecuteWriter(this.context, w)
	return err
}

func (this *Template) Execute(data interface{}) (string, error) {
	this.context = DataToContext(data)
	return this.template.Execute(this.context)
}

// --------------------------------------------------------------------------------
func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func DataToContext(data interface{}) tmplsvr.Context {
	var ctx tmplsvr.Context
	if data != nil {
		switch data.(type) {
		case tmplsvr.Context:
			ctx = data.(tmplsvr.Context)
		case map[string]interface{}:
			ctx = tmplsvr.Context(data.(map[string]interface{}))
		}
	}
	return ctx
}
