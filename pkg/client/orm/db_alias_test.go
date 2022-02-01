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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDataBase(t *testing.T) {
	err := RegisterDataBase("test-params", DBARGS.Driver, DBARGS.Source,
		MaxIdleConnections(20),
		MaxOpenConnections(300),
		ConnMaxLifetime(time.Minute))
	assert.Nil(t, err)

	al := getDbAlias("test-params")
	assert.NotNil(t, al)
	assert.Equal(t, al.MaxIdleConns, 20)
	assert.Equal(t, al.MaxOpenConns, 300)
	assert.Equal(t, al.ConnMaxLifetime, time.Minute)
}

func TestRegisterDataBase_MaxStmtCacheSizeNegative1(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSizeNegative1"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(-1))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBase_MaxStmtCacheSize0(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize0"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(0))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBase_MaxStmtCacheSize1(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize1"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(1))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 1)
}

func TestRegisterDataBase_MaxStmtCacheSize841(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize841"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(841))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 841)
}

func TestDBCache(t *testing.T) {
	dataBaseCache.add("test1", &alias{})
	dataBaseCache.add("default", &alias{})
	al := dataBaseCache.getDefault()
	assert.NotNil(t, al)
	al, ok := dataBaseCache.get("test1")
	assert.NotNil(t, al)
	assert.True(t, ok)
}
