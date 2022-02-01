package bean

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagAutoWireBeanFactory_AutoWire(t *testing.T) {
	factory := NewTagAutoWireBeanFactory()
	bm := &ComplicateStruct{}
	err := factory.AutoWire(context.Background(), nil, bm)
	assert.Nil(t, err)
	assert.Equal(t, 12, bm.IntValue)
	assert.Equal(t, "hello, strValue", bm.StrValue)

	assert.Equal(t, int8(8), bm.Int8Value)
	assert.Equal(t, int16(16), bm.Int16Value)
	assert.Equal(t, int32(32), bm.Int32Value)
	assert.Equal(t, int64(64), bm.Int64Value)

	assert.Equal(t, uint(13), bm.UintValue)
	assert.Equal(t, uint8(88), bm.Uint8Value)
	assert.Equal(t, uint16(1616), bm.Uint16Value)
	assert.Equal(t, uint32(3232), bm.Uint32Value)
	assert.Equal(t, uint64(6464), bm.Uint64Value)

	assert.Equal(t, float32(32.32), bm.Float32Value)
	assert.Equal(t, float64(64.64), bm.Float64Value)

	assert.True(t, bm.BoolValue)
	assert.Equal(t, 0, bm.ignoreInt)

	assert.NotNil(t, bm.TimeValue)
}

type ComplicateStruct struct {
	IntValue   int    `default:"12"`
	StrValue   string `default:"hello, strValue"`
	Int8Value  int8   `default:"8"`
	Int16Value int16  `default:"16"`
	Int32Value int32  `default:"32"`
	Int64Value int64  `default:"64"`

	UintValue   uint   `default:"13"`
	Uint8Value  uint8  `default:"88"`
	Uint16Value uint16 `default:"1616"`
	Uint32Value uint32 `default:"3232"`
	Uint64Value uint64 `default:"6464"`

	Float32Value float32 `default:"32.32"`
	Float64Value float64 `default:"64.64"`

	BoolValue bool `default:"true"`

	ignoreInt int `default:"11"`

	TimeValue time.Time `default:"2018-03-26 12:13:14.000"`
}
