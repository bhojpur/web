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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalDir(t *testing.T) {
	testSkipWasm(t)

	h, _ := LocalDir("test").(localDir)
	require.Equal(t, "test", h.Static())
	require.Equal(t, "test/web/app.wasm", h.AppWASM())

	close := testCreateDir(t, "test/web")
	defer close()

	resources := []string{
		"/web/test",
		"/web/app.wasm",
	}

	for _, r := range resources {
		t.Run(r, func(t *testing.T) {
			path := strings.Replace(r, "/web", "test/web", 1)
			err := ioutil.WriteFile(path, []byte("hello"), 0666)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, r, nil)
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			require.Equal(t, "hello", res.Body.String())
		})
	}
}

func TestRemoteBucket(t *testing.T) {
	utests := []struct {
		scenario string
		provider ResourceProvider
	}{
		{
			scenario: "remote bucket",
			provider: RemoteBucket("https://storage.googleapis.com/test"),
		},
		{
			scenario: "remote bucket with web suffix",
			provider: RemoteBucket("https://storage.googleapis.com/test/web/"),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, "https://storage.googleapis.com/test", u.provider.Static())
			require.Equal(t, "https://storage.googleapis.com/test/web/app.wasm", u.provider.AppWASM())
		})
	}
}
