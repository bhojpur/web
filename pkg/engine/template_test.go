package engine

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
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	assetfs "github.com/bhojpur/web/pkg/synthesis"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/test"
)

var header = `{{define "header"}}
<h1>Hello, Bhojpur!</h1>
{{end}}`

var index = `<!DOCTYPE html>
<html>
  <head>
    <title>Bhojpur Web - Welcome Template</title>
  </head>
  <body>
{{template "block"}}
{{template "header"}}
{{template "blocks/block.tpl"}}
  </body>
</html>
`

var block = `{{define "block"}}
<h1>Hello, blocks!</h1>
{{end}}`

func TestTemplate(t *testing.T) {
	wkdir, err := os.Getwd()
	assert.Nil(t, err)
	dir := filepath.Join(wkdir, "_beeTmp", "TestTemplate")
	files := []string{
		"header.tpl",
		"index.tpl",
		"blocks/block.tpl",
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		t.Fatal(err)
	}
	for k, name := range files {
		dirErr := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0777)
		assert.Nil(t, dirErr)
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			if k == 0 {
				f.WriteString(header)
			} else if k == 1 {
				f.WriteString(index)
			} else if k == 2 {
				f.WriteString(block)
			}

			f.Close()
		}
	}
	if err := AddViewPath(dir); err != nil {
		t.Fatal(err)
	}
	beeTemplates := webViewPathTemplates[dir]
	if len(beeTemplates) != 3 {
		t.Fatalf("should be 3 but got %v", len(beeTemplates))
	}
	if err := beeTemplates["index.tpl"].ExecuteTemplate(os.Stdout, "index.tpl", nil); err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		os.RemoveAll(filepath.Join(dir, name))
	}
	os.RemoveAll(dir)
}

var menu = `<div class="menu">
<ul>
<li>menu1</li>
<li>menu2</li>
<li>menu3</li>
</ul>
</div>
`
var user = `<!DOCTYPE html>
<html>
  <head>
    <title>Bhojpur Web - Welcome Template</title>
  </head>
  <body>
{{template "../public/menu.tpl"}}
  </body>
</html>
`

func TestRelativeTemplate(t *testing.T) {
	wkdir, err := os.Getwd()
	assert.Nil(t, err)
	dir := filepath.Join(wkdir, "_beeTmp")

	//Just add dir to known viewPaths
	if err := AddViewPath(dir); err != nil {
		t.Fatal(err)
	}

	files := []string{
		"easyui/public/menu.tpl",
		"easyui/rbac/user.tpl",
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		t.Fatal(err)
	}
	for k, name := range files {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0777)
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			if k == 0 {
				f.WriteString(menu)
			} else if k == 1 {
				f.WriteString(user)
			}
			f.Close()
		}
	}
	if err := BuildTemplate(dir, files[1]); err != nil {
		t.Fatal(err)
	}
	beeTemplates := beeViewPathTemplates[dir]
	if err := beeTemplates["easyui/rbac/user.tpl"].ExecuteTemplate(os.Stdout, "easyui/rbac/user.tpl", nil); err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		os.RemoveAll(filepath.Join(dir, name))
	}
	os.RemoveAll(dir)
}

var add = `{{ template "layout_blog.tpl" . }}
{{ define "css" }}
        <link rel="stylesheet" href="/static/css/current.css">
{{ end}}


{{ define "content" }}
        <h2>{{ .Title }}</h2>
        <p> This is SomeVar: {{ .SomeVar }}</p>
{{ end }}

{{ define "js" }}
    <script src="/static/js/current.js"></script>
{{ end}}`

var layoutBlog = `<!DOCTYPE html>
<html>
<head>
    <title>Pramila Kumari</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <link rel="stylesheet" href="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
    <link rel="stylesheet" href="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap-theme.min.css">
     {{ block "css" . }}{{ end }}
</head>
<body>

    <div class="container">
        {{ block "content" . }}{{ end }}
    </div>
    <script type="text/javascript" src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js"></script>
     {{ block "js" . }}{{ end }}
</body>
</html>`

var output = `<!DOCTYPE html>
<html>
<head>
    <title>Pramila Kumari</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <link rel="stylesheet" href="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
    <link rel="stylesheet" href="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap-theme.min.css">
     
        <link rel="stylesheet" href="/static/css/current.css">

</head>
<body>

    <div class="container">
        
        <h2>Hello</h2>
        <p> This is SomeVar: val</p>

    </div>
    <script type="text/javascript" src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="http://netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js"></script>
     
    <script src="/static/js/current.js"></script>

</body>
</html>





`

func TestTemplateLayout(t *testing.T) {
	wkdir, err := os.Getwd()
	assert.Nil(t, err)

	dir := filepath.Join(wkdir, "_beeTmp", "TestTemplateLayout")
	files := []string{
		"add.tpl",
		"layout_blog.tpl",
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		t.Fatal(err)
	}

	for k, name := range files {
		dirErr := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0777)
		assert.Nil(t, dirErr)
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			if k == 0 {
				_, writeErr := f.WriteString(add)
				assert.Nil(t, writeErr)
			} else if k == 1 {
				_, writeErr := f.WriteString(layoutBlog)
				assert.Nil(t, writeErr)
			}
			clErr := f.Close()
			assert.Nil(t, clErr)
		}
	}
	if err := AddViewPath(dir); err != nil {
		t.Fatal(err)
	}
	beeTemplates := beeViewPathTemplates[dir]
	if len(beeTemplates) != 2 {
		t.Fatalf("should be 2 but got %v", len(beeTemplates))
	}
	out := bytes.NewBufferString("")

	if err := beeTemplates["add.tpl"].ExecuteTemplate(out, "add.tpl", map[string]string{"Title": "Hello", "SomeVar": "val"}); err != nil {
		t.Fatal(err)
	}
	if out.String() != output {
		t.Log(out.String())
		t.Fatal("Compare failed")
	}
	for _, name := range files {
		os.RemoveAll(filepath.Join(dir, name))
	}
	os.RemoveAll(dir)
}

type TestingFileSystem struct {
	assetfs *assetfs.AssetFS
}

func (d TestingFileSystem) Open(name string) (http.File, error) {
	return d.assetfs.Open(name)
}

var outputBhojpur = `<!DOCTYPE html>
<html>
  <head>
    <title>Bhojpur Web - Welcome Template</title>
  </head>
  <body>

	
<h1>Hello, blocks!</h1>

	
<h1>Hello, Bhojpur!</h1>

	

	<h2>Hello</h2>
	<p> This is SomeVar: val</p>
  </body>
</html>
`

func TestFsSynthesis(t *testing.T) {
	SetTemplateFSFunc(func() http.FileSystem {
		return TestingFileSystem{&assetfs.AssetFS{Asset: test.Asset, AssetDir: test.AssetDir, AssetInfo: test.AssetInfo}}
	})
	dir := "views"
	if err := AddViewPath("views"); err != nil {
		t.Fatal(err)
	}
	beeTemplates := webViewPathTemplates[dir]
	if len(beeTemplates) != 3 {
		t.Fatalf("should be 3 but got %v", len(beeTemplates))
	}
	if err := beeTemplates["index.tpl"].ExecuteTemplate(os.Stdout, "index.tpl", map[string]string{"Title": "Hello", "SomeVar": "val"}); err != nil {
		t.Fatal(err)
	}
	out := bytes.NewBufferString("")
	if err := beeTemplates["index.tpl"].ExecuteTemplate(out, "index.tpl", map[string]string{"Title": "Hello", "SomeVar": "val"}); err != nil {
		t.Fatal(err)
	}

	if out.String() != outputBhojpur {
		t.Log(out.String())
		t.Fatal("Compare failed")
	}
}
