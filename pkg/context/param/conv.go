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
	"reflect"

	logsvr "github.com/bhojpur/logger/pkg/engine"
	ctxutl "github.com/bhojpur/web/pkg/context"
)

// ConvertParams converts http method params to values that will be passed to the method controller as arguments
func ConvertParams(methodParams []*MethodParam, methodType reflect.Type, ctx *ctxutl.Context) (result []reflect.Value) {
	result = make([]reflect.Value, 0, len(methodParams))
	for i := 0; i < len(methodParams); i++ {
		reflectValue := convertParam(methodParams[i], methodType.In(i), ctx)
		result = append(result, reflectValue)
	}
	return
}

func convertParam(param *MethodParam, paramType reflect.Type, ctx *ctxutl.Context) (result reflect.Value) {
	paramValue := getParamValue(param, ctx)
	if paramValue == "" {
		if param.required {
			ctx.Abort(400, fmt.Sprintf("Missing parameter %s", param.name))
		} else {
			paramValue = param.defaultValue
		}
	}

	reflectValue, err := parseValue(param, paramValue, paramType)
	if err != nil {
		logsvr.Debug(fmt.Sprintf("Error converting param %s to type %s. Value: %v, Error: %s", param.name, paramType, paramValue, err))
		ctx.Abort(400, fmt.Sprintf("Invalid parameter %s. Can not convert %v to type %s", param.name, paramValue, paramType))
	}

	return reflectValue
}

func getParamValue(param *MethodParam, ctx *ctxutl.Context) string {
	switch param.in {
	case body:
		return string(ctx.Input.RequestBody)
	case header:
		return ctx.Input.Header(param.name)
	case path:
		return ctx.Input.Query(":" + param.name)
	default:
		return ctx.Input.Query(param.name)
	}
}

func parseValue(param *MethodParam, paramValue string, paramType reflect.Type) (result reflect.Value, err error) {
	if paramValue == "" {
		return reflect.Zero(paramType), nil
	}
	parser := getParser(param, paramType)
	value, err := parser.parse(paramValue, paramType)
	if err != nil {
		return result, err
	}

	return safeConvert(reflect.ValueOf(value), paramType)
}

func safeConvert(value reflect.Value, t reflect.Type) (result reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	result = value.Convert(t)
	return
}
