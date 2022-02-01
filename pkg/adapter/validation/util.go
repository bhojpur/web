package validation

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

	"github.com/bhojpur/web/pkg/core/validation"
)

const (
	// ValidTag struct tag
	ValidTag = validation.ValidTag

	LabelTag = validation.LabelTag
)

var (
	ErrInt64On32 = validation.ErrInt64On32
)

// CustomFunc is for custom validate function
type CustomFunc func(v *Validation, obj interface{}, key string)

// AddCustomFunc Add a custom function to validation
// The name can not be:
//   Clear
//   HasErrors
//   ErrorMap
//   Error
//   Check
//   Valid
//   NoMatch
// If the name is same with exists function, it will replace the origin valid function
func AddCustomFunc(name string, f CustomFunc) error {
	return validation.AddCustomFunc(name, func(v *validation.Validation, obj interface{}, key string) {
		f((*Validation)(v), obj, key)
	})
}

// ValidFunc Valid function type
type ValidFunc validation.ValidFunc

// Funcs Validate function map
type Funcs validation.Funcs

// Call validate values with named type string
func (f Funcs) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	return (validation.Funcs(f)).Call(name, params...)
}
