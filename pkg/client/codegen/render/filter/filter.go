package filter

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
	"reflect"
	"strings"

	"github.com/bhojpur/web/pkg/template"
)

func init() {
	template.RegisterFilter("ValueWithMap", valueWithMap)
	template.RegisterFilter("HasPrefix", hasPrefix)
	template.RegisterFilter("HasSuffix", hasSuffix)
	template.RegisterFilter("CompareString", compareString)
}

////////////////////////////////////////////////////////////////////////////////
func valueWithMap(in, param *template.Value) (*template.Value, *template.Error) {
	var source = in.Interface()
	var key = param.Interface()

	if source == nil {
		return nil, nil
	}

	if key == nil {
		return nil, nil
	}

	var sourceValue = reflect.ValueOf(source)
	if sourceValue.IsNil() {
		return nil, nil
	}

	switch sourceValue.Kind() {
	case reflect.Map:
		var targetValue = reflect.ValueOf(key)
		if targetValue.IsValid() {
			return template.AsValue(sourceValue.MapIndex(targetValue).Interface()), nil
		}
	}
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
func hasPrefix(in, param *template.Value) (*template.Value, *template.Error) {
	return template.AsValue(strings.HasPrefix(in.String(), param.String())), nil
}

func hasSuffix(in, param *template.Value) (*template.Value, *template.Error) {
	return template.AsValue(strings.HasSuffix(in.String(), param.String())), nil
}

func compareString(in, param *template.Value) (*template.Value, *template.Error) {
	return template.AsValue(strings.Compare(in.String(), param.String())), nil
}
