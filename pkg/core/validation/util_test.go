package validation

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
	"log"
	"reflect"
	"testing"
)

type user struct {
	ID    int
	Tag   string `valid:"Maxx(aa)"`
	Name  string `valid:"Required;"`
	Age   int    `valid:"Required; Range(1, 140)"`
	match string `valid:"Required; Match(/^(test)?\\w*@(/test/);com$/);Max(2)"`
}

func TestGetValidFuncs(t *testing.T) {
	u := user{Name: "test", Age: 1}
	tf := reflect.TypeOf(u)
	var vfs []ValidFunc
	var err error

	f, _ := tf.FieldByName("ID")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 0 {
		t.Fatal("should get none ValidFunc")
	}

	f, _ = tf.FieldByName("Tag")
	if _, err = getValidFuncs(f); err.Error() != "doesn't exists Maxx valid function" {
		t.Fatal(err)
	}

	f, _ = tf.FieldByName("Name")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 1 {
		t.Fatal("should get 1 ValidFunc")
	}
	if vfs[0].Name != "Required" && len(vfs[0].Params) != 0 {
		t.Error("Required funcs should be got")
	}

	f, _ = tf.FieldByName("Age")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 2 {
		t.Fatal("should get 2 ValidFunc")
	}
	if vfs[0].Name != "Required" && len(vfs[0].Params) != 0 {
		t.Error("Required funcs should be got")
	}
	if vfs[1].Name != "Range" && len(vfs[1].Params) != 2 {
		t.Error("Range funcs should be got")
	}

	f, _ = tf.FieldByName("match")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 3 {
		t.Fatal("should get 3 ValidFunc but now is", len(vfs))
	}
}

type User struct {
	Name string `valid:"Required;MaxSize(5)" `
	Sex  string `valid:"Required;" label:"sex_label"`
	Age  int    `valid:"Required;Range(1, 140);" label:"age_label"`
}

func TestValidation(t *testing.T) {
	u := User{"man1238888456", "", 1140}
	valid := Validation{}
	b, err := valid.Valid(&u)
	if err != nil {
		// handle error
	}
	if !b {
		// validation does not pass
		// blabla...
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		if len(valid.Errors) != 3 {
			t.Error("must be has 3 error")
		}
	} else {
		t.Error("must be has 3 error")
	}
}

func TestCall(t *testing.T) {
	u := user{Name: "test", Age: 180}
	tf := reflect.TypeOf(u)
	var vfs []ValidFunc
	var err error
	f, _ := tf.FieldByName("Age")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	valid := &Validation{}
	vfs[1].Params = append([]interface{}{valid, u.Age}, vfs[1].Params...)
	if _, err = funcs.Call(vfs[1].Name, vfs[1].Params...); err != nil {
		t.Fatal(err)
	}
	if len(valid.Errors) != 1 {
		t.Error("age out of range should be has an error")
	}
}
