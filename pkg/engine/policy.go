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
	"strings"

	"github.com/bhojpur/web/pkg/context"
)

// PolicyFunc defines a policy function which is invoked before the controller handler is executed.
type PolicyFunc func(*context.Context)

// FindPolicy Find Router info for URL
func (p *ControllerRegister) FindPolicy(cont *context.Context) []PolicyFunc {
	var urlPath = cont.Input.URL()
	if !BasConfig.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}
	httpMethod := cont.Input.Method()
	isWildcard := false
	// Find policy for current method
	t, ok := p.policies[httpMethod]
	// If not found - find policy for whole controller
	if !ok {
		t, ok = p.policies["*"]
		isWildcard = true
	}
	if ok {
		runObjects := t.Match(urlPath, cont)
		if r, ok := runObjects.([]PolicyFunc); ok {
			return r
		} else if !isWildcard {
			// If no policies found and we checked not for "*" method - try to find it
			t, ok = p.policies["*"]
			if ok {
				runObjects = t.Match(urlPath, cont)
				if r, ok = runObjects.([]PolicyFunc); ok {
					return r
				}
			}
		}
	}
	return nil
}

func (p *ControllerRegister) addToPolicy(method, pattern string, r ...PolicyFunc) {
	method = strings.ToUpper(method)
	p.enablePolicy = true
	if !BasConfig.RouterCaseSensitive {
		pattern = strings.ToLower(pattern)
	}
	if t, ok := p.policies[method]; ok {
		t.AddRouter(pattern, r)
	} else {
		t := NewTree()
		t.AddRouter(pattern, r)
		p.policies[method] = t
	}
}

// Policy Register new policy in Bhojpur.NET Platform application
func Policy(pattern, method string, policy ...PolicyFunc) {
	BhojpurApp.Handlers.addToPolicy(method, pattern, policy...)
}

// Find policies and execute if were found
func (p *ControllerRegister) execPolicy(cont *context.Context, urlPath string) (started bool) {
	if !p.enablePolicy {
		return false
	}
	// Find Policy for method
	policyList := p.FindPolicy(cont)
	if len(policyList) > 0 {
		// Run policies
		for _, runPolicy := range policyList {
			runPolicy(cont)
			if cont.ResponseWriter.Started {
				return true
			}
		}
		return false
	}
	return false
}
