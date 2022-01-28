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

// We use this to simply application's development
// for most users, they only need to use those methods
var globalInstance Configure

// InitGlobalInstance will ini the global instance
// If you want to use specific implementation, don't forget to import it.
// e.g. _ import "github.com/bhojpur/web/pkg/core/config/etcd"
// err := InitGlobalInstance("etcd", "someconfig")
func InitGlobalInstance(name string, cfg string) error {
	var err error
	globalInstance, err = NewConfig(name, cfg)
	return err
}

// support section::key type in given key when using ini type.
func Set(key, val string) error {
	return globalInstance.Set(key, val)
}

// support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
func String(key string) (string, error) {
	return globalInstance.String(key)
}

// get string slice
func Strings(key string) ([]string, error) {
	return globalInstance.Strings(key)
}
func Int(key string) (int, error) {
	return globalInstance.Int(key)
}
func Int64(key string) (int64, error) {
	return globalInstance.Int64(key)
}
func Bool(key string) (bool, error) {
	return globalInstance.Bool(key)
}
func Float(key string) (float64, error) {
	return globalInstance.Float(key)
}

// support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
func DefaultString(key string, defaultVal string) string {
	return globalInstance.DefaultString(key, defaultVal)
}

// get string slice
func DefaultStrings(key string, defaultVal []string) []string {
	return globalInstance.DefaultStrings(key, defaultVal)
}
func DefaultInt(key string, defaultVal int) int {
	return globalInstance.DefaultInt(key, defaultVal)
}
func DefaultInt64(key string, defaultVal int64) int64 {
	return globalInstance.DefaultInt64(key, defaultVal)
}
func DefaultBool(key string, defaultVal bool) bool {
	return globalInstance.DefaultBool(key, defaultVal)
}
func DefaultFloat(key string, defaultVal float64) float64 {
	return globalInstance.DefaultFloat(key, defaultVal)
}

// DIY return the original value
func DIY(key string) (interface{}, error) {
	return globalInstance.DIY(key)
}

func GetSection(section string) (map[string]string, error) {
	return globalInstance.GetSection(section)
}

func Unmarshaler(prefix string, obj interface{}, opt ...DecodeOption) error {
	return globalInstance.Unmarshaler(prefix, obj, opt...)
}
func Sub(key string) (Configure, error) {
	return globalInstance.Sub(key)
}

func OnChange(key string, fn func(value string)) {
	globalInstance.OnChange(key, fn)
}

func SaveConfigFile(filename string) error {
	return globalInstance.SaveConfigFile(filename)
}
