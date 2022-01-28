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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalInstance(t *testing.T) {
	cfgStr := `
appname = bhojpurapi
httpport = 8080
mysqlport = 3600
PI = 3.1415926
runmode = "dev"
autorender = false
copyrequestbody = true
session= on
cookieon= off
newreg = OFF
needlogin = ON
enableSession = Y
enableCookie = N
developer="hari;krishna"
flag = 1
path1 = ${GOPATH}
path2 = ${GOPATH||/home/go}
[demo]
key1="bhoj"
key2 = "pur"
CaseInsensitive = true
peers = one;two;three
password = ${GOPATH}
`
	path := os.TempDir() + string(os.PathSeparator) + "test_global_instance.ini"
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(cfgStr)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(path)

	err = InitGlobalInstance("ini", path)
	assert.Nil(t, err)

	val, err := String("appname")
	assert.Nil(t, err)
	assert.Equal(t, "bhojpurapi", val)

	val = DefaultString("appname__", "404")
	assert.Equal(t, "404", val)

	vi, err := Int("httpport")
	assert.Nil(t, err)
	assert.Equal(t, 8080, vi)
	vi = DefaultInt("httpport__", 404)
	assert.Equal(t, 404, vi)

	vi64, err := Int64("mysqlport")
	assert.Nil(t, err)
	assert.Equal(t, int64(3600), vi64)
	vi64 = DefaultInt64("mysqlport__", 404)
	assert.Equal(t, int64(404), vi64)

	vf, err := Float("PI")
	assert.Nil(t, err)
	assert.Equal(t, 3.1415926, vf)
	vf = DefaultFloat("PI__", 4.04)
	assert.Equal(t, 4.04, vf)

	vb, err := Bool("copyrequestbody")
	assert.Nil(t, err)
	assert.True(t, vb)

	vb = DefaultBool("copyrequestbody__", false)
	assert.False(t, vb)

	vss := DefaultStrings("developer__", []string{"hari", ""})
	assert.Equal(t, []string{"hari", ""}, vss)

	vss, err = Strings("developer")
	assert.Nil(t, err)
	assert.Equal(t, []string{"hari", "krishna"}, vss)
}
