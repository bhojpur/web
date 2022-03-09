package dev

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
	"os"

	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

var preCommit = `
goimports -w -format-only ./ \
ineffassign . \
staticcheck -show-ignored -checks "-ST1017,-U1000,-ST1005,-S1034,-S1012,-SA4006,-SA6005,-SA1019,-SA1024" ./ \
`

// for now, we simply override pre-commit file
func initGitHook() {
	// pcf => pre-commit file
	pcfPath := "./.git/hooks/pre-commit"
	pcf, err := os.OpenFile(pcfPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		cliLogger.Log.Errorf("try to create or open file failed: %s, cause: %s", pcfPath, err.Error())
		return
	}

	defer pcf.Close()
	_, err = pcf.Write(([]byte)(preCommit))

	if err != nil {
		cliLogger.Log.Errorf("could not init githooks: %s", err.Error())
	} else {
		cliLogger.Log.Successf("The githooks has been added, the content is:\n %s ", preCommit)
	}
}
