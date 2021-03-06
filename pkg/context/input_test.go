package context

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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBind(t *testing.T) {
	type testItem struct {
		field string
		empty interface{}
		want  interface{}
	}
	type Human struct {
		ID   int
		Nick string
		Pwd  string
		Ms   bool
	}

	cases := []struct {
		request string
		valueGp []testItem
	}{
		{"/?p=str", []testItem{{"p", interface{}(""), interface{}("str")}}},

		{"/?p=", []testItem{{"p", "", ""}}},
		{"/?p=str", []testItem{{"p", "", "str"}}},

		{"/?p=123", []testItem{{"p", 0, 123}}},
		{"/?p=123", []testItem{{"p", uint(0), uint(123)}}},

		{"/?p=1.0", []testItem{{"p", 0.0, 1.0}}},
		{"/?p=1", []testItem{{"p", false, true}}},

		{"/?p=true", []testItem{{"p", false, true}}},
		{"/?p=ON", []testItem{{"p", false, true}}},
		{"/?p=on", []testItem{{"p", false, true}}},
		{"/?p=1", []testItem{{"p", false, true}}},
		{"/?p=2", []testItem{{"p", false, false}}},
		{"/?p=false", []testItem{{"p", false, false}}},

		{"/?p[a]=1&p[b]=2&p[c]=3", []testItem{{"p", map[string]int{}, map[string]int{"a": 1, "b": 2, "c": 3}}}},
		{"/?p[a]=v1&p[b]=v2&p[c]=v3", []testItem{{"p", map[string]string{}, map[string]string{"a": "v1", "b": "v2", "c": "v3"}}}},

		{"/?p[]=8&p[]=9&p[]=10", []testItem{{"p", []int{}, []int{8, 9, 10}}}},
		{"/?p[0]=8&p[1]=9&p[2]=10", []testItem{{"p", []int{}, []int{8, 9, 10}}}},
		{"/?p[0]=8&p[1]=9&p[2]=10&p[5]=14", []testItem{{"p", []int{}, []int{8, 9, 10, 0, 0, 14}}}},
		{"/?p[0]=8.0&p[1]=9.0&p[2]=10.0", []testItem{{"p", []float64{}, []float64{8.0, 9.0, 10.0}}}},

		{"/?p[]=10&p[]=9&p[]=8", []testItem{{"p", []string{}, []string{"10", "9", "8"}}}},
		{"/?p[0]=8&p[1]=9&p[2]=10", []testItem{{"p", []string{}, []string{"8", "9", "10"}}}},

		{"/?p[0]=true&p[1]=false&p[2]=true&p[5]=1&p[6]=ON&p[7]=other", []testItem{{"p", []bool{}, []bool{true, false, true, false, false, true, true, false}}}},

		{"/?human.Nick=bhojpur", []testItem{{"human", Human{}, Human{Nick: "bhojpur"}}}},
		{"/?human.ID=888&human.Nick=bhojpur&human.Ms=true&human[Pwd]=pass", []testItem{{"human", Human{}, Human{ID: 888, Nick: "bhojpur", Ms: true, Pwd: "pass"}}}},
		{"/?human[0].ID=888&human[0].Nick=bhojpur&human[0].Ms=true&human[0][Pwd]=pass01&human[1].ID=999&human[1].Nick=ysqi&human[1].Ms=On&human[1].Pwd=pass02",
			[]testItem{{"human", []Human{}, []Human{
				{ID: 888, Nick: "bhojpur", Ms: true, Pwd: "pass01"},
				{ID: 999, Nick: "ysqi", Ms: true, Pwd: "pass02"},
			}}}},

		{
			"/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&human.Nick=bhojpur",
			[]testItem{
				{"id", 0, 123},
				{"isok", false, true},
				{"ft", 0.0, 1.2},
				{"ol", []int{}, []int{1, 2}},
				{"ul", []string{}, []string{"str", "array"}},
				{"human", Human{}, Human{Nick: "bhojpur"}},
			},
		},
	}
	for _, c := range cases {
		r, _ := http.NewRequest("GET", c.request, nil)
		bhojpurInput := NewInput()
		bhojpurInput.Context = NewContext()
		bhojpurInput.Context.Reset(httptest.NewRecorder(), r)

		for _, item := range c.valueGp {
			got := item.empty
			err := bhojpurInput.Bind(&got, item.field)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, item.want) {
				t.Fatalf("Bind %q error,should be:\n%#v \ngot:\n%#v", item.field, item.want, got)
			}
		}

	}
}

func TestSubDomain(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.bhojpur.net/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=bhojpur", nil)
	bhojpurInput := NewInput()
	bhojpurInput.Context = NewContext()
	bhojpurInput.Context.Reset(httptest.NewRecorder(), r)

	subdomain := bhojpurInput.SubDomains()
	if subdomain != "www" {
		t.Fatal("Subdomain parse error, got" + subdomain)
	}

	r, _ = http.NewRequest("GET", "http://localhost/", nil)
	bhojpurInput.Context.Request = r
	if bhojpurInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, should be empty, got " + bhojpurInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.bhojpur.net/", nil)
	bhojpurInput.Context.Request = r
	if bhojpurInput.SubDomains() != "aa.bb" {
		t.Fatal("Subdomain parse error, got " + bhojpurInput.SubDomains())
	}

	/* TODO Fix this
	r, _ = http.NewRequest("GET", "http://127.0.0.1/", nil)
	bhojpurInput.Context.Request = r
	if bhojpurInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + bhojpurInput.SubDomains())
	}
	*/

	r, _ = http.NewRequest("GET", "http://bhojpur.net/", nil)
	bhojpurInput.Context.Request = r
	if bhojpurInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + bhojpurInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.cc.dd.bhojpur.net/", nil)
	bhojpurInput.Context.Request = r
	if bhojpurInput.SubDomains() != "aa.bb.cc.dd" {
		t.Fatal("Subdomain parse error, got " + bhojpurInput.SubDomains())
	}
}

func TestParams(t *testing.T) {
	inp := NewInput()

	inp.SetParam("p1", "val1_ver1")
	inp.SetParam("p2", "val2_ver1")
	inp.SetParam("p3", "val3_ver1")
	if l := inp.ParamsLen(); l != 3 {
		t.Fatalf("Input.ParamsLen wrong value: %d, expected %d", l, 3)
	}

	if val := inp.Param("p1"); val != "val1_ver1" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver1")
	}
	if val := inp.Param("p3"); val != "val3_ver1" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val3_ver1")
	}
	vals := inp.Params()
	expected := map[string]string{
		"p1": "val1_ver1",
		"p2": "val2_ver1",
		"p3": "val3_ver1",
	}
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("Input.Params wrong value: %s, expected %s", vals, expected)
	}

	// overwriting existing params
	inp.SetParam("p1", "val1_ver2")
	inp.SetParam("p2", "val2_ver2")
	expected = map[string]string{
		"p1": "val1_ver2",
		"p2": "val2_ver2",
		"p3": "val3_ver1",
	}
	vals = inp.Params()
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("Input.Params wrong value: %s, expected %s", vals, expected)
	}

	if l := inp.ParamsLen(); l != 3 {
		t.Fatalf("Input.ParamsLen wrong value: %d, expected %d", l, 3)
	}

	if val := inp.Param("p1"); val != "val1_ver2" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver2")
	}

	if val := inp.Param("p2"); val != "val2_ver2" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver2")
	}
}

func BenchmarkQuery(b *testing.B) {
	bhojpurInput := NewInput()
	bhojpurInput.Context = NewContext()
	bhojpurInput.Context.Request, _ = http.NewRequest("POST", "http://www.bhojpur.net/?q=foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bhojpurInput.Query("q")
		}
	})
}
