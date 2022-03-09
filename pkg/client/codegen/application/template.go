package application

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
	"strings"
	"unicode/utf8"

	"github.com/bhojpur/web/pkg/client/utils"
	"github.com/bhojpur/web/pkg/template"
)

func init() {
	_ = template.RegisterFilter("lowerFirst", templateLowerFirst)
	_ = template.RegisterFilter("upperFirst", templateUpperFirst)
	_ = template.RegisterFilter("snakeString", templateSnakeString)
	_ = template.RegisterFilter("camelString", templateCamelString)
}

func templateLowerFirst(in *template.Value, param *template.Value) (*template.Value, *template.Error) {
	if in.Len() <= 0 {
		return template.AsValue(""), nil
	}
	t := in.String()
	r, size := utf8.DecodeRuneInString(t)
	return template.AsValue(strings.ToLower(string(r)) + t[size:]), nil
}

func templateUpperFirst(in *template.Value, param *template.Value) (*template.Value, *template.Error) {
	if in.Len() <= 0 {
		return template.AsValue(""), nil
	}
	t := in.String()
	return template.AsValue(strings.Replace(t, string(t[0]), strings.ToUpper(string(t[0])), 1)), nil
}

// snake string, XxYy to xx_yy
func templateSnakeString(in *template.Value, param *template.Value) (*template.Value, *template.Error) {
	if in.Len() <= 0 {
		return template.AsValue(""), nil
	}
	t := in.String()
	return template.AsValue(utils.SnakeString(t)), nil
}

// snake string, XxYy to xx_yy
func templateCamelString(in *template.Value, param *template.Value) (*template.Value, *template.Error) {
	if in.Len() <= 0 {
		return template.AsValue(""), nil
	}
	t := in.String()
	return template.AsValue(utils.CamelString(t)), nil
}

//func upperFirst(str string) string {
//	return strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
//}

func lowerFirst(str string) string {
	return strings.Replace(str, string(str[0]), strings.ToLower(string(str[0])), 1)
}
