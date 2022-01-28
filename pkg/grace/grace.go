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

// Package grace use to hot reload
//
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
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
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
	regLock              *sync.Mutex
	runningServers       map[string]*Server
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	// DefaultReadTimeOut is the HTTP read timeout
	DefaultReadTimeOut time.Duration
	// DefaultWriteTimeOut is the HTTP Write timeout
	DefaultWriteTimeOut time.Duration
	// DefaultMaxHeaderBytes is the Max HTTP Header size, default is 0, no limit
	DefaultMaxHeaderBytes int
	// DefaultTimeout is the shutdown server's timeout. default is 60s
	DefaultTimeout = 60 * time.Second

	isChild     bool
	socketOrder string

	hookableSignals []os.Signal
)

func init() {
	flag.BoolVar(&isChild, "graceful", false, "listen on open fd (after forking)")
	flag.StringVar(&socketOrder, "socketorder", "", "previous initialization order - used when more than one listener was started")

	regLock = &sync.Mutex{}
	runningServers = make(map[string]*Server)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
}

// NewServer returns a new graceServer.
func NewServer(addr string, handler http.Handler) (srv *Server) {
	regLock.Lock()
	defer regLock.Unlock()

	if !flag.Parsed() {
		flag.Parse()
	}
	if len(socketOrder) > 0 {
		for i, addr := range strings.Split(socketOrder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &Server{
		sigChan: make(chan os.Signal),
		isChild: isChild,
		SignalHooks: map[int]map[os.Signal][]func(){
			PreSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
			PostSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
		},
		state:        StateInit,
		Network:      "tcp",
		terminalChan: make(chan error), //no cache channel
	}
	srv.Server = &http.Server{
		Addr:           addr,
		ReadTimeout:    DefaultReadTimeOut,
		WriteTimeout:   DefaultWriteTimeOut,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
		Handler:        handler,
	}

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv
	return srv
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
