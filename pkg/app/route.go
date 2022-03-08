package app

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
	"reflect"
	"regexp"
	"sync"
)

var (
	routes = makeRouter()
)

// Route associates the type of the given component to the given path.
//
// When a page is requested and matches the route, a new instance of the given
// component is created before being displayed.
func Route(path string, c Composer) {
	routes.route(path, c)
}

// RouteWithRegexp associates the type of the given component to the given
// regular expression pattern.
//
// Patterns use the Go standard regexp format.
//
// When a page is requested and matches the pattern, a new instance of the given
// component is created before being displayed.
func RouteWithRegexp(pattern string, c Composer) {
	routes.routeWithRegexp(pattern, c)
}

type router struct {
	mu               sync.RWMutex
	routes           map[string]reflect.Type
	routesWithRegexp []regexpRoute
}

func makeRouter() router {
	return router{
		routes: make(map[string]reflect.Type),
	}
}

func (r *router) route(path string, c Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes[path] = reflect.TypeOf(c)
}

func (r *router) routeWithRegexp(pattern string, c Composer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routesWithRegexp = append(r.routesWithRegexp, regexpRoute{
		regexp:    regexp.MustCompile(pattern),
		compoType: reflect.TypeOf(c),
	})
}

func (r *router) createComponent(path string) (Composer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	compoType, isRouted := r.routes[path]
	if !isRouted {
		for _, rwr := range r.routesWithRegexp {
			if rwr.regexp.MatchString(path) {
				compoType = rwr.compoType
				isRouted = true
				break
			}
		}
	}
	if !isRouted {
		return nil, false
	}

	compo := reflect.New(compoType.Elem()).Interface().(Composer)
	return compo, true
}

func (r *router) len() int {
	return len(r.routes) + len(r.routesWithRegexp)
}

type regexpRoute struct {
	regexp    *regexp.Regexp
	compoType reflect.Type
}
