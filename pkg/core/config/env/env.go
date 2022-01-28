package env

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
	"os"
	"strings"

	"github.com/bhojpur/web/pkg/core/utils"
)

var env *utils.BhojpurMap

func init() {
	env = utils.NewBhojpurMap()
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		env.Set(splits[0], os.Getenv(splits[0]))
	}
}

// Get returns a value for a given key.
// If the key does not exist, the default value will be returned.
func Get(key string, defVal string) string {
	if val := env.Get(key); val != nil {
		return val.(string)
	}
	return defVal
}

// MustGet returns a value by key.
// If the key does not exist, it will return an error.
func MustGet(key string) (string, error) {
	if val := env.Get(key); val != nil {
		return val.(string), nil
	}
	return "", fmt.Errorf("no env variable with %s", key)
}

// Set sets a value in the ENV copy.
// This does not affect the child process environment.
func Set(key string, value string) {
	env.Set(key, value)
}

// MustSet sets a value in the ENV copy and the child process environment.
// It returns an error in case the set operation failed.
func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env.Set(key, value)
	return nil
}

// GetAll returns all keys/values in the current child process environment.
func GetAll() map[string]string {
	items := env.Items()
	envs := make(map[string]string, env.Count())

	for key, val := range items {
		switch key := key.(type) {
		case string:
			switch val := val.(type) {
			case string:
				envs[key] = val
			}
		}
	}
	return envs
}
