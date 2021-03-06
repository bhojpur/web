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
	"os"
	"testing"
)

func TestJsonStartsWithArray(t *testing.T) {

	const jsoncontextwitharray = `[
	{
		"url": "user",
		"serviceAPI": "http://www.test.com/user"
	},
	{
		"url": "employee",
		"serviceAPI": "http://www.test.com/employee"
	}
]`
	f, err := os.Create("testjsonWithArray.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontextwitharray)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testjsonWithArray.conf")
	jsonconf, err := NewConfig("json", "testjsonWithArray.conf")
	if err != nil {
		t.Fatal(err)
	}
	rootArray, err := jsonconf.DIY("rootArray")
	if err != nil {
		t.Error("array does not exist as element")
	}
	rootArrayCasted := rootArray.([]interface{})
	if rootArrayCasted == nil {
		t.Error("array from root is nil")
	} else {
		elem := rootArrayCasted[0].(map[string]interface{})
		if elem["url"] != "user" || elem["serviceAPI"] != "http://www.test.com/user" {
			t.Error("array[0] values are not valid")
		}

		elem2 := rootArrayCasted[1].(map[string]interface{})
		if elem2["url"] != "employee" || elem2["serviceAPI"] != "http://www.test.com/employee" {
			t.Error("array[1] values are not valid")
		}
	}
}

func TestJson(t *testing.T) {

	var (
		jsoncontext = `{
"appname": "bhojpurapi",
"testnames": "foo;bar",
"httpport": 8080,
"mysqlport": 3600,
"PI": 3.1415976, 
"runmode": "dev",
"autorender": false,
"copyrequestbody": true,
"session": "on",
"cookieon": "off",
"newreg": "OFF",
"needlogin": "ON",
"enableSession": "Y",
"enableCookie": "N",
"flag": 1,
"path1": "${GOPATH}",
"path2": "${GOPATH||/home/go}",
"database": {
        "host": "host",
        "port": "port",
        "database": "database",
        "username": "username",
        "password": "${GOPATH}",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info",
			"root": "${GOPATH}"
		}
    }
}`
		keyValue = map[string]interface{}{
			"appname":                         "bhojpurapi",
			"testnames":                       []string{"foo", "bar"},
			"httpport":                        8080,
			"mysqlport":                       int64(3600),
			"PI":                              3.1415976,
			"runmode":                         "dev",
			"autorender":                      false,
			"copyrequestbody":                 true,
			"session":                         true,
			"cookieon":                        false,
			"newreg":                          false,
			"needlogin":                       true,
			"enableSession":                   true,
			"enableCookie":                    false,
			"flag":                            true,
			"path1":                           os.Getenv("GOPATH"),
			"path2":                           os.Getenv("GOPATH"),
			"database::host":                  "host",
			"database::port":                  "port",
			"database::database":              "database",
			"database::password":              os.Getenv("GOPATH"),
			"database::conns::maxconnection":  12,
			"database::conns::autoconnect":    true,
			"database::conns::connectioninfo": "info",
			"database::conns::root":           os.Getenv("GOPATH"),
			"unknown":                         "",
		}
	)

	f, err := os.Create("testjson.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testjson.conf")
	jsonconf, err := NewConfig("json", "testjson.conf")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = jsonconf.Int(k)
		case int64:
			value, err = jsonconf.Int64(k)
		case float64:
			value, err = jsonconf.Float(k)
		case bool:
			value, err = jsonconf.Bool(k)
		case []string:
			value = jsonconf.Strings(k)
		case string:
			value = jsonconf.String(k)
		default:
			value, err = jsonconf.DIY(k)
		}
		if err != nil {
			t.Fatalf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Fatalf("get key %q value, want %v got %v .", k, v, value)
		}

	}
	if err = jsonconf.Set("name", "bhojpur"); err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("name") != "bhojpur" {
		t.Fatal("get name error")
	}

	if db, err := jsonconf.DIY("database"); err != nil {
		t.Fatal(err)
	} else if m, ok := db.(map[string]interface{}); !ok {
		t.Log(db)
		t.Fatal("db not map[string]interface{}")
	} else {
		if m["host"].(string) != "host" {
			t.Fatal("get host err")
		}
	}

	if _, err := jsonconf.Int("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int")
	}

	if _, err := jsonconf.Int64("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int64")
	}

	if _, err := jsonconf.Float("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Float")
	}

	if _, err := jsonconf.DIY("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an interface{}")
	}

	if val := jsonconf.String("unknown"); val != "" {
		t.Error("unknown keys should return an empty string when expecting a String")
	}

	if _, err := jsonconf.Bool("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Bool")
	}

	if !jsonconf.DefaultBool("unknown", true) {
		t.Error("unknown keys with default value wrong")
	}
}
