//go:build !wasm

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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/web/pkg/app/errors"
)

// GenerateStaticWebsite generates the files to run a PWA built with Bhojpur Web application
// as a static website in the specified directory. Static websites can be used with hosts
// such as Github Pages.
//
// Note that app.wasm must still be built separately and put into the web directory.
func GenerateStaticWebsite(dir string, h *Handler, pages ...string) error {
	if dir == "" {
		dir = "."
	}

	resources := map[string]struct{}{
		"/":                     {},
		"/wasm_exec.js":         {},
		"/app.js":               {},
		"/app-worker.js":        {},
		"/manifest.webmanifest": {},
		"/app.css":              {},
		"/web":                  {},
	}

	for path := range routes.routes {
		resources[path] = struct{}{}
	}

	for _, p := range pages {
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		resources[p] = struct{}{}
	}

	server := httptest.NewServer(h)
	defer server.Close()

	for path := range resources {
		switch path {
		case "/web":
			if err := createStaticDir(filepath.Join(dir, path), ""); err != nil {
				return errors.New("creating web directory failed").Wrap(err)
			}

		default:
			filename := path
			if filename == "/" {
				filename = "/index.html"
			}

			f, err := createStaticFile(dir, filename)
			if err != nil {
				return errors.New("creating file failed").
					Tag("path", path).
					Tag("filename", filename).
					Wrap(err)
			}
			defer f.Close()

			page, err := createStaticPage(server.URL + path)
			if err != nil {
				return errors.New("creating page failed").
					Tag("path", path).
					Tag("filename", filename).
					Wrap(err)
			}

			if n, err := f.Write(page); err != nil {
				return errors.New("writing page failed").
					Tag("path", path).
					Tag("filename", filename).
					Tag("bytes-written", n).
					Wrap(err)
			}
		}
	}

	return nil
}

func createStaticDir(dir, path string) error {
	dir = filepath.Join(dir, filepath.Dir(path))
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}
	return os.MkdirAll(filepath.Join(dir), 0755)
}

func createStaticFile(dir, path string) (*os.File, error) {
	if err := createStaticDir(dir, path); err != nil {
		return nil, errors.New("creating file directory failed").Wrap(err)
	}

	filename := filepath.Join(dir, path)
	if filepath.Ext(filename) == "" {
		filename += ".html"
	}

	return os.Create(filename)
}

func createStaticPage(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("creating http request failed").
			Tag("path", path).
			Wrap(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("http request failed").
			Tag("path", path).
			Wrap(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("reading request body failed").
			Tag("path", path).
			Wrap(err)
	}
	return body, nil
}
