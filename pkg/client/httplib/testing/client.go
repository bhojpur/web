package testing

import (
	"github.com/bhojpur/web/pkg/client/httplib"
)

var port = ""
var baseURL = "http://localhost:"

// TestHTTPRequest bhojpur test request client
type TestHTTPRequest struct {
	httplib.BhojpurHTTPRequest
}

func SetTestingPort(p string) {
	port = p
}

func getPort() string {
	if port == "" {
		port = "8080"
		return port
	}
	return port
}

// Get returns test client in GET method
func Get(path string) *TestHTTPRequest {
	return &TestHTTPRequest{*httplib.Get(baseURL + getPort() + path)}
}

// Post returns test client in POST method
func Post(path string) *TestHTTPRequest {
	return &TestHTTPRequest{*httplib.Post(baseURL + getPort() + path)}
}

// Put returns test client in PUT method
func Put(path string) *TestHTTPRequest {
	return &TestHTTPRequest{*httplib.Put(baseURL + getPort() + path)}
}

// Delete returns test client in DELETE method
func Delete(path string) *TestHTTPRequest {
	return &TestHTTPRequest{*httplib.Delete(baseURL + getPort() + path)}
}

// Head returns test client in HEAD method
func Head(path string) *TestHTTPRequest {
	return &TestHTTPRequest{*httplib.Head(baseURL + getPort() + path)}
}
