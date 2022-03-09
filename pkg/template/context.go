package template

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
	"regexp"

	"errors"
)

var reIdentifiers = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var autoescape = true

func SetAutoescape(newValue bool) {
	autoescape = newValue
}

// A Context type provides constants, variables, instances or functions to a template.
//
// It automatically provides meta-information or functions through the "bhojpur"-key.
// Currently, context["bhojpur"] contains the following keys:
//  1. version: returns the version string
//
// Template examples for accessing items from your context:
//     {{ myconstant }}
//     {{ myfunc("test", 42) }}
//     {{ user.name }}
//     {{ bhojpur.version }}
type Context map[string]interface{}

func (c Context) checkForValidIdentifiers() *Error {
	for k, v := range c {
		if !reIdentifiers.MatchString(k) {
			return &Error{
				Sender:    "checkForValidIdentifiers",
				OrigError: fmt.Errorf("context-key '%s' (value: '%+v') is not a valid identifier", k, v),
			}
		}
	}
	return nil
}

// Update updates this context with the key/value-pairs from another context.
func (c Context) Update(other Context) Context {
	for k, v := range other {
		c[k] = v
	}
	return c
}

// ExecutionContext contains all data important for the current rendering state.
//
// If you're writing a custom tag, your tag's Execute()-function will
// have access to the ExecutionContext. This struct stores anything
// about the current rendering process's Context including
// the Context provided by the user (field Public).
// You can safely use the Private context to provide data to the user's
// template (like a 'forloop'-information). The Shared-context is used
// to share data between tags. All ExecutionContexts share this context.
//
// Please be careful when accessing the Public data.
// PLEASE DO NOT MODIFY THE PUBLIC CONTEXT (read-only).
//
// To create your own execution context within tags, use the
// NewChildExecutionContext(parent) function.
type ExecutionContext struct {
	template *Template

	Autoescape bool
	Public     Context
	Private    Context
	Shared     Context
}

var templateMetaContext = Context{
	"version": Version,
}

func newExecutionContext(tpl *Template, ctx Context) *ExecutionContext {
	privateCtx := make(Context)

	// Make the Bhojpur Web template related funcs/vars available to the context
	privateCtx["template"] = templateMetaContext

	return &ExecutionContext{
		template: tpl,

		Public:     ctx,
		Private:    privateCtx,
		Autoescape: autoescape,
	}
}

func NewChildExecutionContext(parent *ExecutionContext) *ExecutionContext {
	newctx := &ExecutionContext{
		template: parent.template,

		Public:     parent.Public,
		Private:    make(Context),
		Autoescape: parent.Autoescape,
	}
	newctx.Shared = parent.Shared

	// Copy all existing private items
	newctx.Private.Update(parent.Private)

	return newctx
}

func (ctx *ExecutionContext) Error(msg string, token *Token) *Error {
	return ctx.OrigError(errors.New(msg), token)
}

func (ctx *ExecutionContext) OrigError(err error, token *Token) *Error {
	filename := ctx.template.name
	var line, col int
	if token != nil {
		// No tokens available
		// TODO: Add location (from where?)
		filename = token.Filename
		line = token.Line
		col = token.Col
	}
	return &Error{
		Template:  ctx.template,
		Filename:  filename,
		Line:      line,
		Column:    col,
		Token:     token,
		Sender:    "execution",
		OrigError: err,
	}
}

func (ctx *ExecutionContext) Logf(format string, args ...interface{}) {
	ctx.template.set.logf(format, args...)
}
