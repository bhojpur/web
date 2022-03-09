package swaggergen

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
	"go/ast"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//package model
//
//import (
//"sync"
//
//"bhojpur.net/pkgnotexist"
//"github.com/shopspring/decimal"
//)
//
//type Object struct {
//	Field1 decimal.Decimal
//	Field2 pkgnotexist.TestType
//	Field3 sync.Map
//}
func TestCheckAndLoadPackageOnGoMod(t *testing.T) {
	defer os.Setenv("GO111MODULE", os.Getenv("GO111MODULE"))
	os.Setenv("GO111MODULE", "on")

	testCases := []struct {
		pkgName       string
		pkgImportPath string
		imports       []*ast.ImportSpec
		realType      string
		curPkgName    string
		expected      bool
	}{
		{
			pkgName:       "decimal",
			pkgImportPath: "github.com/shopspring/decimal",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "github.com/shopspring/decimal",
					},
				},
			},
			realType:   "decimal.Decimal",
			curPkgName: "model",
			expected:   true,
		},
		{
			pkgName:       "pkgnotexist",
			pkgImportPath: "bhojpur.net/pkgnotexist",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "bhojpur.net/pkgnotexist",
					},
				},
			},
			realType:   "pkgnotexist.TestType",
			curPkgName: "model",
			expected:   false,
		},
		{
			pkgName:       "sync",
			pkgImportPath: "sync",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "sync",
					},
				},
			},
			realType:   "sync.Map",
			curPkgName: "model",
			expected:   false,
		},
	}

	for _, test := range testCases {
		checkAndLoadPackage(test.imports, test.realType, test.curPkgName)
		result := false
		for _, v := range astPkgs {
			if v.Name == test.pkgName {
				result = true
				break
			}
		}
		if test.expected != result {
			t.Fatalf("load module error, expected: %v, result: %v", test.expected, result)
		}
	}
}

//package model
//
//import (
//"sync"
//
//"bhojpur.net/comm"
//"bhojpur.net/pkgnotexist"
//)
//
//type Object struct {
//	Field1 comm.Common
//	Field2 pkgnotexist.TestType
//	Field3 sync.Map
//}
func TestCheckAndLoadPackageOnGoPath(t *testing.T) {
	var (
		testCommPkg = `
package comm

type Common struct {
	Code  string
	Error string
}
`
	)

	gopath, err := ioutil.TempDir("", "gobuild-gopath")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(gopath)

	if err := os.MkdirAll(filepath.Join(gopath, "src/bhojpur.net/comm"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(gopath, "src/bhojpur.net/comm/comm.go"), []byte(testCommPkg), 0666); err != nil {
		t.Fatal(err)
	}

	defer os.Setenv("GO111MODULE", os.Getenv("GO111MODULE"))
	os.Setenv("GO111MODULE", "off")
	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", gopath)
	build.Default.GOPATH = gopath

	testCases := []struct {
		pkgName       string
		pkgImportPath string
		imports       []*ast.ImportSpec
		realType      string
		curPkgName    string
		expected      bool
	}{
		{
			pkgName:       "comm",
			pkgImportPath: "bhojpur.net/comm",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "bhojpur.net/comm",
					},
				},
			},
			realType:   "comm.Common",
			curPkgName: "model",
			expected:   true,
		},
		{
			pkgName:       "pkgnotexist",
			pkgImportPath: "bhojpur.net/pkgnotexist",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "bhojpur.net/pkgnotexist",
					},
				},
			},
			realType:   "pkgnotexist.TestType",
			curPkgName: "model",
			expected:   false,
		},
		{
			pkgName:       "sync",
			pkgImportPath: "sync",
			imports: []*ast.ImportSpec{
				{
					Path: &ast.BasicLit{
						Value: "sync",
					},
				},
			},
			realType:   "sync.Map",
			curPkgName: "model",
			expected:   false,
		},
	}

	for _, test := range testCases {
		checkAndLoadPackage(test.imports, test.realType, test.curPkgName)
		result := false
		for _, v := range astPkgs {
			if v.Name == test.pkgName {
				result = true
				break
			}
		}
		if test.expected != result {
			t.Fatalf("load module error, expected: %v, result: %v", test.expected, result)
		}
	}
}
