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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClientStaticResourceResolver(t *testing.T) {
	utests := []struct {
		scenario           string
		staticResourcesURL string
		path               string
		expected           string
	}{
		{
			scenario: "non-static resource is skipped",
			path:     "/hello",
			expected: "/hello",
		},
		{
			scenario: "non-static resource without slash is skipped",
			path:     "hello",
			expected: "hello",
		},
		{
			scenario:           "non-static resource with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/bhojpur",
			path:               "/hello",
			expected:           "/hello",
		},
		{
			scenario:           "non-static resource without slash and with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/bhojpur",
			path:               "hello",
			expected:           "hello",
		},
		{
			scenario: "static resource",
			path:     "/web/hello.css",
			expected: "/web/hello.css",
		},
		{
			scenario: "static resource without slash",
			path:     "web/hello.css",
			expected: "/web/hello.css",
		},
		{
			scenario:           "static resource with remote root dir is resolved",
			staticResourcesURL: "https://storage.googleapis.com/bhojpur",
			path:               "/web/hello.css",
			expected:           "https://storage.googleapis.com/bhojpur/web/hello.css",
		},
		{
			scenario:           "static resource without slash and with remote root dir is resolved",
			staticResourcesURL: "https://storage.googleapis.com/bhojpur",
			path:               "web/hello.css",
			expected:           "https://storage.googleapis.com/bhojpur/web/hello.css",
		},
		{
			scenario: "resolved static resource is skipped",
			path:     "https://storage.googleapis.com/bhojpur/web/hello.css",
			expected: "https://storage.googleapis.com/bhojpur/web/hello.css",
		},
		{
			scenario:           "resolved static resource with remote root dir is skipped",
			staticResourcesURL: "https://storage.googleapis.com/bhojpur",
			path:               "https://storage.googleapis.com/bhojpur/web/hello.css",
			expected:           "https://storage.googleapis.com/bhojpur/web/hello.css",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			res := newClientStaticResourceResolver(u.staticResourcesURL)(u.path)
			require.Equal(t, u.expected, res)
		})
	}
}
