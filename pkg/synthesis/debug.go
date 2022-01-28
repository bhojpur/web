package synthesis

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
	"io"
)

// writeDebug writes the debug code file.
func writeDebug(w io.Writer, c *Config, toc []Asset) error {
	err := writeDebugHeader(w, c)
	if err != nil {
		return err
	}

	err = writeAssetFS(w, c)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeDebugAsset(w, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeDebugHeader writes output file headers.
// This targets debug builds.
func writeDebugHeader(w io.Writer, c *Config) error {
	var header string

	if c.HttpFileSystem {
		header = `import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	} else {
		header = `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	}

	_, err := fmt.Fprintf(w, `%s
)

// bhojpurRead reads the given file from disk. It returns an error on failure.
func bhojpurRead(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %%s at %%s: %%v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bhojpurFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bhojpurFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bhojpurFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bhojpurFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bhojpurFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bhojpurFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bhojpurFileInfo) Sys() interface{} {
	return nil
}

`, header)
	return err
}

// writeDebugAsset write a debug entry for the given asset.
// A debug entry is simply a function which reads the asset from
// the original file (e.g.: from disk).
func writeDebugAsset(w io.Writer, c *Config, asset *Asset) error {
	pathExpr := fmt.Sprintf("%q", asset.Path)
	if c.Dev {
		pathExpr = fmt.Sprintf("filepath.Join(rootDir, %q)", asset.Name)
	}

	_, err := fmt.Fprintf(w, `// %s reads file data from disk. It returns an error on failure.
func %s() (*asset, error) {
	path := %s
	name := %q
	bytes, err := bhojpurRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %%s at %%s: %%v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

`, asset.Func, asset.Func, pathExpr, asset.Name)
	return err
}
