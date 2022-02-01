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

	"github.com/stretchr/testify/assert"
)

type Interface struct {
	Id   int
	Name string

	Index1 string
	Index2 string

	Unique1 string
	Unique2 string
}

func (i *Interface) TableIndex() [][]string {
	return [][]string{{"index1"}, {"index2"}}
}

func (i *Interface) TableUnique() [][]string {
	return [][]string{{"unique1"}, {"unique2"}}
}

func (i *Interface) TableName() string {
	return "INTERFACE_"
}

func (i *Interface) TableEngine() string {
	return "innodb"
}

func TestDbBase_GetTables(t *testing.T) {
	RegisterModel(&Interface{})
	mi, ok := modelCache.get("INTERFACE_")
	assert.True(t, ok)
	assert.NotNil(t, mi)

	engine := getTableEngine(mi.addrField)
	assert.Equal(t, "innodb", engine)
	uniques := getTableUnique(mi.addrField)
	assert.Equal(t, [][]string{{"unique1"}, {"unique2"}}, uniques)
	indexes := getTableIndex(mi.addrField)
	assert.Equal(t, [][]string{{"index1"}, {"index2"}}, indexes)
}
