package common

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
	webapp "github.com/bhojpur/web/pkg/app"
	"github.com/bhojpur/web/pkg/app/errors"
)

const (
	getMarkdown = "/markdown/get"
)

func handleGetMarkdown(ctx webapp.Context, a webapp.Action) {
	path := a.Tags.Get("path")
	if path == "" {
		webapp.Log(errors.New("getting markdown failed").
			Tag("reason", "empty path"))
		return
	}
	state := markdownState(path)

	var md markdownContent
	ctx.GetState(state, &md)
	switch md.Status {
	case loading, loaded:
		return
	}

	md.Status = loading
	md.Err = nil
	ctx.SetState(state, md)

	res, err := get(ctx, path)
	if err != nil {
		md.Status = loadingErr
		md.Err = errors.New("getting markdown failed").Wrap(err)
		ctx.SetState(state, md)
		return
	}

	md.Status = loaded
	md.Data = string(res)
	ctx.SetState(state, md)
}

func markdownState(src string) string {
	return src
}

type markdownContent struct {
	Status status
	Err    error
	Data   string
}
