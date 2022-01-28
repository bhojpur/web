package config

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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseConfigure_DefaultBool(t *testing.T) {
	bc := newBaseConfigure("true")
	assert.True(t, bc.DefaultBool("key1", false))
	assert.True(t, bc.DefaultBool("key2", true))
}

func TestBaseConfigure_DefaultFloat(t *testing.T) {
	bc := newBaseConfigure("12.3")
	assert.Equal(t, 12.3, bc.DefaultFloat("key1", 0.1))
	assert.Equal(t, 0.1, bc.DefaultFloat("key2", 0.1))
}

func TestBaseConfigure_DefaultInt(t *testing.T) {
	bc := newBaseConfigure("10")
	assert.Equal(t, 10, bc.DefaultInt("key1", 8))
	assert.Equal(t, 8, bc.DefaultInt("key2", 8))
}

func TestBaseConfigure_DefaultInt64(t *testing.T) {
	bc := newBaseConfigure("64")
	assert.Equal(t, int64(64), bc.DefaultInt64("key1", int64(8)))
	assert.Equal(t, int64(8), bc.DefaultInt64("key2", int64(8)))
}

func TestBaseConfigure_DefaultString(t *testing.T) {
	bc := newBaseConfigure("Hello")
	assert.Equal(t, "Hello", bc.DefaultString("key1", "world"))
	assert.Equal(t, "world", bc.DefaultString("key2", "world"))
}

func TestBaseConfigure_DefaultStrings(t *testing.T) {
	bc := newBaseConfigure("Hello;world")
	assert.Equal(t, []string{"Hello", "world"}, bc.DefaultStrings("key1", []string{"world"}))
	assert.Equal(t, []string{"world"}, bc.DefaultStrings("key2", []string{"world"}))
}

func newBaseConfigure(str1 string) *BaseConfigure {
	return &BaseConfigure{
		reader: func(ctx context.Context, key string) (string, error) {
			if key == "key1" {
				return str1, nil
			} else {
				return "", errors.New("mock error")
			}

		},
	}
}
