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
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

// writeRelease writes the release code file.
func writeRelease(w io.Writer, c *Config, toc []Asset) error {
	err := writeReleaseHeader(w, c)
	if err != nil {
		return err
	}

	err = writeAssetFS(w, c)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeReleaseAsset(w, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeReleaseHeader writes output file headers.
// This targets release builds.
func writeReleaseHeader(w io.Writer, c *Config) error {
	var err error
	if c.NoCompress {
		if c.NoMemCopy {
			err = header_uncompressed_nomemcopy(w, c)
		} else {
			err = header_uncompressed_memcopy(w, c)
		}
	} else {
		if c.NoMemCopy {
			err = header_compressed_nomemcopy(w, c)
		} else {
			err = header_compressed_memcopy(w, c)
		}
	}
	if err != nil {
		return err
	}
	return header_release_common(w)
}

// writeReleaseAsset write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeReleaseAsset(w io.Writer, c *Config, asset *Asset) error {
	fd, err := os.Open(asset.Path)
	if err != nil {
		return err
	}

	defer fd.Close()

	if c.NoCompress {
		if c.NoMemCopy {
			err = uncompressed_nomemcopy(w, asset, fd)
		} else {
			err = uncompressed_memcopy(w, asset, fd)
		}
	} else {
		if c.NoMemCopy {
			err = compressed_nomemcopy(w, asset, fd)
		} else {
			err = compressed_memcopy(w, asset, fd)
		}
	}
	if err != nil {
		return err
	}
	return asset_release_common(w, c, asset)
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reasons
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

func header_compressed_nomemcopy(w io.Writer, c *Config) error {
	var header string

	if c.HttpFileSystem {
		header = `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	} else {
		header = `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	}

	_, err := fmt.Fprintf(w, `%s
)

func bhojpurRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %%q: %%v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

`, header)
	return err
}

func header_compressed_memcopy(w io.Writer, c *Config) error {
	var header string

	if c.HttpFileSystem {
		header = `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	} else {
		header = `import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"`
	}

	_, err := fmt.Fprintf(w, `%s
)

func bhojpurRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %%q: %%v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

`, header)
	return err
}

func header_uncompressed_nomemcopy(w io.Writer, c *Config) error {
	var header string

	if c.HttpFileSystem {
		header = `import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"`
	} else {
		header = `import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"`
	}

	_, err := fmt.Fprintf(w, `%s
)

func bhojpurRead(data, name string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b, nil
}

`, header)
	return err
}

func header_uncompressed_memcopy(w io.Writer, c *Config) error {
	var header string

	if c.HttpFileSystem {
		header = `import (
	"bytes"
	"fmt"
	"net/http"
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
`, header)
	return err
}

func header_release_common(w io.Writer) error {
	_, err := fmt.Fprintf(w, `type asset struct {
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

`)
	return err
}

func compressed_nomemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = "`, asset.Func)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	gz.Close()

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `"

func %sBytes() ([]byte, error) {
	return bhojpurRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func compressed_memcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = []byte("`, asset.Func)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	_, err = io.Copy(gz, r)
	gz.Close()

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `")

func %sBytes() ([]byte, error) {
	return bhojpurRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func uncompressed_nomemcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = "`, asset.Func)
	if err != nil {
		return err
	}

	_, err = io.Copy(&StringWriter{Writer: w}, r)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `"

func %sBytes() ([]byte, error) {
	return bhojpurRead(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
	return err
}

func uncompressed_memcopy(w io.Writer, asset *Asset, r io.Reader) error {
	_, err := fmt.Fprintf(w, `var _%s = []byte(`, asset.Func)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if utf8.Valid(b) && !bytes.Contains(b, []byte{0}) {
		fmt.Fprintf(w, "`%s`", sanitize(b))
	} else {
		fmt.Fprintf(w, "%+q", b)
	}

	_, err = fmt.Fprintf(w, `)

func %sBytes() ([]byte, error) {
	return _%s, nil
}

`, asset.Func, asset.Func)
	return err
}

func asset_release_common(w io.Writer, c *Config, asset *Asset) error {
	fi, err := os.Stat(asset.Path)
	if err != nil {
		return err
	}

	mode := uint(fi.Mode())
	modTime := fi.ModTime().Unix()
	size := fi.Size()
	if c.NoMetadata {
		mode = 0
		modTime = 0
		size = 0
	}
	if c.Mode > 0 {
		mode = uint(os.ModePerm) & c.Mode
	}
	if c.ModTime > 0 {
		modTime = c.ModTime
	}
	_, err = fmt.Fprintf(w, `func %s() (*asset, error) {
	bytes, err := %sBytes()
	if err != nil {
		return nil, err
	}

	info := bhojpurFileInfo{name: %q, size: %d, mode: os.FileMode(%d), modTime: time.Unix(%d, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

`, asset.Func, asset.Func, asset.Name, size, mode, modTime)
	return err
}
