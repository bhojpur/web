package virtual

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
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"testing"
	"time"
)

// TestErrors ensures that error types are implemented properly
func TestErrors(t *testing.T) {
	// create local err vars of error interface type
	var notFoundErr error
	var badRequestErr error
	var closedErr error

	// stuff local error types in to interface values
	notFoundErr = NotFound{fmt.Errorf("test NotFound")}
	badRequestErr = BadRequest{fmt.Errorf("test BadRequest")}
	closedErr = Closed{fmt.Errorf("test Closed")}

	// assert the types
	switch errType := notFoundErr.(type) {
	case NotFound:
	default:
		t.Fatalf("expected NotFound, got: %s", errType)
	}
	switch errType := badRequestErr.(type) {
	case BadRequest:
	default:
		t.Fatalf("expected BadRequest, got: %s", errType)
	}
	switch errType := closedErr.(type) {
	case Closed:
	default:
		t.Fatalf("expected Closed, got: %s", errType)
	}
}

func localListener(t *testing.T) (net.Listener, string) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	return l, strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

func TestHTTPMux(t *testing.T) {
	l, port := localListener(t)
	mux, err := NewHTTPMuxer(l, time.Second)
	if err != nil {
		t.Fatalf("failed to start muxer: %v", err)
	}
	go mux.HandleErrors()

	muxed, err := mux.Listen("bhojpur.net")
	if err != nil {
		t.Fatalf("failed to listen on muxer: %v", muxed)
	}

	go http.Serve(muxed, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	}))

	msg := "test"
	url := "http://localhost:" + port
	resp, err := http.Post(url, "text/plain", strings.NewReader(msg))
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("sent incorrect host header, expected 404 but got %d", resp.StatusCode)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(msg))
	if err != nil {
		t.Fatalf("failed to construct HTTP request: %v", err)
	}
	req.Host = "bhojpur.net"
	req.Header.Set("Content-Type", "text/plain")

	resp, err = new(http.Client).Do(req)
	if err != nil {
		t.Fatalf("failed to make HTTP request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	got := string(body)
	if got != msg {
		t.Fatalf("unexpected resposne. got: %v, expected: %v", got, msg)
	}
}

func testMux(t *testing.T, listen, dial string) {
	muxFn := func(c net.Conn) (Conn, error) {
		return fakeConn{c, dial}, nil
	}

	fakel := make(fakeListener, 1)
	mux, err := NewVirtualHostMuxer(fakel, muxFn, time.Second)
	if err != nil {
		t.Fatalf("failed to start virtual host muxer: %v", err)
	}

	l, err := mux.Listen(listen)
	if err != nil {
		t.Fatalf("failed to listen for %s", err)
	}

	done := make(chan struct{})
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("failed to accept connection: %v", err)
			return
		}

		got := conn.(Conn).Host()
		expected := dial
		if got != expected {
			t.Fatalf("got connection with unexpected host. got: %s, expected: %s", got, expected)
			return
		}

		close(done)
	}()

	go func() {
		_, err := mux.NextError()
		if err != nil {
			t.Fatalf("muxing error: %v", err)
		}
	}()

	fakel <- struct{}{}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("test timed out: dial: %s listen: %s", dial, listen)
	}
}

func TestMuxingPatterns(t *testing.T) {
	var tests = []struct {
		listen string
		dial   string
	}{
		{"bhojpur.net", "bhojpur.net"},
		{"sub.bhojpur.net", "sub.bhojpur.net"},
		{"*.bhojpur.net", "sub.bhojpur.net"},
		{"*.bhojpur.net", "nested.sub.bhojpur.net"},
	}

	for _, test := range tests {
		testMux(t, test.listen, test.dial)
	}
}

type fakeConn struct {
	net.Conn
	host string
}

func (c fakeConn) SetDeadline(d time.Time) error { return nil }
func (c fakeConn) Host() string                  { return c.host }
func (c fakeConn) Free()                         {}

type fakeNetConn struct {
	net.Conn
}

func (fakeNetConn) SetDeadline(time.Time) error { return nil }

type fakeListener chan struct{}

func (l fakeListener) Accept() (net.Conn, error) {
	for _ = range l {
		return fakeNetConn{nil}, nil
	}
	select {}
}
func (fakeListener) Addr() net.Addr { return nil }
func (fakeListener) Close() error   { return nil }
