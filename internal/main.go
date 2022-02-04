package main

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
	"net/http"

	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/bhojpur/web/pkg/synthesis"
	test "github.com/bhojpur/web/test"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "भोजपुर जिला घर बा, तब कौना बात के डर बा !!")
}

func namasteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "नमस्ते, %s!", r.FormValue("नाम"))
}

func main() {
	websvr.Run("127.0.0.1:3000")
	websvr.AddFuncMap("/", http.HandlerFunc(indexHandler))
	websvr.AddFuncMap("/अभिवादन/:नाम", http.HandlerFunc(namasteHandler))
	websvr.AddFuncMap("/",
		http.FileServer(
			&synthesis.AssetFS{
				Asset:     test.Asset,
				AssetDir:  test.AssetDir,
				AssetInfo: test.AssetInfo,
				Prefix:    "data",
				Fallback:  "index.html",
			}))
}
