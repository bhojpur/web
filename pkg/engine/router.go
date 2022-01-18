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
	"io"
	"net/http"
	"strings"
)

type endpoint struct {
	handler    http.Handler
	parameters []string
}

type Router struct {
	methodEndpoints map[string]*endpoint
	staticBranches  map[string]*Router
	parameterBranch *Router
}

// ServeHTTP makes Router implement standard http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if endpoint, arguments, ok := r.findEndpoint(req.Method, segments, []string{}); ok {
		addRouteArgumentsToRequest(endpoint.parameters, arguments, req)
		endpoint.handler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) Get(path string, handler http.Handler) {
	r.AddRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler http.Handler) {
	r.AddRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler http.Handler) {
	r.AddRoute("PUT", path, handler)
}

func (r *Router) Delete(path string, handler http.Handler) {
	r.AddRoute("DELETE", path, handler)
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
	segments := segmentizePath(path)
	parameters := extractParameterNames(segments)
	endpoint := &endpoint{handler: handler, parameters: parameters}
	r.addRouteFromSegments(method, segments, endpoint)
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.methodEndpoints = make(map[string]*endpoint)
	r.staticBranches = make(map[string]*Router)
	return r
}

func (r *Router) addRouteFromSegments(method string, segments []string, endpoint *endpoint) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]

		var subrouter *Router
		if strings.HasPrefix(head, ":") {
			if r.parameterBranch == nil {
				r.parameterBranch = NewRouter()
			}
			subrouter = r.parameterBranch
		} else {
			if _, present := r.staticBranches[head]; !present {
				r.staticBranches[head] = NewRouter()
			}
			subrouter = r.staticBranches[head]
		}

		subrouter.addRouteFromSegments(method, tail, endpoint)
	} else {
		r.methodEndpoints[method] = endpoint
	}
}

func segmentizePath(path string) (segments []string) {
	for _, s := range strings.Split(path, "/") {
		if len(s) != 0 {
			segments = append(segments, s)
		}
	}
	return
}

func (r *Router) findEndpoint(method string, segments []string, pathArguments []string) (*endpoint, []string, bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticBranches[head]; present {
			return subrouter.findEndpoint(method, tail, pathArguments)
		} else if r.parameterBranch != nil {
			pathArguments = append(pathArguments, head)
			return r.parameterBranch.findEndpoint(method, tail, pathArguments)
		} else {
			return nil, nil, false
		}
	}
	endpoint, present := r.methodEndpoints[method]
	return endpoint, pathArguments, present
}

func addRouteArgumentsToRequest(names, values []string, req *http.Request) {
	query := req.URL.Query()
	for i := 0; i < len(names); i++ {
		query.Set(names[i], values[i])
	}
	req.URL.RawQuery = query.Encode()
}

func extractParameterNames(segments []string) (parameters []string) {
	for _, s := range segments {
		if strings.HasPrefix(s, ":") {
			parameters = append(parameters, s[1:])
		}
	}
	return
}
