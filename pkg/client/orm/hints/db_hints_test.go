package hints

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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHint_time(t *testing.T) {
	key := "qweqwe"
	value := time.Second
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestNewHint_int(t *testing.T) {
	key := "qweqwe"
	value := 281230
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestNewHint_float(t *testing.T) {
	key := "qweqwe"
	value := 21.2459753
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestForceIndex(t *testing.T) {
	s := []string{`f_index1`, `f_index2`, `f_index3`}
	hint := ForceIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyForceIndex)
}

func TestForceIndex_0(t *testing.T) {
	var s []string
	hint := ForceIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyForceIndex)
}

func TestIgnoreIndex(t *testing.T) {
	s := []string{`i_index1`, `i_index2`, `i_index3`}
	hint := IgnoreIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyIgnoreIndex)
}

func TestIgnoreIndex_0(t *testing.T) {
	var s []string
	hint := IgnoreIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyIgnoreIndex)
}

func TestUseIndex(t *testing.T) {
	s := []string{`u_index1`, `u_index2`, `u_index3`}
	hint := UseIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyUseIndex)
}

func TestUseIndex_0(t *testing.T) {
	var s []string
	hint := UseIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyUseIndex)
}

func TestForUpdate(t *testing.T) {
	hint := ForUpdate()
	assert.Equal(t, hint.GetValue(), true)
	assert.Equal(t, hint.GetKey(), KeyForUpdate)
}

func TestDefaultRelDepth(t *testing.T) {
	hint := DefaultRelDepth()
	assert.Equal(t, hint.GetValue(), true)
	assert.Equal(t, hint.GetKey(), KeyRelDepth)
}

func TestRelDepth(t *testing.T) {
	hint := RelDepth(157965)
	assert.Equal(t, hint.GetValue(), 157965)
	assert.Equal(t, hint.GetKey(), KeyRelDepth)
}

func TestLimit(t *testing.T) {
	hint := Limit(1579625)
	assert.Equal(t, hint.GetValue(), int64(1579625))
	assert.Equal(t, hint.GetKey(), KeyLimit)
}

func TestOffset(t *testing.T) {
	hint := Offset(int64(1572123965))
	assert.Equal(t, hint.GetValue(), int64(1572123965))
	assert.Equal(t, hint.GetKey(), KeyOffset)
}

func TestOrderBy(t *testing.T) {
	hint := OrderBy(`-ID`)
	assert.Equal(t, hint.GetValue(), `-ID`)
	assert.Equal(t, hint.GetKey(), KeyOrderBy)
}
