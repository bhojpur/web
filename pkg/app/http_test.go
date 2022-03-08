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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	Route("/", &preRenderTestCompo{})
}

type preRenderTestCompo struct {
	Compo
}

func (c *preRenderTestCompo) Render() UI {
	return Div().
		ID("pre-render-ok").
		Body(
			Img().Src("/web/resolve-static-resource-test.jpg"),
		)
}

func TestHandlerServePageWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: LocalDir(""),
		Title:     "Handler testing",
		Scripts: []string{
			"web/hello.js",
			"http://yunica.in/bar.js",
		},
		Styles: []string{
			"web/foo.css",
			"/web/bar.css",
			"http://yunica.in/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
		Image: "/web/test.png",
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `<html lang="en">`)
	require.Contains(t, body, `href="/web/foo.css"`)
	require.Contains(t, body, `href="/web/bar.css"`)
	require.Contains(t, body, `href="http://yunica.in/bar.css"`)
	require.Contains(t, body, `src="/web/hello.js"`)
	require.Contains(t, body, `src="http://yunica.in/bar.js"`)
	require.Contains(t, body, `href="/manifest.webmanifest"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `content="/web/test.png"`)
	require.Contains(t, body, `<img src="/web/resolve-static-resource-test.jpg">`)

	t.Log(body)
}

func TestHandlerServePageWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title:     "Handler testing",
		Resources: RemoteBucket("https://storage.googleapis.com/bhojpur/"),
		Scripts: []string{
			"/web/hello.js",
			"http://yunica.in/bar.js",
		},
		Styles: []string{
			"web/foo.css",
			"/web/bar.css",
			"http://yunica.in/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
		Image: "/web/test.png",
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="https://storage.googleapis.com/bhojpur/web/foo.css"`)
	require.Contains(t, body, `href="https://storage.googleapis.com/bhojpur/web/bar.css"`)
	require.Contains(t, body, `href="http://yunica.in/bar.css"`)
	require.Contains(t, body, `src="https://storage.googleapis.com/bhojpur/web/hello.js"`)
	require.Contains(t, body, `src="http://yunica.in/bar.js"`)
	require.Contains(t, body, `href="/manifest.webmanifest"`)
	require.Contains(t, body, `href="/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `content="https://storage.googleapis.com/bhojpur/web/test.png"`)
	require.Contains(t, body, `<img src="https://storage.googleapis.com/bhojpur/web/resolve-static-resource-test.jpg">`)

	t.Log(body)
}

func TestHandlerServePageWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Title:     "Handler testing",
		Resources: GitHubPages("bhojpur"),
		Scripts: []string{
			"/web/hello.js",
			"http://yunica.in/bar.js",
		},
		Styles: []string{
			"web/foo.css",
			"/web/bar.css",
			"http://yunica.in/bar.css",
		},
		RawHeaders: []string{
			`<meta http-equiv="refresh" content="30">`,
		},
	}
	h.Icon.AppleTouch = "ios.png"

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, body, `href="/bhojpur/web/foo.css"`)
	require.Contains(t, body, `href="/bhojpur/web/bar.css"`)
	require.Contains(t, body, `href="http://yunica.in/bar.css"`)
	require.Contains(t, body, `src="/bhojpur/web/hello.js"`)
	require.Contains(t, body, `src="http://yunica.in/bar.js"`)
	require.Contains(t, body, `href="/bhojpur/manifest.webmanifest"`)
	require.Contains(t, body, `href="/bhojpur/app.css"`)
	require.Contains(t, body, `<meta http-equiv="refresh" content="30">`)
	require.Contains(t, body, `<div id="pre-render-ok">`)
	require.Contains(t, body, `<img src="/bhojpur/web/resolve-static-resource-test.jpg">`)
	t.Log(body)
}

func TestHandlerServeWasmExecJS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/wasm_exec.js", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Equal(t, wasmExecJS, w.Body.String())
	t.Log(w.Body.String())
}

func TestHandlerServeAppJSWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		InternalURLs: []string{"https://redirect.me"},
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `register("/app-worker.js"`)
	require.Contains(t, body, `fetch("/web/app.wasm"`)
	require.Contains(t, body, "BHOJPUR_WEB_APP_VERSION")
	require.Contains(t, body, `"BHOJPUR_STATIC_RESOURCES_URL":""`)
	require.Contains(t, body, `"BHOJPUR_WEB_ROOT_PREFIX":""`)
	require.Contains(t, body, `"BHOJPUR_WEB_INTERNAL_URLS":"[\"https://redirect.me\"]"`)
}

func TestHandlerServeAppJSWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: RemoteBucket("https://storage.googleapis.com/bhojpur/"),
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `register("/app-worker.js"`)
	require.Contains(t, body, `fetch("https://storage.googleapis.com/bhojpur/web/app.wasm"`)
	require.Contains(t, body, "BHOJPUR_WEB_APP_VERSION")
	require.Contains(t, body, `"BHOJPUR_STATIC_RESOURCES_URL":"https://storage.googleapis.com/bhojpur"`)
	require.Contains(t, body, `"BHOJPUR_WEB_ROOT_PREFIX":""`)
}

func TestHandlerServeAppJSWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: GitHubPages("bhojpur"),
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `register("/bhojpur/app-worker.js"`)
	require.Contains(t, body, `fetch("/bhojpur/web/app.wasm"`)
	require.Contains(t, body, "BHOJPUR_WEB_APP_VERSION")
	require.Contains(t, body, `"BHOJPUR_STATIC_RESOURCES_URL":"/bhojpur"`)
	require.Contains(t, body, `"BHOJPUR_WEB_ROOT_PREFIX":"/bhojpur"`)
}

func TestHandlerServeAppJSWithEnv(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Env: Environment{
			"FOO": "foo",
			"BAR": "bar",
		},
	}
	h.ServeHTTP(w, r)
	body := w.Body.String()

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, "BHOJPUR_WEB_APP_VERSION")
	require.Contains(t, body, `"FOO":"foo"`)
	require.Contains(t, body, `"BAR":"bar"`)
	require.Contains(t, body, `"BHOJPUR_STATIC_RESOURCES_URL":""`)
	require.Contains(t, body, `"BHOJPUR_WEB_ROOT_PREFIX":""`)
}

func TestHandlerServeAppWorkerJSWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Scripts: []string{"web/hello.js"},
		Styles:  []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"/web/hello.css",`)
	require.Contains(t, body, `"/web/hello.js",`)
	require.Contains(t, body, `"/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/wasm_exec.js",`)
	require.Contains(t, body, `"/app.js",`)
	require.Contains(t, body, `"/web/app.wasm",`)
	require.Contains(t, body, `"/",`)
}

func TestHandlerServeAppWorkerJSWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: RemoteBucket("https://storage.googleapis.com/bhojpur/"),
		Scripts:   []string{"web/hello.js"},
		Styles:    []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"https://storage.googleapis.com/bhojpur/web/hello.css",`)
	require.Contains(t, body, `"https://storage.googleapis.com/bhojpur/web/hello.js",`)
	require.Contains(t, body, `"https://storage.googleapis.com/bhojpur/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/wasm_exec.js",`)
	require.Contains(t, body, `"/app.js",`)
	require.Contains(t, body, `"https://storage.googleapis.com/bhojpur/web/app.wasm",`)
	require.Contains(t, body, `"/",`)
}

func TestHandlerServeAppWorkerJSWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app-worker.js", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources: GitHubPages("bhojpur"),
		Scripts:   []string{"web/hello.js"},
		Styles:    []string{"/web/hello.css"},
		CacheableResources: []string{
			"web/hello.png",
			"http://test.io/hello.png",
		},
	}
	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/javascript", w.Header().Get("Content-Type"))
	require.Contains(t, body, `self.addEventListener("install", event => {`)
	require.Contains(t, body, `self.addEventListener("activate", event => {`)
	require.Contains(t, body, `self.addEventListener("fetch", event => {`)
	require.Contains(t, body, `"/bhojpur/web/hello.css",`)
	require.Contains(t, body, `"/bhojpur/web/hello.js",`)
	require.Contains(t, body, `"/bhojpur/web/hello.png",`)
	require.Contains(t, body, `"http://test.io/hello.png",`)
	require.Contains(t, body, `"/bhojpur/wasm_exec.js",`)
	require.Contains(t, body, `"/bhojpur/app.js",`)
	require.Contains(t, body, `"/bhojpur/web/app.wasm",`)
	require.Contains(t, body, `"/bhojpur",`)
}

func TestHandlerServeManifestJSONWithLocalDir(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/"`)
	require.Contains(t, body, `"start_url": "/"`)
}

func TestHandlerServeManifestJSONWithRemoteBucket(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources:       RemoteBucket("https://storage.googleapis.com/bhojpur/"),
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/"`)
	require.Contains(t, body, `"start_url": "/"`)
}

func TestHandlerServeManifestJSONWithGitHubPages(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/manifest.webmanifest", nil)
	w := httptest.NewRecorder()

	h := Handler{
		Resources:       GitHubPages("bhojpur"),
		Name:            "foobar",
		ShortName:       "foo",
		BackgroundColor: "#0000f0",
		ThemeColor:      "#0000ff",
	}

	h.ServeHTTP(w, r)

	body := w.Body.String()
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/manifest+json", w.Header().Get("Content-Type"))
	require.Contains(t, body, `"short_name": "foo"`)
	require.Contains(t, body, `"name": "foobar"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"src": "https://static.bhojpur.net/favicon.ico"`)
	require.Contains(t, body, `"background_color": "#0000f0"`)
	require.Contains(t, body, `"theme_color": "#0000ff"`)
	require.Contains(t, body, `"scope": "/bhojpur/"`)
	require.Contains(t, body, `"start_url": "/bhojpur/"`)
}

func TestHandlerServeAppCSS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/app.css", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "text/css", w.Header().Get("Content-Type"))
	require.Equal(t, appCSS, w.Body.String())
}

func TestHandlerServeAppWasm(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()
	testCreateFile(t, filepath.Join("web", "app.wasm"), "wasm!")

	h := Handler{}
	h.init()

	utests := []struct {
		scenario string
		path     string
	}{
		{
			scenario: "from resource provider path",
			path:     h.Resources.AppWASM(),
		},
		{
			scenario: "from legacy v6 path",
			path:     "/app.wasm",
		},
		{
			scenario: "from legacy v6 path",
			path:     "/bhojpur.wasm",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, u.path, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)
			require.Equal(t, "application/wasm", w.Header().Get("Content-Type"))
			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "wasm!", w.Body.String())
		})
	}
}

func TestHandlerServeFile(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()
	testCreateFile(t, filepath.Join("web", "hello.txt"), "hello!")

	r := httptest.NewRequest(http.MethodGet, "/web/hello.txt", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello!", w.Body.String())
}

func TestHandlerProxyResources(t *testing.T) {
	close := testCreateDir(t, "web")
	defer close()

	s := httptest.NewServer(&Handler{
		ProxyResources: []ProxyResource{
			{
				Path:         "/hello.txt",
				ResourcePath: "/web/hello.txt",
			},
			{
				Path:         "/plop.txt",
				ResourcePath: "/web/plop.txt",
			},
			{
				Path:         "/app.js",
				ResourcePath: "/web/app.js",
			},
		},
	})
	defer s.Close()

	utests := []struct {
		scenario string
		file     string
		body     string
		code     int
		notProxy bool
	}{
		{
			scenario: "robots.txt is fetched",
			file:     "robots.txt",
			code:     http.StatusOK,
			body:     "robots!",
		},
		{
			scenario: "sitemap.xml is fetched",
			file:     "sitemap.xml",
			code:     http.StatusOK,
			body:     "sitemap!",
		},
		{
			scenario: "ads.txt is fetched",
			file:     "ads.txt",
			code:     http.StatusOK,
			body:     "ads!",
		},
		{
			scenario: "proxy resource is fetched",
			file:     "hello.txt",
			code:     http.StatusOK,
			body:     "hello!",
		},
		{
			scenario: "proxy resource is not found",
			file:     "plop.txt",
			code:     http.StatusNotFound,
		},
		{
			scenario: "no proxy resource is not fetched",
			file:     "bye.txt",
			code:     http.StatusNotFound,
			body:     "bye!",
			notProxy: true,
		},
		{
			scenario: "app.js is not a proxy resource",
			file:     "app.js",
			code:     http.StatusOK,
			body:     "wasm!",
			notProxy: true,
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			for i := 0; i < 2; i++ {
				if u.body != "" {
					testCreateFile(t, filepath.Join("web", u.file), u.body)
				}

				url := s.URL + "/" + u.file

				res, err := http.Get(url)
				require.NoError(t, err)
				defer res.Body.Close()

				require.Equal(t, u.code, res.StatusCode)
				if u.code != http.StatusOK {
					return
				}

				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)

				if u.notProxy {
					require.NotEqual(t, u.body, btos(body))
					return
				}

				require.Equal(t, u.body, btos(body))
			}
		})
	}
}

func BenchmarkHandlerColdRun(b *testing.B) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		h := Handler{}
		h.ServeHTTP(w, r)
		h.ServeHTTP(w, r)
	}
}

func BenchmarkHandlerHotRun(b *testing.B) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	h := Handler{}
	h.ServeHTTP(w, r)

	for i := 0; i < b.N; i++ {
		h.ServeHTTP(w, r)
	}
}

func TestIsRemoteLocation(t *testing.T) {
	tests := []struct {
		scenario string
		path     string
		expected bool
	}{
		{
			scenario: "path with http scheme is a remote location",
			path:     "http://localhost/hello",
			expected: true,
		},
		{
			scenario: "path with https scheme is a remote location",
			path:     "https://localhost/hello",
			expected: true,
		},
		{
			scenario: "empty path is not a remote location",
			path:     "",
			expected: false,
		},
		{
			scenario: "working dir path is not a remote location",
			path:     ".",
			expected: false,
		},
		{
			scenario: "absolute path is not a remote location",
			path:     "/User/hello",
			expected: false,
		},
		{
			scenario: "relative path is not a remote location",
			path:     "./hello",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := isRemoteLocation(test.path)
			require.Equal(t, test.expected, res)
		})
	}
}

func TestIsStaticResourcePath(t *testing.T) {
	tests := []struct {
		scenario string
		path     string
		expected bool
	}{
		{
			scenario: "static resource path",
			path:     "/web/hello",
			expected: true,
		},
		{
			scenario: "static resource path with prefix slash",
			path:     "web/hello",
			expected: true,
		},
		{
			scenario: "static resource directory",
			path:     "/web",
			expected: false,
		},
		{
			scenario: "static resource directory without prefix slash",
			path:     "web",
			expected: false,
		},
		{
			scenario: "non static resource",
			path:     "/app.js",
			expected: false,
		},
		{
			scenario: "remote resource",
			path:     "https://localhost/hello",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := isStaticResourcePath(test.path)
			require.Equal(t, test.expected, res)
		})
	}
}
