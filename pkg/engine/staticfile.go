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
	"bytes"
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"

	logsvr "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/context"
)

var errNotStaticRequest = errors.New("request not a static file request")

func serverStaticRouter(ctx *context.Context) {
	if ctx.Input.Method() != "GET" && ctx.Input.Method() != "HEAD" {
		return
	}

	forbidden, filePath, fileInfo, err := lookupFile(ctx)
	if err == errNotStaticRequest {
		return
	}

	if forbidden {
		exception("403", ctx)
		return
	}

	if filePath == "" || fileInfo == nil {
		if BConfig.RunMode == DEV {
			logsvr.Warn("Can't find/open the file:", filePath, err)
		}
		http.NotFound(ctx.ResponseWriter, ctx.Request)
		return
	}
	if fileInfo.IsDir() {
		requestURL := ctx.Input.URL()
		if requestURL[len(requestURL)-1] != '/' {
			redirectURL := requestURL + "/"
			if ctx.Request.URL.RawQuery != "" {
				redirectURL = redirectURL + "?" + ctx.Request.URL.RawQuery
			}
			ctx.Redirect(302, redirectURL)
		} else {
			// serveFile will list dir
			http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
		}
		return
	} else if fileInfo.Size() > int64(BConfig.WebConfig.StaticCacheFileSize) {
		// over size file serve with http module
		http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
		return
	}

	enableCompress := BConfig.EnableGzip && isStaticCompress(filePath)
	var acceptEncoding string
	if enableCompress {
		acceptEncoding = context.ParseEncoding(ctx.Request)
	}
	b, n, sch, reader, err := openFile(filePath, fileInfo, acceptEncoding)
	if err != nil {
		if BConfig.RunMode == DEV {
			logsvr.Warn("Can't compress the file:", filePath, err)
		}
		http.NotFound(ctx.ResponseWriter, ctx.Request)
		return
	}

	if b {
		ctx.Output.Header("Content-Encoding", n)
	} else {
		ctx.Output.Header("Content-Length", strconv.FormatInt(sch.size, 10))
	}

	http.ServeContent(ctx.ResponseWriter, ctx.Request, filePath, sch.modTime, reader)
}

type serveContentHolder struct {
	data       []byte
	modTime    time.Time
	size       int64
	originSize int64 // original file size:to judge file changed
	encoding   string
}

type serveContentReader struct {
	*bytes.Reader
}

var (
	staticFileLruCache *lru.Cache
	lruLock            sync.RWMutex
)

func openFile(filePath string, fi os.FileInfo, acceptEncoding string) (bool, string, *serveContentHolder, *serveContentReader, error) {
	if staticFileLruCache == nil {
		// avoid lru cache error
		if BConfig.WebConfig.StaticCacheFileNum >= 1 {
			staticFileLruCache, _ = lru.New(BConfig.WebConfig.StaticCacheFileNum)
		} else {
			staticFileLruCache, _ = lru.New(1)
		}
	}
	mapKey := acceptEncoding + ":" + filePath
	lruLock.RLock()
	var mapFile *serveContentHolder
	if cacheItem, ok := staticFileLruCache.Get(mapKey); ok {
		mapFile = cacheItem.(*serveContentHolder)
	}
	lruLock.RUnlock()
	if isOk(mapFile, fi) {
		reader := &serveContentReader{Reader: bytes.NewReader(mapFile.data)}
		return mapFile.encoding != "", mapFile.encoding, mapFile, reader, nil
	}
	lruLock.Lock()
	defer lruLock.Unlock()
	if cacheItem, ok := staticFileLruCache.Get(mapKey); ok {
		mapFile = cacheItem.(*serveContentHolder)
	}
	if !isOk(mapFile, fi) {
		file, err := os.Open(filePath)
		if err != nil {
			return false, "", nil, nil, err
		}
		defer file.Close()
		var bufferWriter bytes.Buffer
		_, n, err := context.WriteFile(acceptEncoding, &bufferWriter, file)
		if err != nil {
			return false, "", nil, nil, err
		}
		mapFile = &serveContentHolder{data: bufferWriter.Bytes(), modTime: fi.ModTime(), size: int64(bufferWriter.Len()), originSize: fi.Size(), encoding: n}
		if isOk(mapFile, fi) {
			staticFileLruCache.Add(mapKey, mapFile)
		}
	}

	reader := &serveContentReader{Reader: bytes.NewReader(mapFile.data)}
	return mapFile.encoding != "", mapFile.encoding, mapFile, reader, nil
}

func isOk(s *serveContentHolder, fi os.FileInfo) bool {
	if s == nil {
		return false
	} else if s.size > int64(BConfig.WebConfig.StaticCacheFileSize) {
		return false
	}
	return s.modTime == fi.ModTime() && s.originSize == fi.Size()
}

// isStaticCompress detect static files
func isStaticCompress(filePath string) bool {
	for _, statExtension := range BConfig.WebConfig.StaticExtensionsToGzip {
		if strings.HasSuffix(strings.ToLower(filePath), strings.ToLower(statExtension)) {
			return true
		}
	}
	return false
}

// searchFile search the file by url path
// if none the static file prefix matches ,return notStaticRequestErr
func searchFile(ctx *context.Context) (string, os.FileInfo, error) {
	requestPath := filepath.ToSlash(filepath.Clean(ctx.Request.URL.Path))
	// special processing : favicon.ico/robots.txt  can be in any static dir
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
		file := path.Join(".", requestPath)
		if fi, _ := os.Stat(file); fi != nil {
			return file, fi, nil
		}
		for _, staticDir := range BConfig.WebConfig.StaticDir {
			filePath := path.Join(staticDir, requestPath)
			if fi, _ := os.Stat(filePath); fi != nil {
				return filePath, fi, nil
			}
		}
		return "", nil, errNotStaticRequest
	}

	for prefix, staticDir := range BConfig.WebConfig.StaticDir {
		if !strings.Contains(requestPath, prefix) {
			continue
		}
		if prefix != "/" && len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
			continue
		}
		filePath := path.Join(staticDir, requestPath[len(prefix):])
		if fi, err := os.Stat(filePath); fi != nil {
			return filePath, fi, err
		}
	}
	return "", nil, errNotStaticRequest
}

// lookupFile find the file to serve
// if the file is dir ,search the index.html as default file( MUST NOT A DIR also)
// if the index.html not exist or is a dir, give a forbidden response depending on  DirectoryIndex
func lookupFile(ctx *context.Context) (bool, string, os.FileInfo, error) {
	fp, fi, err := searchFile(ctx)
	if fp == "" || fi == nil {
		return false, "", nil, err
	}
	if !fi.IsDir() {
		return false, fp, fi, err
	}
	if requestURL := ctx.Input.URL(); requestURL[len(requestURL)-1] == '/' {
		ifp := filepath.Join(fp, "index.html")
		if ifi, _ := os.Stat(ifp); ifi != nil && ifi.Mode().IsRegular() {
			return false, ifp, ifi, err
		}
	}
	return !BConfig.WebConfig.DirectoryIndex, fp, fi, err
}
