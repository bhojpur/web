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
	"github.com/pkg/errors"

	"github.com/bhojpur/web/pkg/core/config"
)

type newToOldConfigureAdapter struct {
	delegate config.Configure
}

func (c *newToOldConfigureAdapter) Set(key, val string) error {
	return c.delegate.Set(key, val)
}

func (c *newToOldConfigureAdapter) String(key string) string {
	res, _ := c.delegate.String(key)
	return res
}

func (c *newToOldConfigureAdapter) Strings(key string) []string {
	res, _ := c.delegate.Strings(key)
	return res
}

func (c *newToOldConfigureAdapter) Int(key string) (int, error) {
	return c.delegate.Int(key)
}

func (c *newToOldConfigureAdapter) Int64(key string) (int64, error) {
	return c.delegate.Int64(key)
}

func (c *newToOldConfigureAdapter) Bool(key string) (bool, error) {
	return c.delegate.Bool(key)
}

func (c *newToOldConfigureAdapter) Float(key string) (float64, error) {
	return c.delegate.Float(key)
}

func (c *newToOldConfigureAdapter) DefaultString(key string, defaultVal string) string {
	return c.delegate.DefaultString(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DefaultStrings(key string, defaultVal []string) []string {
	return c.delegate.DefaultStrings(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DefaultInt(key string, defaultVal int) int {
	return c.delegate.DefaultInt(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DefaultInt64(key string, defaultVal int64) int64 {
	return c.delegate.DefaultInt64(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DefaultBool(key string, defaultVal bool) bool {
	return c.delegate.DefaultBool(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DefaultFloat(key string, defaultVal float64) float64 {
	return c.delegate.DefaultFloat(key, defaultVal)
}

func (c *newToOldConfigureAdapter) DIY(key string) (interface{}, error) {
	return c.delegate.DIY(key)
}

func (c *newToOldConfigureAdapter) GetSection(section string) (map[string]string, error) {
	return c.delegate.GetSection(section)
}

func (c *newToOldConfigureAdapter) SaveConfigFile(filename string) error {
	return c.delegate.SaveConfigFile(filename)
}

type oldToNewConfigureAdapter struct {
	delegate Configure
}

func (o *oldToNewConfigureAdapter) Set(key, val string) error {
	return o.delegate.Set(key, val)
}

func (o *oldToNewConfigureAdapter) String(key string) (string, error) {
	return o.delegate.String(key), nil
}

func (o *oldToNewConfigureAdapter) Strings(key string) ([]string, error) {
	return o.delegate.Strings(key), nil
}

func (o *oldToNewConfigureAdapter) Int(key string) (int, error) {
	return o.delegate.Int(key)
}

func (o *oldToNewConfigureAdapter) Int64(key string) (int64, error) {
	return o.delegate.Int64(key)
}

func (o *oldToNewConfigureAdapter) Bool(key string) (bool, error) {
	return o.delegate.Bool(key)
}

func (o *oldToNewConfigureAdapter) Float(key string) (float64, error) {
	return o.delegate.Float(key)
}

func (o *oldToNewConfigureAdapter) DefaultString(key string, defaultVal string) string {
	return o.delegate.DefaultString(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DefaultStrings(key string, defaultVal []string) []string {
	return o.delegate.DefaultStrings(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DefaultInt(key string, defaultVal int) int {
	return o.delegate.DefaultInt(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DefaultInt64(key string, defaultVal int64) int64 {
	return o.delegate.DefaultInt64(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DefaultBool(key string, defaultVal bool) bool {
	return o.delegate.DefaultBool(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DefaultFloat(key string, defaultVal float64) float64 {
	return o.delegate.DefaultFloat(key, defaultVal)
}

func (o *oldToNewConfigureAdapter) DIY(key string) (interface{}, error) {
	return o.delegate.DIY(key)
}

func (o *oldToNewConfigureAdapter) GetSection(section string) (map[string]string, error) {
	return o.delegate.GetSection(section)
}

func (o *oldToNewConfigureAdapter) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	return errors.New("unsupported operation, please use actual config.Configure")
}

func (o *oldToNewConfigureAdapter) Sub(key string) (config.Configure, error) {
	return nil, errors.New("unsupported operation, please use actual config.Configure")
}

func (o *oldToNewConfigureAdapter) OnChange(key string, fn func(value string)) {
	// do nothing
}

func (o *oldToNewConfigureAdapter) SaveConfigFile(filename string) error {
	return o.delegate.SaveConfigFile(filename)
}

type oldToNewConfigAdapter struct {
	delegate Config
}

func (o *oldToNewConfigAdapter) Parse(key string) (config.Configure, error) {
	old, err := o.delegate.Parse(key)
	if err != nil {
		return nil, err
	}
	return &oldToNewConfigureAdapter{delegate: old}, nil
}

func (o *oldToNewConfigAdapter) ParseData(data []byte) (config.Configure, error) {
	old, err := o.delegate.ParseData(data)
	if err != nil {
		return nil, err
	}
	return &oldToNewConfigureAdapter{delegate: old}, nil
}
