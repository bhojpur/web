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

// Usage:
//  import "github.com/bhojpur/web/pkg/core/config"
//
// Examples.
//
//  cnf, err := config.NewConfig("ini", "config.conf")
//
//  cnf APIS:
//
//  cnf.Set(key, val string) error
//  cnf.String(key string) string
//  cnf.Strings(key string) []string
//  cnf.Int(key string) (int, error)
//  cnf.Int64(key string) (int64, error)
//  cnf.Bool(key string) (bool, error)
//  cnf.Float(key string) (float64, error)
//  cnf.DefaultString(key string, defaultVal string) string
//  cnf.DefaultStrings(key string, defaultVal []string) []string
//  cnf.DefaultInt(key string, defaultVal int) int
//  cnf.DefaultInt64(key string, defaultVal int64) int64
//  cnf.DefaultBool(key string, defaultVal bool) bool
//  cnf.DefaultFloat(key string, defaultVal float64) float64
//  cnf.DIY(key string) (interface{}, error)
//  cnf.GetSection(section string) (map[string]string, error)
//  cnf.SaveConfigFile(filename string) error

import (
	"github.com/bhojpur/web/pkg/core/config"
)

// Configure defines how to get and set value from configuration raw data.
type Configure interface {
	Set(key, val string) error   // support section::key type in given key when using ini type.
	String(key string) string    // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string // get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultVal []string) []string // get string slice
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

// Config is the Adapter interface for parsing config file to get raw data to Configer.
type Config interface {
	Parse(key string) (Configure, error)
	ParseData(data []byte) (Configure, error)
}

var adapters = make(map[string]Config)

// Register makes a config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Config) {
	config.Register(name, &oldToNewConfigAdapter{delegate: adapter})
}

// NewConfig adapterName is ini/json/xml/yaml.
// filename is the config file path.
func NewConfig(adapterName, filename string) (Configure, error) {
	cfg, err := config.NewConfig(adapterName, filename)
	if err != nil {
		return nil, err
	}

	// it was registered by using Register method
	res, ok := cfg.(*oldToNewConfigureAdapter)
	if ok {
		return res.delegate, nil
	}

	return &newToOldConfigureAdapter{
		delegate: cfg,
	}, nil
}

// NewConfigData adapterName is ini/json/xml/yaml.
// data is the config data.
func NewConfigData(adapterName string, data []byte) (Configure, error) {
	cfg, err := config.NewConfigData(adapterName, data)
	if err != nil {
		return nil, err
	}

	// it was registered by using Register method
	res, ok := cfg.(*oldToNewConfigureAdapter)
	if ok {
		return res.delegate, nil
	}

	return &newToOldConfigureAdapter{
		delegate: cfg,
	}, nil
}

// ExpandValueEnvForMap convert all string value with environment variable.
func ExpandValueEnvForMap(m map[string]interface{}) map[string]interface{} {
	return config.ExpandValueEnvForMap(m)
}

// ExpandValueEnv returns value of convert with environment variable.
//
// Return environment variable if value start with "${" and end with "}".
// Return default value if environment variable is empty or not exist.
//
// It accept value formats "${env}" , "${env||}}" , "${env||defaultValue}" , "defaultvalue".
// Examples:
//	v1 := config.ExpandValueEnv("${GOPATH}")			// return the GOPATH environment variable.
//	v2 := config.ExpandValueEnv("${GOBhojpur||/usr/local/go}")	// return the default value "/usr/local/go/".
//	v3 := config.ExpandValueEnv("Bhojpur")				// return the value "Bhojpur".
func ExpandValueEnv(value string) string {
	return config.ExpandValueEnv(value)
}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	return config.ParseBool(val)
}

// ToString converts values of any type to string.
func ToString(x interface{}) string {
	return config.ToString(x)
}
