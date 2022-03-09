//go:build !wasm

package jsutil

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

import "github.com/bhojpur/web/pkg/app"

// JS2Bytes convert from TypedArray for JS to byte slice for Go.
func JS2Bytes(v app.Value) []byte {
	panic("not implemented")
	return nil
}

// Bytes2JS convert from byte slice for Go to Uint8Array for app.
func Bytes2JS(b []byte) app.Value {
	panic("not implemented")
	return nil
}

// Callback0 make auto-release callback without params.
func Callback0(fn func() interface{}) app.Func {
	panic("not implemented")
	return nil
}

// Callback1 make auto-release callback with 1 param.
func Callback1(fn func(res app.Value) interface{}) app.Func {
	panic("not implemented")
	return nil
}

// CallbackN make auto-release callback with multiple params.
func CallbackN(fn func(res []app.Value) interface{}) app.Func {
	panic("not implemented")
	return nil
}

// WrapError to wrap golang standard error to app.Value.
func WrapError(err error) app.Value {
	panic("not implemented")
	return nil
}

// UnwrapError to unwrap app.Value to golang standard error.
func UnwrapError(v app.Value) error {
	panic("not implemented")
	return nil
}

// IsArray checking value is array type.
func IsArray(item app.Value) bool {
	panic("not implemented")
	return false
}

// JS2Go JS values convert to Go values.
func JS2Go(obj app.Value) interface{} {
	panic("not implemented")
	return nil
}

// Form2Go retrieve form values from form element.
func Form2Go(form app.Value) map[string]interface{} {
	panic("not implemented")
	return nil
}
