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
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/bhojpur/web/pkg/app/errors"
)

func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v

	case []byte:
		return btos(v)

	case int:
		return strconv.Itoa(v)

	case float64:
		return strconv.FormatFloat(v, 'f', 4, 64)

	case bool:
		return strconv.FormatBool(v)

	case nil:
		return ""

	default:
		return fmt.Sprint(v)
	}
}

func toPath(v ...interface{}) string {
	var b strings.Builder

	for _, o := range v {
		s := toString(o)
		if s == "" {
			continue
		}
		b.WriteByte('/')
		b.WriteString(s)
	}

	return b.String()
}

func writeIndent(w io.Writer, indent int) {
	for i := 0; i < indent*2; i++ {
		w.Write([]byte(" "))
	}
}

func ln() []byte {
	return []byte("\n")
}

func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func pxToString(px int) string {
	return strconv.Itoa(px) + "px"
}

func stringTo(s string, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return errors.New("receiver in not a pointer").Tag("receiver-type", val.Type())
	}
	val = val.Elem()

	switch val.Kind() {
	case reflect.String:
		val.SetString(s)

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, _ := strconv.ParseInt(s, 10, 0)
		val.SetInt(i)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, _ := strconv.ParseUint(s, 10, 0)
		val.SetUint(i)

	case reflect.Float64:
		f, _ := strconv.ParseFloat(s, 64)
		val.SetFloat(f)

	case reflect.Float32:
		f, _ := strconv.ParseFloat(s, 32)
		val.SetFloat(f)

	default:
		return errors.New("string cannot be converted to receiver type").
			Tag("string", s).
			Tag("receiver-type", val.Type())
	}

	return nil
}

// AppendClass adds c to the given class string.
func AppendClass(class, c string) string {
	if c == "" {
		return class
	}
	if class != "" {
		class += " "
	}
	class += c
	return class
}
