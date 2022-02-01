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
	"github.com/bhojpur/web/pkg/adapter/context"
	bcontext "github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
)

// PolicyFunc defines a policy function which is invoked before the controller handler is executed.
type PolicyFunc func(*context.Context)

// FindPolicy Find Router info for URL
func (p *ControllerRegister) FindPolicy(cont *context.Context) []PolicyFunc {
	pf := (*web.ControllerRegister)(p).FindPolicy((*bcontext.Context)(cont))
	npf := newToOldPolicyFunc(pf)
	return npf
}

func newToOldPolicyFunc(pf []web.PolicyFunc) []PolicyFunc {
	npf := make([]PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *context.Context) {
			f((*bcontext.Context)(c))
		})
	}
	return npf
}

func oldToNewPolicyFunc(pf []PolicyFunc) []web.PolicyFunc {
	npf := make([]web.PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *bcontext.Context) {
			f((*context.Context)(c))
		})
	}
	return npf
}

// Policy Register new policy in Bhojpur
func Policy(pattern, method string, policy ...PolicyFunc) {
	pf := oldToNewPolicyFunc(policy)
	web.Policy(pattern, method, pf...)
}
