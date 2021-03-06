package grace

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

// Usage:
//
// import(
//   "log"
//	 "net/http"
//	 "os"
//
//   "github.com/bhojpur/web/pkg/grace"
// )
//
//  func handler(w http.ResponseWriter, r *http.Request) {
//	  w.Write([]byte("WORLD!"))
//  }
//
//  func main() {
//      mux := http.NewServeMux()
//      mux.HandleFunc("/hello", handler)
//
//	    err := grace.ListenAndServe("localhost:8080", mux)
//      if err != nil {
//		   log.Println(err)
//	    }
//      log.Println("Server on 8080 stopped")
//	     os.Exit(0)
//    }

import (
	"net/http"
	"time"

	"github.com/bhojpur/web/pkg/grace"
)

const (
	// PreSignal is the position to add filter before signal
	PreSignal = iota
	// PostSignal is the position to add filter after signal
	PostSignal
	// StateInit represent the application inited
	StateInit
	// StateRunning represent the application is running
	StateRunning
	// StateShuttingDown represent the application is shutting down
	StateShuttingDown
	// StateTerminate represent the application is killed
	StateTerminate
)

var (

	// DefaultReadTimeOut is the HTTP read timeout
	DefaultReadTimeOut time.Duration
	// DefaultWriteTimeOut is the HTTP Write timeout
	DefaultWriteTimeOut time.Duration
	// DefaultMaxHeaderBytes is the Max HTTP Header size, default is 0, no limit
	DefaultMaxHeaderBytes int
	// DefaultTimeout is the shutdown server's timeout. default is 60s
	DefaultTimeout = grace.DefaultTimeout
)

// NewServer returns a new graceServer.
func NewServer(addr string, handler http.Handler) (srv *Server) {
	return (*Server)(grace.NewServer(addr, handler))
}

// ListenAndServe refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

// ListenAndServeTLS refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}
