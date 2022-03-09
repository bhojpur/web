package template_test

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
	"testing"

	"github.com/bhojpur/web/pkg/template"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	tpl *template.Template
}

var (
	_          = Suite(&TestSuite{})
	testSuite2 = template.NewSet("test suite 2", template.MustNewLocalFileSystemLoader(""))
)

func parseTemplate(s string, c template.Context) string {
	t, err := testSuite2.FromString(s)
	if err != nil {
		panic(err)
	}
	out, err := t.Execute(c)
	if err != nil {
		panic(err)
	}
	return out
}

func parseTemplateFn(s string, c template.Context) func() {
	return func() {
		parseTemplate(s, c)
	}
}

func (s *TestSuite) TestMisc(c *C) {
	// Must
	// TODO: Add better error message
	c.Check(
		func() { template.Must(testSuite2.FromFile("template_tests/inheritance/base2.tpl")) },
		PanicMatches,
		`\[Error \(where: fromfile\) in .*template_tests[/\\]inheritance[/\\]doesnotexist.tpl | Line 1 Col 12 near 'doesnotexist.tpl'\] open .*template_tests[/\\]inheritance[/\\]doesnotexist.tpl: no such file or directory`,
	)

	// Context
	c.Check(parseTemplateFn("", template.Context{"'illegal": nil}), PanicMatches, ".*not a valid identifier.*")

	// Registers
	c.Check(template.RegisterFilter("escape", nil).Error(), Matches, ".*is already registered")
	c.Check(template.RegisterTag("for", nil).Error(), Matches, ".*is already registered")

	// ApplyFilter
	v, err := template.ApplyFilter("title", template.AsValue("this is a title"), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Check(v.String(), Equals, "This Is A Title")
	c.Check(func() {
		_, err := template.ApplyFilter("doesnotexist", nil, nil)
		if err != nil {
			panic(err)
		}
	}, PanicMatches, `\[Error \(where: applyfilter\)\] filter with name 'doesnotexist' not found`)
}

func (s *TestSuite) TestImplicitExecCtx(c *C) {
	tpl, err := template.FromString("{{ ImplicitExec }}")
	if err != nil {
		c.Fatalf("Error in FromString: %v", err)
	}

	val := "a stringy thing"

	res, err := tpl.Execute(template.Context{
		"Value": val,
		"ImplicitExec": func(ctx *template.ExecutionContext) string {
			return ctx.Public["Value"].(string)
		},
	})
	if err != nil {
		c.Fatalf("Error executing template: %v", err)
	}

	c.Check(res, Equals, val)

	// The implicit ctx should not be persisted from call-to-call
	res, err = tpl.Execute(template.Context{
		"ImplicitExec": func() string {
			return val
		},
	})

	if err != nil {
		c.Fatalf("Error executing template: %v", err)
	}

	c.Check(res, Equals, val)
}
