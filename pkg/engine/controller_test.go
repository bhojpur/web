package engine

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
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	ctxsvr "github.com/bhojpur/web/pkg/context"
)

func TestGetInt(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", "40")
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt("age")
	if val != 40 {
		t.Errorf("TestGetInt expect 40,get %T,%v", val, val)
	}
}

func TestGetInt8(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", "40")
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt8("age")
	if val != 40 {
		t.Errorf("TestGetInt8 expect 40,get %T,%v", val, val)
	}
	//Output: int8
}

func TestGetInt16(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", "40")
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt16("age")
	if val != 40 {
		t.Errorf("TestGetInt16 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt32(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", "40")
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt32("age")
	if val != 40 {
		t.Errorf("TestGetInt32 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt64(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", "40")
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt64("age")
	if val != 40 {
		t.Errorf("TestGeetInt64 expect 40,get %T,%v", val, val)
	}
}

func TestGetUint8(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", strconv.FormatUint(math.MaxUint8, 10))
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetUint8("age")
	if val != math.MaxUint8 {
		t.Errorf("TestGetUint8 expect %v,get %T,%v", math.MaxUint8, val, val)
	}
}

func TestGetUint16(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", strconv.FormatUint(math.MaxUint16, 10))
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetUint16("age")
	if val != math.MaxUint16 {
		t.Errorf("TestGetUint16 expect %v,get %T,%v", math.MaxUint16, val, val)
	}
}

func TestGetUint32(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", strconv.FormatUint(math.MaxUint32, 10))
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetUint32("age")
	if val != math.MaxUint32 {
		t.Errorf("TestGetUint32 expect %v,get %T,%v", math.MaxUint32, val, val)
	}
}

func TestGetUint64(t *testing.T) {
	i := ctxsvr.NewInput()
	i.SetParam("age", strconv.FormatUint(math.MaxUint64, 10))
	ctx := &ctxsvr.Context{Input: i}
	ctrlr := ctxsvr.Controller{Ctx: ctx}
	val, _ := ctrlr.GetUint64("age")
	if val != math.MaxUint64 {
		t.Errorf("TestGetUint64 expect %v,get %T,%v", uint64(math.MaxUint64), val, val)
	}
}

func TestAdditionalViewPaths(t *testing.T) {
	wkdir, err := os.Getwd()
	assert.Nil(t, err)
	dir1 := filepath.Join(wkdir, "_beeTmp", "TestAdditionalViewPaths")
	dir2 := filepath.Join(wkdir, "_beeTmp2", "TestAdditionalViewPaths")
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)

	dir1file := "file1.tpl"
	dir2file := "file2.tpl"

	genFile := func(dir string, name string, content string) {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0777)
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			defer f.Close()
			f.WriteString(content)
			f.Close()
		}

	}
	genFile(dir1, dir1file, `<div>{{.Content}}</div>`)
	genFile(dir2, dir2file, `<html>{{.Content}}</html>`)

	ctxsvr.AddViewPath(dir1)
	ctxsvr.AddViewPath(dir2)

	ctrl := ctxsvr.Controller{
		TplName:  "file1.tpl",
		ViewPath: dir1,
	}
	ctrl.Data = map[interface{}]interface{}{
		"Content": "value2",
	}
	if result, err := ctrl.RenderString(); err != nil {
		t.Fatal(err)
	} else {
		if result != "<div>value2</div>" {
			t.Fatalf("TestAdditionalViewPaths expect %s got %s", "<div>value2</div>", result)
		}
	}

	func() {
		ctrl.TplName = "file2.tpl"
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("TestAdditionalViewPaths expected error")
			}
		}()
		ctrl.RenderString()
	}()

	ctrl.TplName = "file2.tpl"
	ctrl.ViewPath = dir2
	ctrl.RenderString()
}
