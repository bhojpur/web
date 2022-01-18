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
)

// MethodParamOption defines a func which apply options on a MethodParam
type MethodParamOption func(*MethodParam)

// IsRequired indicates that this param is required and can not be omitted from the http request
var IsRequired MethodParamOption = func(p *MethodParam) {
	p.required = true
}

// InHeader indicates that this param is passed via an http header
var InHeader MethodParamOption = func(p *MethodParam) {
	p.in = header
}

// InPath indicates that this param is part of the URL path
var InPath MethodParamOption = func(p *MethodParam) {
	p.in = path
}

// InBody indicates that this param is passed as an http request body
var InBody MethodParamOption = func(p *MethodParam) {
	p.in = body
}

// Default provides a default value for the http param
func Default(defaultValue interface{}) MethodParamOption {
	return func(p *MethodParam) {
		if defaultValue != nil {
			p.defaultValue = fmt.Sprint(defaultValue)
		}
	}
}
