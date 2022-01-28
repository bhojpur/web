package yaml

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/core/config"
)

func TestYaml(t *testing.T) {

	var (
		yamlcontext = `
"appname": bhojpurapi
"httpport": 8080
"mysqlport": 3600
"PI": 3.1415976
"runmode": dev
"autorender": false
"copyrequestbody": true
"PATH": GOPATH
"path1": ${GOPATH}
"path2": ${GOPATH||/home/go}
"empty": "" 
"user":
  "name": "tom"
  "age": 13
`

		keyValue = map[string]interface{}{
			"appname":         "bhojpurapi",
			"httpport":        8080,
			"mysqlport":       int64(3600),
			"PI":              3.1415976,
			"runmode":         "dev",
			"autorender":      false,
			"copyrequestbody": true,
			"PATH":            "GOPATH",
			"path1":           os.Getenv("GOPATH"),
			"path2":           os.Getenv("GOPATH"),
			"error":           "",
			"emptystrings":    []string{},
		}
	)
	f, err := os.Create("testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(yamlcontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testyaml.conf")
	yamlconf, err := config.NewConfig("yaml", "testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}

	res, _ := yamlconf.String("appname")
	if res != "bhojpurapi" {
		t.Fatal("appname not equal to bhojpurapi")
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

		switch v.(type) {
		case int:
			value, err = yamlconf.Int(k)
		case int64:
			value, err = yamlconf.Int64(k)
		case float64:
			value, err = yamlconf.Float(k)
		case bool:
			value, err = yamlconf.Bool(k)
		case []string:
			value, err = yamlconf.Strings(k)
		case string:
			value, err = yamlconf.String(k)
		default:
			value, err = yamlconf.DIY(k)
		}
		if err != nil {
			t.Errorf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}

	}

	if err = yamlconf.Set("name", "bhojpur"); err != nil {
		t.Fatal(err)
	}
	res, _ = yamlconf.String("name")
	if res != "bhojpur" {
		t.Fatal("get name error")
	}

	sub, err := yamlconf.Sub("user")
	assert.Nil(t, err)
	assert.NotNil(t, sub)
	name, err := sub.String("name")
	assert.Nil(t, err)
	assert.Equal(t, "tom", name)

	age, err := sub.Int("age")
	assert.Nil(t, err)
	assert.Equal(t, 13, age)

	user := &User{}

	err = sub.Unmarshaler("", user)
	assert.Nil(t, err)
	assert.Equal(t, "tom", user.Name)
	assert.Equal(t, 13, user.Age)

	user = &User{}

	err = yamlconf.Unmarshaler("user", user)
	assert.Nil(t, err)
	assert.Equal(t, "tom", user.Name)
	assert.Equal(t, 13, user.Age)
}

type User struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}
