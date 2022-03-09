package template

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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// FSLoader supports the fs.FS interface for loading templates
type FSLoader struct {
	fs fs.FS
}

func NewFSLoader(fs fs.FS) *FSLoader {
	return &FSLoader{
		fs: fs,
	}
}

func (l *FSLoader) Abs(base, name string) string {
	return filepath.Join(filepath.Dir(base), name)
}

func (l *FSLoader) Get(path string) (io.Reader, error) {
	return l.fs.Open(path)
}

// LocalFilesystemLoader represents a local filesystem loader with basic
// BaseDirectory capabilities. The access to the local filesystem is unrestricted.
type LocalFilesystemLoader struct {
	baseDir string
}

// MustNewLocalFileSystemLoader creates a new LocalFilesystemLoader instance
// and panics if there's any error during instantiation. The parameters
// are the same like NewLocalFileSystemLoader.
func MustNewLocalFileSystemLoader(baseDir string) *LocalFilesystemLoader {
	fs, err := NewLocalFileSystemLoader(baseDir)
	if err != nil {
		log.Panic(err)
	}
	return fs
}

// NewLocalFileSystemLoader creates a new LocalFilesystemLoader and allows
// templatesto be loaded from disk (unrestricted). If any base directory
// is given (or being set using SetBaseDir), this base directory is being used
// for path calculation in template inclusions/imports. Otherwise the path
// is calculated based relatively to the including template's path.
func NewLocalFileSystemLoader(baseDir string) (*LocalFilesystemLoader, error) {
	fs := &LocalFilesystemLoader{}
	if baseDir != "" {
		if err := fs.SetBaseDir(baseDir); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

// SetBaseDir sets the template's base directory. This directory will
// be used for any relative path in filters, tags and From*-functions to determine
// your template. See the comment for NewLocalFileSystemLoader as well.
func (fs *LocalFilesystemLoader) SetBaseDir(path string) error {
	// Make the path absolute
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		path = abs
	}

	// Check for existence
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("the given path '%s' is not a directory", path)
	}

	fs.baseDir = path
	return nil
}

// Get reads the path's content from your local filesystem.
func (fs *LocalFilesystemLoader) Get(path string) (io.Reader, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

// Abs resolves a filename relative to the base directory. Absolute paths are allowed.
// When there's no base dir set, the absolute path to the filename
// will be calculated based on either the provided base directory (which
// might be a path of a template which includes another template) or
// the current working directory.
func (fs *LocalFilesystemLoader) Abs(base, name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	// Our own base dir has always priority; if there's none
	// we use the path provided in base.
	var err error
	if fs.baseDir == "" {
		if base == "" {
			base, err = os.Getwd()
			if err != nil {
				panic(err)
			}
			return filepath.Join(base, name)
		}

		return filepath.Join(filepath.Dir(base), name)
	}

	return filepath.Join(fs.baseDir, name)
}

// SandboxedFilesystemLoader is still WIP.
type SandboxedFilesystemLoader struct {
	*LocalFilesystemLoader
}

// NewSandboxedFilesystemLoader creates a new sandboxed local file system instance.
func NewSandboxedFilesystemLoader(baseDir string) (*SandboxedFilesystemLoader, error) {
	fs, err := NewLocalFileSystemLoader(baseDir)
	if err != nil {
		return nil, err
	}
	return &SandboxedFilesystemLoader{
		LocalFilesystemLoader: fs,
	}, nil
}

// Move sandbox to a virtual fs

/*
if len(set.SandboxDirectories) > 0 {
    defer func() {
        // Remove any ".." or other crap
        resolvedPath = filepath.Clean(resolvedPath)

        // Make the path absolute
        absPath, err := filepath.Abs(resolvedPath)
        if err != nil {
            panic(err)
        }
        resolvedPath = absPath

        // Check against the sandbox directories (once one pattern matches, we're done and can allow it)
        for _, pattern := range set.SandboxDirectories {
            matched, err := filepath.Match(pattern, resolvedPath)
            if err != nil {
                panic("Wrong sandbox directory match pattern (see http://golang.org/pkg/path/filepath/#Match).")
            }
            if matched {
                // OK!
                return
            }
        }

        // No pattern matched, we have to log+deny the request
        set.logf("Access attempt outside of the sandbox directories (blocked): '%s'", resolvedPath)
        resolvedPath = ""
    }()
}
*/

// HttpFilesystemLoader supports loading templates
// from an http.FileSystem - useful for using one of several
// file-to-code generators that packs static files into
// a go binary (ex: https://github.com/jteeuwen/go-bindata)
type HttpFilesystemLoader struct {
	fs      http.FileSystem
	baseDir string
}

// MustNewHttpFileSystemLoader creates a new HttpFilesystemLoader instance
// and panics if there's any error during instantiation. The parameters
// are the same like NewHttpFilesystemLoader.
func MustNewHttpFileSystemLoader(httpfs http.FileSystem, baseDir string) *HttpFilesystemLoader {
	fs, err := NewHttpFileSystemLoader(httpfs, baseDir)
	if err != nil {
		log.Panic(err)
	}
	return fs
}

// NewHttpFileSystemLoader creates a new HttpFileSystemLoader and allows
// templates to be loaded from the virtual filesystem. The path
// is calculated based relatively from the root of the http.Filesystem.
func NewHttpFileSystemLoader(httpfs http.FileSystem, baseDir string) (*HttpFilesystemLoader, error) {
	hfs := &HttpFilesystemLoader{
		fs:      httpfs,
		baseDir: baseDir,
	}
	if httpfs == nil {
		err := errors.New("httpfs cannot be nil")
		return nil, err
	}
	return hfs, nil
}

// Abs in this instance simply returns the filename, since
// there's no potential for an unexpanded path in an http.FileSystem
func (h *HttpFilesystemLoader) Abs(base, name string) string {
	return name
}

// Get returns an io.Reader where the template's content can be read from.
func (h *HttpFilesystemLoader) Get(path string) (io.Reader, error) {
	fullPath := path
	if h.baseDir != "" {
		fullPath = fmt.Sprintf(
			"%s/%s",
			h.baseDir,
			fullPath,
		)
	}

	return h.fs.Open(fullPath)
}