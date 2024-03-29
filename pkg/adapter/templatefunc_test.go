package adapter

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
	"html/template"
	"net/url"
	"testing"
	"time"
)

func TestSubstr(t *testing.T) {
	s := `012345`
	if Substr(s, 0, 2) != "01" {
		t.Error("should be equal")
	}
	if Substr(s, 0, 100) != "012345" {
		t.Error("should be equal")
	}
	if Substr(s, 12, 100) != "012345" {
		t.Error("should be equal")
	}
}

func TestHtml2str(t *testing.T) {
	h := `<HTML><style></style><script>x<x</script></HTML><123>  123\n


	\n`
	if HTML2str(h) != "123\\n\n\\n" {
		t.Error("should be equal")
	}
}

func TestDateFormat(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	if ss := DateFormat(tt, "2006-01-02 15:04:05"); ss != "2013-07-01 13:27:42" {
		t.Errorf("2013-07-01 13:27:42 does not equal %v", ss)
	}
}

func TestDate(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	if ss := Date(tt, "Y-m-d H:i:s"); ss != "2013-07-01 13:27:42" {
		t.Errorf("2013-07-01 13:27:42 does not equal %v", ss)
	}
	if ss := Date(tt, "y-n-j h:i:s A"); ss != "13-7-1 01:27:42 PM" {
		t.Errorf("13-7-1 01:27:42 PM does not equal %v", ss)
	}
	if ss := Date(tt, "D, d M Y g:i:s a"); ss != "Mon, 01 Jul 2013 1:27:42 pm" {
		t.Errorf("Mon, 01 Jul 2013 1:27:42 pm does not equal %v", ss)
	}
	if ss := Date(tt, "l, d F Y G:i:s"); ss != "Monday, 01 July 2013 13:27:42" {
		t.Errorf("Monday, 01 July 2013 13:27:42 does not equal %v", ss)
	}
}

func TestCompareRelated(t *testing.T) {
	if !Compare("abc", "abc") {
		t.Error("should be equal")
	}
	if Compare("abc", "aBc") {
		t.Error("should be not equal")
	}
	if !Compare("1", 1) {
		t.Error("should be equal")
	}
	if CompareNot("abc", "abc") {
		t.Error("should be equal")
	}
	if !CompareNot("abc", "aBc") {
		t.Error("should be not equal")
	}
	if !NotNil("a string") {
		t.Error("should not be nil")
	}
}

func TestHtmlquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	if Htmlquote(s) != h {
		t.Error("should be equal")
	}
}

func TestHtmlunquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	if Htmlunquote(h) != s {
		t.Error("should be equal")
	}
}

func TestParseForm(t *testing.T) {
	type ExtendInfo struct {
		Hobby []string `form:"hobby"`
		Memo  string
	}

	type OtherInfo struct {
		Organization string `form:"organization"`
		Title        string `form:"title"`
		ExtendInfo
	}

	type user struct {
		ID      int         `form:"-"`
		tag     string      `form:"tag"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text"`
		Email   string
		Intro   string    `form:",textarea"`
		StrBool bool      `form:"strbool"`
		Date    time.Time `form:"date,2006-01-02"`
		OtherInfo
	}

	u := user{}
	form := url.Values{
		"ID":           []string{"1"},
		"-":            []string{"1"},
		"tag":          []string{"no"},
		"username":     []string{"test"},
		"age":          []string{"40"},
		"Email":        []string{"test@bhojpur.net"},
		"Intro":        []string{"I am an Engineer!"},
		"strbool":      []string{"yes"},
		"date":         []string{"2014-11-12"},
		"organization": []string{"bhojpur"},
		"title":        []string{"CXO"},
		"hobby":        []string{"", "Basketball", "Football"},
		"memo":         []string{"nothing"},
	}
	if err := ParseForm(form, u); err == nil {
		t.Fatal("nothing will be changed")
	}
	if err := ParseForm(form, &u); err != nil {
		t.Fatal(err)
	}
	if u.ID != 0 {
		t.Errorf("ID should equal 0 but got %v", u.ID)
	}
	if len(u.tag) != 0 {
		t.Errorf("tag's length should equal 0 but got %v", len(u.tag))
	}
	if u.Name.(string) != "test" {
		t.Errorf("Name should equal `test` but got `%v`", u.Name.(string))
	}
	if u.Age != 40 {
		t.Errorf("Age should equal 40 but got %v", u.Age)
	}
	if u.Email != "test@bhojpur.net" {
		t.Errorf("Email should equal `test@bhojpur.net` but got `%v`", u.Email)
	}
	if u.Intro != "I am an Engineer!" {
		t.Errorf("Intro should equal `I am an Engineer!` but got `%v`", u.Intro)
	}
	if !u.StrBool {
		t.Errorf("strboll should equal `true`, but got `%v`", u.StrBool)
	}
	y, m, d := u.Date.Date()
	if y != 2014 || m.String() != "November" || d != 12 {
		t.Errorf("Date should equal `2014-11-12`, but got `%v`", u.Date.String())
	}
	if u.Organization != "bhojpur" {
		t.Errorf("Organization should equal `bhojpur`, but got `%v`", u.Organization)
	}
	if u.Title != "CXO" {
		t.Errorf("Title should equal `CXO`, but got `%v`", u.Title)
	}
	if u.Hobby[0] != "" {
		t.Errorf("Hobby should equal ``, but got `%v`", u.Hobby[0])
	}
	if u.Hobby[1] != "Basketball" {
		t.Errorf("Hobby should equal `Basketball`, but got `%v`", u.Hobby[1])
	}
	if u.Hobby[2] != "Football" {
		t.Errorf("Hobby should equal `Football`, but got `%v`", u.Hobby[2])
	}
	if len(u.Memo) != 0 {
		t.Errorf("Memo's length should equal 0 but got %v", len(u.Memo))
	}
}

func TestRenderForm(t *testing.T) {
	type user struct {
		ID      int         `form:"-"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text,年龄："`
		Sex     string
		Email   []string
		Intro   string `form:",textarea"`
		Ignored string `form:"-"`
	}

	u := user{Name: "test", Intro: "Some Text"}
	output := RenderForm(u)
	if output != template.HTML("") {
		t.Errorf("output should be empty but got %v", output)
	}
	output = RenderForm(&u)
	result := template.HTML(
		`Name: <input name="username" type="text" value="test"></br>` +
			`年龄：<input name="age" type="text" value="0"></br>` +
			`Sex: <input name="Sex" type="text" value=""></br>` +
			`Intro: <textarea name="Intro">Some Text</textarea>`)
	if output != result {
		t.Errorf("output should equal `%v` but got `%v`", result, output)
	}
}

func TestMapGet(t *testing.T) {
	// test one level map
	m1 := map[string]int64{
		"a": 1,
		"1": 2,
	}

	if res, err := MapGet(m1, "a"); err == nil {
		if res.(int64) != 1 {
			t.Errorf("Should return 1, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, "1"); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, 1); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 2 level map
	m2 := M{
		"1": map[string]float64{
			"2": 3.5,
		},
	}

	if res, err := MapGet(m2, 1, 2); err == nil {
		if res.(float64) != 3.5 {
			t.Errorf("Should return 3.5, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 5 level map
	m5 := M{
		"1": M{
			"2": M{
				"3": M{
					"4": M{
						"5": 1.2,
					},
				},
			},
		},
	}

	if res, err := MapGet(m5, 1, 2, 3, 4, 5); err == nil {
		if res.(float64) != 1.2 {
			t.Errorf("Should return 1.2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// check whether element not exists in map
	if res, err := MapGet(m5, 5, 4, 3, 2, 1); err == nil {
		if res != nil {
			t.Errorf("Should return nil, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}
}
