package httplib

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
	"io"
	"io/ioutil"
	"net/http"
)

// Client provides an HTTP client supporting chain call
type Client struct {
	Name       string
	Endpoint   string
	CommonOpts []BhojpurHTTPRequestOption

	Setting BhojpurHTTPSettings
}

// HTTPResponseCarrier If value implement HTTPResponseCarrier. http.Response will pass to SetHTTPResponse
type HTTPResponseCarrier interface {
	SetHTTPResponse(resp *http.Response)
}

// HTTPBodyCarrier If value implement HTTPBodyCarrier. http.Response.Body will pass to SetReader
type HTTPBodyCarrier interface {
	SetReader(r io.ReadCloser)
}

// HTTPBytesCarrier If value implement HTTPBytesCarrier.
// All the byte in http.Response.Body will pass to SetBytes
type HTTPBytesCarrier interface {
	SetBytes(bytes []byte)
}

// HTTPStatusCarrier If value implement HTTPStatusCarrier. http.Response.StatusCode will pass to SetStatusCode
type HTTPStatusCarrier interface {
	SetStatusCode(status int)
}

// HttpHeaderCarrier If value implement HttpHeaderCarrier. http.Response.Header will pass to SetHeader
type HTTPHeadersCarrier interface {
	SetHeader(header map[string][]string)
}

// NewClient return a new http client
func NewClient(name string, endpoint string, opts ...ClientOption) (*Client, error) {
	res := &Client{
		Name:     name,
		Endpoint: endpoint,
	}
	setting := GetDefaultSetting()
	res.Setting = setting
	for _, o := range opts {
		o(res)
	}
	return res, nil
}

func (c *Client) customReq(req *BhojpurHTTPRequest, opts []BhojpurHTTPRequestOption) {
	req.Setting(c.Setting)
	opts = append(c.CommonOpts, opts...)
	for _, o := range opts {
		o(req)
	}
}

// handleResponse try to parse body to meaningful value
func (c *Client) handleResponse(value interface{}, req *BhojpurHTTPRequest) error {
	// make sure req.resp is not nil
	_, err := req.Bytes()
	if err != nil {
		return err
	}

	err = c.handleCarrier(value, req)
	if err != nil {
		return err
	}

	return req.ToValue(value)
}

// handleCarrier set http data to value
func (c *Client) handleCarrier(value interface{}, req *BhojpurHTTPRequest) error {
	if value == nil {
		return nil
	}

	if carrier, ok := value.(HTTPResponseCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		req.resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetHTTPResponse(req.resp)
	}
	if carrier, ok := value.(HTTPBodyCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		reader := ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetReader(reader)
	}
	if carrier, ok := value.(HTTPBytesCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		carrier.SetBytes(b)
	}
	if carrier, ok := value.(HTTPStatusCarrier); ok {
		carrier.SetStatusCode(req.resp.StatusCode)
	}
	if carrier, ok := value.(HTTPHeadersCarrier); ok {
		carrier.SetHeader(req.resp.Header)
	}
	return nil
}

// Get Send a GET request and try to give its result value
func (c *Client) Get(value interface{}, path string, opts ...BhojpurHTTPRequestOption) error {
	req := Get(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Post Send a POST request and try to give its result value
func (c *Client) Post(value interface{}, path string, body interface{}, opts ...BhojpurHTTPRequestOption) error {
	req := Post(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Put Send a Put request and try to give its result value
func (c *Client) Put(value interface{}, path string, body interface{}, opts ...BhojpurHTTPRequestOption) error {
	req := Put(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Delete Send a Delete request and try to give its result value
func (c *Client) Delete(value interface{}, path string, opts ...BhojpurHTTPRequestOption) error {
	req := Delete(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Head Send a Head request and try to give its result value
func (c *Client) Head(value interface{}, path string, opts ...BhojpurHTTPRequestOption) error {
	req := Head(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}
