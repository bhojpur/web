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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestIni(t *testing.T) {

	var (
		inicontext = `
;comment one
#comment two
appname = bhojpurapi
httpport = 8080
mysqlport = 3600
PI = 3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
session= on
cookieon= off
newreg = OFF
needlogin = ON
enableSession = Y
enableCookie = N
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

		keyValue = map[string]interface{}{
			"appname":               "bhojpurapi",
			"httpport":              8080,
			"mysqlport":             int64(3600),
			"pi":                    3.1415976,
			"runmode":               "dev",
			"autorender":            false,
			"copyrequestbody":       true,
			"session":               true,
			"cookieon":              false,
			"newreg":                false,
			"needlogin":             true,
			"enableSession":         true,
			"enableCookie":          false,
			"flag":                  true,
			"path1":                 os.Getenv("GOPATH"),
			"path2":                 os.Getenv("GOPATH"),
			"demo::key1":            "bhoj",
			"demo::key2":            "pur",
			"demo::CaseInsensitive": true,
			"demo::peers":           []string{"one", "two", "three"},
			"demo::password":        os.Getenv("GOPATH"),
			"null":                  "",
			"demo2::key1":           "",
			"error":                 "",
			"emptystrings":          []string{},
		}
	)

	f, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(inicontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testini.conf")
	iniconf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = iniconf.Int(k)
		case int64:
			value, err = iniconf.Int64(k)
		case float64:
			value, err = iniconf.Float(k)
		case bool:
			value, err = iniconf.Bool(k)
		case []string:
			value = iniconf.Strings(k)
		case string:
			value = iniconf.String(k)
		default:
			value, err = iniconf.DIY(k)
		}
		if err != nil {
			t.Fatalf("get key %q value fail,err %s", k, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Fatalf("get key %q value, want %v got %v .", k, v, value)
		}

	}
	if err = iniconf.Set("name", "bhojpur"); err != nil {
		t.Fatal(err)
	}
	if iniconf.String("name") != "bhojpur" {
		t.Fatal("get name error")
	}

}

func TestIniSave(t *testing.T) {

	const (
		inicontext = `
app = app
;comment one
#comment two
# comment three
appname = bhojpurapi
httpport = 8080
# DB Info
# enable db
[dbinfo]
# db type name
# suport mysql,sqlserver
name = mysql
`

		saveResult = `
app=app
#comment one
#comment two
# comment three
appname=bhojpurapi
httpport=8080

# DB Info
# enable db
[dbinfo]
# db type name
# suport mysql,sqlserver
name=mysql
`
	)
	cfg, err := NewConfigData("ini", []byte(inicontext))
	if err != nil {
		t.Fatal(err)
	}
	name := "newIniConfig.ini"
	if err := cfg.SaveConfigFile(name); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)

	if data, err := ioutil.ReadFile(name); err != nil {
		t.Fatal(err)
	} else {
		cfgData := string(data)
		datas := strings.Split(saveResult, "\n")
		for _, line := range datas {
			if !strings.Contains(cfgData, line+"\n") {
				t.Fatalf("different after save ini config file. need contains %q", line)
			}
		}

	}
}
