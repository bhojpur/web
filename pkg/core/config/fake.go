package config

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
	"context"
	"errors"
	"strconv"
	"strings"
)

type fakeConfigContainer struct {
	BaseConfigure
	data map[string]string
}

func (c *fakeConfigContainer) getData(key string) string {
	return c.data[strings.ToLower(key)]
}

func (c *fakeConfigContainer) Set(key, val string) error {
	c.data[strings.ToLower(key)] = val
	return nil
}

func (c *fakeConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}

func (c *fakeConfigContainer) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}

func (c *fakeConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Bool(key string) (bool, error) {
	return ParseBool(c.getData(key))
}

func (c *fakeConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}

func (c *fakeConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) DIY(key string) (interface{}, error) {
	if v, ok := c.data[strings.ToLower(key)]; ok {
		return v, nil
	}
	return nil, errors.New("key not find")
}

func (c *fakeConfigContainer) GetSection(section string) (map[string]string, error) {
	return nil, errors.New("not implement in the fakeConfigContainer")
}

func (c *fakeConfigContainer) SaveConfigFile(filename string) error {
	return errors.New("not implement in the fakeConfigContainer")
}

func (c *fakeConfigContainer) Unmarshaler(prefix string, obj interface{}, opt ...DecodeOption) error {
	return errors.New("unsupported operation")
}

var _ Configure = new(fakeConfigContainer)

// NewFakeConfig return a fake Configure
func NewFakeConfig() Configure {
	res := &fakeConfigContainer{
		data: make(map[string]string),
	}
	res.BaseConfigure = NewBaseConfigure(func(ctx context.Context, key string) (string, error) {
		return res.getData(key), nil
	})
	return res
}
