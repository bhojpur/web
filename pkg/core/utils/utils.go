package utils

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
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && compareGoVersion(runtime.Version(), "go1.8") >= 0 {
		gopath = defaultGOPATH()
	}
	return filepath.SplitList(gopath)
}

func compareGoVersion(a, b string) int {
	reg := regexp.MustCompile("^\\d*")

	a = strings.TrimPrefix(a, "go")
	b = strings.TrimPrefix(b, "go")

	versionsA := strings.Split(a, ".")
	versionsB := strings.Split(b, ".")

	for i := 0; i < len(versionsA) && i < len(versionsB); i++ {
		versionA := versionsA[i]
		versionB := versionsB[i]

		vA, err := strconv.Atoi(versionA)
		if err != nil {
			str := reg.FindString(versionA)
			if str != "" {
				vA, _ = strconv.Atoi(str)
			} else {
				vA = -1
			}
		}

		vB, err := strconv.Atoi(versionB)
		if err != nil {
			str := reg.FindString(versionB)
			if str != "" {
				vB, _ = strconv.Atoi(str)
			} else {
				vB = -1
			}
		}

		if vA > vB {
			// vA = 12, vB = 8
			return 1
		} else if vA < vB {
			// vA = 6, vB = 8
			return -1
		} else if vA == -1 {
			// vA = rc1, vB = rc3
			return strings.Compare(versionA, versionB)
		}

		// vA = vB = 8
		continue
	}

	if len(versionsA) > len(versionsB) {
		return 1
	} else if len(versionsA) == len(versionsB) {
		return 0
	}

	return -1
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}
