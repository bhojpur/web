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
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestPage(t *testing.T) {
	testPage(t, &requestPage{
		width:  42,
		height: 21,
	})
}

func TestBrowserPage(t *testing.T) {
	testSkipNonWasm(t)
	testPage(t, browserPage{})
}

func testPage(t *testing.T, p Page) {
	p.SetTitle("Bhojpur Web")
	require.Equal(t, "Bhojpur Web", p.Title())

	p.SetLang("fr")
	require.Equal(t, "fr", p.Lang())

	p.SetDescription("test")
	require.Equal(t, "test", p.Description())

	p.SetAuthor("Pramila")
	require.Equal(t, "Pramila", p.Author())

	p.SetKeywords("go", "app")
	require.Equal(t, "go, app", p.Keywords())

	p.SetLoadingLabel("loading test")

	p.SetImage("image")
	require.Equal(t, "image", p.Image())

	u, _ := url.Parse("https://bhojpur.net")
	p.ReplaceURL(u)
	require.Equal(t, u.String(), p.URL().String())

	w, h := p.Size()
	require.NotZero(t, w)
	require.NotZero(t, h)
}
