package analytics

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

import "testing"

func TestAnalytics(t *testing.T) {
	testingProps := func() map[string]interface{} {
		return map[string]interface{}{
			"string": 42,
			"uint":   uint(23),
			"int":    42,
			"float":  42.2,
			"slice":  []interface{}{"hello", 42},
			"map":    map[string]interface{}{"foo": "bar"},
			"struct": struct{ Foo string }{Foo: "bar"},
		}
	}

	providers := []struct {
		name    string
		backend Backend
	}{
		{
			name:    "google analytics",
			backend: NewGoogleAnalytics(),
		},
	}

	for _, p := range providers {
		t.Run(p.name, func(t *testing.T) {
			Add(p.backend)
			defer func() {
				backends = nil
			}()

			t.Run("identify", func(t *testing.T) {
				Identify("Maxoo", nil)
			})
			t.Run("identify with traits", func(t *testing.T) {
				Identify("Maxoo", testingProps())
			})

			t.Run("event", func(t *testing.T) {
				Track("test", nil)
			})
			t.Run("event with properties", func(t *testing.T) {
				Track("test", testingProps())
			})

			t.Run("page", func(t *testing.T) {
				Page("Test", nil)
			})
			t.Run("page with properties", func(t *testing.T) {
				Page("Test", testingProps())
			})
		})
	}
}
