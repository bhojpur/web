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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateStaticWebsite(t *testing.T) {
	testSkipWasm(t)

	dir := "static-test"
	defer os.RemoveAll(dir)

	err := GenerateStaticWebsite(dir,
		&Handler{
			Name:      "Static Bhojpur",
			Title:     "Static test",
			Resources: GitHubPages("bhojpur"),
		},
		"/hello",
		"world",
		"/nested/foo",
	)
	require.NoError(t, err)

	files := []string{
		filepath.Join(dir),
		filepath.Join(dir, "web"),
		filepath.Join(dir, "index.html"),
		filepath.Join(dir, "wasm_exec.js"),
		filepath.Join(dir, "app.js"),
		filepath.Join(dir, "app-worker.js"),
		filepath.Join(dir, "manifest.webmanifest"),
		filepath.Join(dir, "app.css"),
		filepath.Join(dir, "hello.html"),
		filepath.Join(dir, "world.html"),
		filepath.Join(dir, "nested", "foo.html"),
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			_, err := os.Stat(f)
			require.NoError(t, err)
		})
	}
}
