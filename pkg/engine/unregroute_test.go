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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The unregroute_test.go contains tests for the unregister route
// functionality, that allows overriding route paths in children project
// that embed parent routers.

const contentRootOriginal = "ok-original-root"
const contentLevel1Original = "ok-original-level1"
const contentLevel2Original = "ok-original-level2"

const contentRootReplacement = "ok-replacement-root"
const contentLevel1Replacement = "ok-replacement-level1"
const contentLevel2Replacement = "ok-replacement-level2"

// TestPreUnregController will supply content for the original routes,
// before unregistration
type TestPreUnregController struct {
	Controller
}

func (tc *TestPreUnregController) GetFixedRoot() {
	tc.Ctx.Output.Body([]byte(contentRootOriginal))
}
func (tc *TestPreUnregController) GetFixedLevel1() {
	tc.Ctx.Output.Body([]byte(contentLevel1Original))
}
func (tc *TestPreUnregController) GetFixedLevel2() {
	tc.Ctx.Output.Body([]byte(contentLevel2Original))
}

// TestPostUnregController will supply content for the overriding routes,
// after the original ones are unregistered.
type TestPostUnregController struct {
	Controller
}

func (tc *TestPostUnregController) GetFixedRoot() {
	tc.Ctx.Output.Body([]byte(contentRootReplacement))
}
func (tc *TestPostUnregController) GetFixedLevel1() {
	tc.Ctx.Output.Body([]byte(contentLevel1Replacement))
}
func (tc *TestPostUnregController) GetFixedLevel2() {
	tc.Ctx.Output.Body([]byte(contentLevel2Replacement))
}

// TestUnregisterFixedRouteRoot replaces just the root fixed route path.
// In this case, for a path like "/level1/level2" or "/level1", those actions
// should remain intact, and continue to serve the original content.
func TestUnregisterFixedRouteRoot(t *testing.T) {

	var method = "GET"

	handler := NewControllerRegister()
	handler.Add("/", &TestPreUnregController{}, "get:GetFixedRoot")
	handler.Add("/level1", &TestPreUnregController{}, "get:GetFixedLevel1")
	handler.Add("/level1/level2", &TestPreUnregController{}, "get:GetFixedLevel2")

	// Test original root
	testHelperFnContentCheck(t, handler, "Test original root",
		method, "/", contentRootOriginal)

	// Test original level 1
	testHelperFnContentCheck(t, handler, "Test original level 1",
		method, "/level1", contentLevel1Original)

	// Test original level 2
	testHelperFnContentCheck(t, handler, "Test original level 2",
		method, "/level1/level2", contentLevel2Original)

	// Remove only the root path
	findAndRemoveSingleTree(handler.routers[method])

	// Replace the root path TestPreUnregController action with the action from
	// TestPostUnregController
	handler.Add("/", &TestPostUnregController{}, "get:GetFixedRoot")

	// Test replacement root (expect change)
	testHelperFnContentCheck(t, handler, "Test replacement root (expect change)", method, "/", contentRootReplacement)

	// Test level 1 (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test level 1 (expect no change from the original)", method, "/level1", contentLevel1Original)

	// Test level 2 (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test level 2 (expect no change from the original)", method, "/level1/level2", contentLevel2Original)

}

// TestUnregisterFixedRouteLevel1 replaces just the "/level1" fixed route path.
// In this case, for a path like "/level1/level2" or "/", those actions
// should remain intact, and continue to serve the original content.
func TestUnregisterFixedRouteLevel1(t *testing.T) {

	var method = "GET"

	handler := NewControllerRegister()
	handler.Add("/", &TestPreUnregController{}, "get:GetFixedRoot")
	handler.Add("/level1", &TestPreUnregController{}, "get:GetFixedLevel1")
	handler.Add("/level1/level2", &TestPreUnregController{}, "get:GetFixedLevel2")

	// Test original root
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original root",
		method, "/", contentRootOriginal)

	// Test original level 1
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original level 1",
		method, "/level1", contentLevel1Original)

	// Test original level 2
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original level 2",
		method, "/level1/level2", contentLevel2Original)

	// Remove only the level1 path
	subPaths := splitPath("/level1")
	if handler.routers[method].prefix == strings.Trim("/level1", "/ ") {
		findAndRemoveSingleTree(handler.routers[method])
	} else {
		findAndRemoveTree(subPaths, handler.routers[method], method)
	}

	// Replace the "level1" path TestPreUnregController action with the action from
	// TestPostUnregController
	handler.Add("/level1", &TestPostUnregController{}, "get:GetFixedLevel1")

	// Test replacement root (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test replacement root (expect no change from the original)", method, "/", contentRootOriginal)

	// Test level 1 (expect change)
	testHelperFnContentCheck(t, handler, "Test level 1 (expect change)", method, "/level1", contentLevel1Replacement)

	// Test level 2 (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test level 2 (expect no change from the original)", method, "/level1/level2", contentLevel2Original)

}

// TestUnregisterFixedRouteLevel2 unregisters just the "/level1/level2" fixed
// route path. In this case, for a path like "/level1" or "/", those actions
// should remain intact, and continue to serve the original content.
func TestUnregisterFixedRouteLevel2(t *testing.T) {

	var method = "GET"

	handler := NewControllerRegister()
	handler.Add("/", &TestPreUnregController{}, "get:GetFixedRoot")
	handler.Add("/level1", &TestPreUnregController{}, "get:GetFixedLevel1")
	handler.Add("/level1/level2", &TestPreUnregController{}, "get:GetFixedLevel2")

	// Test original root
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original root",
		method, "/", contentRootOriginal)

	// Test original level 1
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original level 1",
		method, "/level1", contentLevel1Original)

	// Test original level 2
	testHelperFnContentCheck(t, handler,
		"TestUnregisterFixedRouteLevel1.Test original level 2",
		method, "/level1/level2", contentLevel2Original)

	// Remove only the level2 path
	subPaths := splitPath("/level1/level2")
	if handler.routers[method].prefix == strings.Trim("/level1/level2", "/ ") {
		findAndRemoveSingleTree(handler.routers[method])
	} else {
		findAndRemoveTree(subPaths, handler.routers[method], method)
	}

	// Replace the "/level1/level2" path TestPreUnregController action with the action from
	// TestPostUnregController
	handler.Add("/level1/level2", &TestPostUnregController{}, "get:GetFixedLevel2")

	// Test replacement root (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test replacement root (expect no change from the original)", method, "/", contentRootOriginal)

	// Test level 1 (expect no change from the original)
	testHelperFnContentCheck(t, handler, "Test level 1 (expect no change from the original)", method, "/level1", contentLevel1Original)

	// Test level 2 (expect change)
	testHelperFnContentCheck(t, handler, "Test level 2 (expect change)", method, "/level1/level2", contentLevel2Replacement)

}

func testHelperFnContentCheck(t *testing.T, handler *ControllerRegister,
	testName, method, path, expectedBodyContent string) {

	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Errorf("httpRecorderBodyTest NewRequest error: %v", err)
		return
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	if body != expectedBodyContent {
		t.Errorf("%s: expected [%s], got [%s];", testName, expectedBodyContent, body)
	}
}
