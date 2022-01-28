package yaml2

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
	"regexp"
	"strconv"
)

var (
	RE_INT, _   = regexp.Compile("^[0-9,]+$")
	RE_FLOAT, _ = regexp.Compile("^[0-9]+[.][0-9]+$")
	RE_DATE, _  = regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
	RE_TIME, _  = regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$")
)

func string2Val(str string) interface{} {
	tmp := []byte(str)
	switch {
	case str == "false":
		return false
	case str == "true":
		return true
	case RE_INT.Match(tmp):
		// TODO check err
		_int, _ := strconv.ParseInt(str, 10, 64)
		return _int
	case RE_FLOAT.Match(tmp):
		_float, _ := strconv.ParseFloat(str, 64)
		return _float
		//TODO support time or Not?
		/*
			case RE_DATE.Match(tmp):
				_date, _ := time.Parse("2018-03-26", str)
				return _date
			case RE_TIME.Match(tmp):
				_time, _ := time.Parse("2018-03-26 03:04:05", str)
				return _time
		*/
	}
	return str
}
