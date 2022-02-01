package xml

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

	"github.com/bhojpur/web/pkg/adapter/config"
)

func TestXML(t *testing.T) {

	var (
		//xml parse should incluce in <config></config> tags
		xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>bhojpurapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
<path1>${GOPATH}</path1>
<path2>${GOPATH||/home/go}</path2>
<mysection>
<id>1</id>
<name>MySection</name>
</mysection>
</config>
`
		keyValue = map[string]interface{}{
			"appname":         "bhojpurapi",
			"httpport":        8080,
			"mysqlport":       int64(3600),
			"PI":              3.1415976,
			"runmode":         "dev",
			"autorender":      false,
			"copyrequestbody": true,
			"path1":           os.Getenv("GOPATH"),
			"path2":           os.Getenv("GOPATH"),
			"error":           "",
			"emptystrings":    []string{},
		}
	)

	f, err := os.Create("testxml.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(xmlcontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testxml.conf")

	xmlconf, err := config.NewConfig("xml", "testxml.conf")
	if err != nil {
		t.Fatal(err)
	}

	var xmlsection map[string]string
	xmlsection, err = xmlconf.GetSection("mysection")
	if err != nil {
		t.Fatal(err)
	}

	if len(xmlsection) == 0 {
		t.Error("section should not be empty")
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

		switch v.(type) {
		case int:
			value, err = xmlconf.Int(k)
		case int64:
			value, err = xmlconf.Int64(k)
		case float64:
			value, err = xmlconf.Float(k)
		case bool:
			value, err = xmlconf.Bool(k)
		case []string:
			value = xmlconf.Strings(k)
		case string:
			value = xmlconf.String(k)
		default:
			value, err = xmlconf.DIY(k)
		}
		if err != nil {
			t.Errorf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}

	}

	if err = xmlconf.Set("name", "bhojpur"); err != nil {
		t.Fatal(err)
	}
	if xmlconf.String("name") != "bhojpur" {
		t.Fatal("get name error")
	}
}
