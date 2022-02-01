package berror

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
	"sync"
)

// A Code is an unsigned 32-bit error code as defined in the Bhojpur Web spec.
type Code interface {
	Code() uint32
	Module() string
	Desc() string
	Name() string
}

var defaultCodeRegistry = &codeRegistry{
	codes: make(map[uint32]*codeDefinition, 127),
}

// DefineCode defining a new Code
// Before defining a new code, please read Bhojpur Web specification.
// desc could be markdown doc
func DefineCode(code uint32, module string, name string, desc string) Code {
	res := &codeDefinition{
		code:   code,
		module: module,
		desc:   desc,
	}
	defaultCodeRegistry.lock.Lock()
	defer defaultCodeRegistry.lock.Unlock()

	if _, ok := defaultCodeRegistry.codes[code]; ok {
		panic(fmt.Sprintf("duplicate code, code %d has been registered", code))
	}
	defaultCodeRegistry.codes[code] = res
	return res
}

type codeRegistry struct {
	lock  sync.RWMutex
	codes map[uint32]*codeDefinition
}

func (cr *codeRegistry) Get(code uint32) (Code, bool) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()
	c, ok := cr.codes[code]
	return c, ok
}

type codeDefinition struct {
	code   uint32
	module string
	desc   string
	name   string
}

func (c *codeDefinition) Name() string {
	return c.name
}

func (c *codeDefinition) Code() uint32 {
	return c.code
}

func (c *codeDefinition) Module() string {
	return c.module
}

func (c *codeDefinition) Desc() string {
	return c.desc
}
