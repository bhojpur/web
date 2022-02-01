package orm

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
	"testing"
)

func TestCamelString(t *testing.T) {
	snake := []string{"pic_url", "hello_world_", "hello__World", "_HelLO_Word", "pic_url_1", "pic_url__1"}
	camel := []string{"PicUrl", "HelloWorld", "HelloWorld", "HelLOWord", "PicUrl1", "PicUrl1"}

	answer := make(map[string]string)
	for i, v := range snake {
		answer[v] = camel[i]
	}

	for _, v := range snake {
		res := camelString(v)
		if res != answer[v] {
			t.Error("Unit Test Fail:", v, res, answer[v])
		}
	}
}

func TestSnakeString(t *testing.T) {
	camel := []string{"PicUrl", "HelloWorld", "HelloWorld", "HelLOWord", "PicUrl1", "XyXX"}
	snake := []string{"pic_url", "hello_world", "hello_world", "hel_l_o_word", "pic_url1", "xy_x_x"}

	answer := make(map[string]string)
	for i, v := range camel {
		answer[v] = snake[i]
	}

	for _, v := range camel {
		res := snakeString(v)
		if res != answer[v] {
			t.Error("Unit Test Fail:", v, res, answer[v])
		}
	}
}

func TestSnakeStringWithAcronym(t *testing.T) {
	camel := []string{"ID", "PicURL", "HelloWorld", "HelloWorld", "HelLOWord", "PicUrl1", "XyXX"}
	snake := []string{"id", "pic_url", "hello_world", "hello_world", "hel_lo_word", "pic_url1", "xy_xx"}

	answer := make(map[string]string)
	for i, v := range camel {
		answer[v] = snake[i]
	}

	for _, v := range camel {
		res := snakeStringWithAcronym(v)
		if res != answer[v] {
			t.Error("Unit Test Fail:", v, res, answer[v])
		}
	}
}
