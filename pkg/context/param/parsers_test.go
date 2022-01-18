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
	"reflect"
	"testing"
	"time"
)

type testDefinition struct {
	strValue       string
	expectedValue  interface{}
	expectedParser paramParser
}

func Test_Parsers(t *testing.T) {

	//ints
	checkParser(testDefinition{"1", 1, intParser{}}, t)
	checkParser(testDefinition{"-1", int64(-1), intParser{}}, t)
	checkParser(testDefinition{"1", uint64(1), intParser{}}, t)

	//floats
	checkParser(testDefinition{"1.0", float32(1.0), floatParser{}}, t)
	checkParser(testDefinition{"-1.0", float64(-1.0), floatParser{}}, t)

	//strings
	checkParser(testDefinition{"AB", "AB", stringParser{}}, t)
	checkParser(testDefinition{"AB", []byte{65, 66}, stringParser{}}, t)

	//bools
	checkParser(testDefinition{"true", true, boolParser{}}, t)
	checkParser(testDefinition{"0", false, boolParser{}}, t)

	//timeParser
	checkParser(testDefinition{"2018-03-26T13:54:53Z", time.Date(2018, 3, 26, 13, 54, 53, 0, time.UTC), timeParser{}}, t)
	checkParser(testDefinition{"2018-03-26", time.Date(2018, 3, 26, 0, 0, 0, 0, time.UTC), timeParser{}}, t)

	//json
	checkParser(testDefinition{`{"X": 5, "Y":"Z"}`, struct {
		X int
		Y string
	}{5, "Z"}, jsonParser{}}, t)

	//slice in query is parsed as comma delimited
	checkParser(testDefinition{`1,2`, []int{1, 2}, sliceParser(intParser{})}, t)

	//slice in body is parsed as json
	checkParser(testDefinition{`["a","b"]`, []string{"a", "b"}, jsonParser{}}, t, MethodParam{in: body})

	//pointers
	var someInt = 1
	checkParser(testDefinition{`1`, &someInt, ptrParser(intParser{})}, t)

	var someStruct = struct{ X int }{5}
	checkParser(testDefinition{`{"X": 5}`, &someStruct, jsonParser{}}, t)

}

func checkParser(def testDefinition, t *testing.T, methodParam ...MethodParam) {
	toType := reflect.TypeOf(def.expectedValue)
	var mp MethodParam
	if len(methodParam) == 0 {
		mp = MethodParam{}
	} else {
		mp = methodParam[0]
	}
	parser := getParser(&mp, toType)

	if reflect.TypeOf(parser) != reflect.TypeOf(def.expectedParser) {
		t.Errorf("Invalid parser for value %v. Expected: %v, actual: %v", def.strValue, reflect.TypeOf(def.expectedParser).Name(), reflect.TypeOf(parser).Name())
		return
	}
	result, err := parser.parse(def.strValue, toType)
	if err != nil {
		t.Errorf("Parsing error for value %v. Expected result: %v, error: %v", def.strValue, def.expectedValue, err)
		return
	}
	convResult, err := safeConvert(reflect.ValueOf(result), toType)
	if err != nil {
		t.Errorf("Conversion error for %v. from value: %v, toType: %v, error: %v", def.strValue, result, toType, err)
		return
	}
	if !reflect.DeepEqual(convResult.Interface(), def.expectedValue) {
		t.Errorf("Parsing error for value %v. Expected result: %v, actual: %v", def.strValue, def.expectedValue, result)
	}
}
