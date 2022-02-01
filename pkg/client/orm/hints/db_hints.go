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
	"github.com/bhojpur/web/pkg/core/utils"
)

const (
	//query level
	KeyForceIndex = iota
	KeyUseIndex
	KeyIgnoreIndex
	KeyForUpdate
	KeyLimit
	KeyOffset
	KeyOrderBy
	KeyRelDepth
)

type Hint struct {
	key   interface{}
	value interface{}
}

var _ utils.KV = new(Hint)

// GetKey return key
func (s *Hint) GetKey() interface{} {
	return s.key
}

// GetValue return value
func (s *Hint) GetValue() interface{} {
	return s.value
}

var _ utils.KV = new(Hint)

// ForceIndex return a hint about ForceIndex
func ForceIndex(indexes ...string) *Hint {
	return NewHint(KeyForceIndex, indexes)
}

// UseIndex return a hint about UseIndex
func UseIndex(indexes ...string) *Hint {
	return NewHint(KeyUseIndex, indexes)
}

// IgnoreIndex return a hint about IgnoreIndex
func IgnoreIndex(indexes ...string) *Hint {
	return NewHint(KeyIgnoreIndex, indexes)
}

// ForUpdate return a hint about ForUpdate
func ForUpdate() *Hint {
	return NewHint(KeyForUpdate, true)
}

// DefaultRelDepth return a hint about DefaultRelDepth
func DefaultRelDepth() *Hint {
	return NewHint(KeyRelDepth, true)
}

// RelDepth return a hint about RelDepth
func RelDepth(d int) *Hint {
	return NewHint(KeyRelDepth, d)
}

// Limit return a hint about Limit
func Limit(d int64) *Hint {
	return NewHint(KeyLimit, d)
}

// Offset return a hint about Offset
func Offset(d int64) *Hint {
	return NewHint(KeyOffset, d)
}

// OrderBy return a hint about OrderBy
func OrderBy(s string) *Hint {
	return NewHint(KeyOrderBy, s)
}

// NewHint return a hint
func NewHint(key interface{}, value interface{}) *Hint {
	return &Hint{
		key:   key,
		value: value,
	}
}
