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
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type paramParser interface {
	parse(value string, toType reflect.Type) (interface{}, error)
}

func getParser(param *MethodParam, t reflect.Type) paramParser {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return intParser{}
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 { // treat []byte as string
			return stringParser{}
		}
		if param.in == body {
			return jsonParser{}
		}
		elemParser := getParser(param, t.Elem())
		if elemParser == (jsonParser{}) {
			return elemParser
		}
		return sliceParser(elemParser)
	case reflect.Bool:
		return boolParser{}
	case reflect.String:
		return stringParser{}
	case reflect.Float32, reflect.Float64:
		return floatParser{}
	case reflect.Ptr:
		elemParser := getParser(param, t.Elem())
		if elemParser == (jsonParser{}) {
			return elemParser
		}
		return ptrParser(elemParser)
	default:
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return timeParser{}
		}
		return jsonParser{}
	}
}

type parserFunc func(value string, toType reflect.Type) (interface{}, error)

func (f parserFunc) parse(value string, toType reflect.Type) (interface{}, error) {
	return f(value, toType)
}

type boolParser struct{}

func (p boolParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return strconv.ParseBool(value)
}

type stringParser struct{}

func (p stringParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return value, nil
}

type intParser struct{}

func (p intParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return strconv.Atoi(value)
}

type floatParser struct{}

func (p floatParser) parse(value string, toType reflect.Type) (interface{}, error) {
	if toType.Kind() == reflect.Float32 {
		res, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(res), nil
	}
	return strconv.ParseFloat(value, 64)
}

type timeParser struct{}

func (p timeParser) parse(value string, toType reflect.Type) (result interface{}, err error) {
	result, err = time.Parse(time.RFC3339, value)
	if err != nil {
		result, err = time.Parse("2006-01-02", value)
	}
	return
}

type jsonParser struct{}

func (p jsonParser) parse(value string, toType reflect.Type) (interface{}, error) {
	pResult := reflect.New(toType)
	v := pResult.Interface()
	err := json.Unmarshal([]byte(value), v)
	if err != nil {
		return nil, err
	}
	return pResult.Elem().Interface(), nil
}

func sliceParser(elemParser paramParser) paramParser {
	return parserFunc(func(value string, toType reflect.Type) (interface{}, error) {
		values := strings.Split(value, ",")
		result := reflect.MakeSlice(toType, 0, len(values))
		elemType := toType.Elem()
		for _, v := range values {
			parsedValue, err := elemParser.parse(v, elemType)
			if err != nil {
				return nil, err
			}
			result = reflect.Append(result, reflect.ValueOf(parsedValue))
		}
		return result.Interface(), nil
	})
}

func ptrParser(elemParser paramParser) paramParser {
	return parserFunc(func(value string, toType reflect.Type) (interface{}, error) {
		parsedValue, err := elemParser.parse(value, toType.Elem())
		if err != nil {
			return nil, err
		}
		newValPtr := reflect.New(toType.Elem())
		newVal := reflect.Indirect(newValPtr)
		convertedVal, err := safeConvert(reflect.ValueOf(parsedValue), toType.Elem())
		if err != nil {
			return nil, err
		}

		newVal.Set(convertedVal)
		return newValPtr.Interface(), nil
	})
}
