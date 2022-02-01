package log

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
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/client/httplib"
)

func TestFilterChain(t *testing.T) {
	next := func(ctx context.Context, req *httplib.BhojpurHTTPRequest) (*http.Response, error) {
		time.Sleep(100 * time.Millisecond)
		return &http.Response{
			StatusCode: 404,
		}, nil
	}
	builder := NewFilterChainBuilder()
	filter := builder.FilterChain(next)
	req := httplib.Get("https://github.com/notifications?query=repo%3Abhojpur%2Fweb")
	resp, err := filter(context.Background(), req)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
}

func TestContains(t *testing.T) {
	jsonType := "application/json"
	cases := []struct {
		Name        string
		Types       []string
		ContentType string
		Expected    bool
	}{
		{"case1", []string{jsonType}, jsonType, true},
		{"case2", []string{"text/plain"}, jsonType, false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := contains(c.Types, c.ContentType); ans != c.Expected {
				t.Fatalf("Types: %v, ContentType: %v, expected %v, but %v got",
					c.Types, c.ContentType, c.Expected, ans)
			}
		})
	}
}
