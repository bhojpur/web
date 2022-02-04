package param

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
	"strings"
)

// MethodParam keeps param information to be auto passed to controller methods
type MethodParam struct {
	name         string
	in           paramType
	required     bool
	defaultValue string
}

type paramType byte

const (
	param paramType = iota
	path
	body
	header
)

// New creates a new MethodParam with name and specific options
func New(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, nil, opts)
}

func newParam(name string, parser paramParser, opts []MethodParamOption) (param *MethodParam) {
	param = &MethodParam{name: name}
	for _, option := range opts {
		option(param)
	}
	return
}

// Make creates an array of MethodParmas or an empty array
func Make(list ...*MethodParam) []*MethodParam {
	if len(list) > 0 {
		return list
	}
	return nil
}

func (mp *MethodParam) String() string {
	options := []string{}
	result := "param.New(\"" + mp.name + "\""
	if mp.required {
		options = append(options, "param.IsRequired")
	}
	switch mp.in {
	case path:
		options = append(options, "param.InPath")
	case body:
		options = append(options, "param.InBody")
	case header:
		options = append(options, "param.InHeader")
	}
	if mp.defaultValue != "" {
		options = append(options, fmt.Sprintf(`param.Default("%s")`, mp.defaultValue))
	}
	if len(options) > 0 {
		result += ", "
	}
	result += strings.Join(options, ", ")
	result += ")"
	return result
}
