package xml

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

// It is for config provider.
//
// depend on github.com/bhojpur/web/pkg/core/x2j.
//
// go install github.com/bhojpur/web/pkg/core/x2j.
//
// Usage:
//  import(
//    _ "github.com/bhojpur/web/pkg/core/config/xml"
//      "github.com/bhojpur/web/pkg/core/config"
//  )
//
//  cnf, err := config.NewConfig("xml", "config.xml")

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"

	logs "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/core/config"

	"github.com/bhojpur/web/pkg/core/x2j"
)

// Config is a xml config parser and implements Config interface.
// xml configurations should be included in <config></config> tag.
// only support key/value pair as <key>value</key> as each item.
type Config struct{}

// Parse returns a ConfigContainer with parsed xml config map.
func (xc *Config) Parse(filename string) (config.Configure, error) {
	context, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return xc.ParseData(context)
}

// ParseData xml data
func (xc *Config) ParseData(data []byte) (config.Configure, error) {
	x := &ConfigContainer{data: make(map[string]interface{})}

	d, err := x2j.DocToMap(string(data))
	if err != nil {
		return nil, err
	}

	x.data = config.ExpandValueEnvForMap(d["config"].(map[string]interface{}))

	return x, nil
}

// ConfigContainer is a Config which represents the xml configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

// Unmarshaler is a little be inconvenient since the xml library doesn't know type.
// So when you use
// <id>1</id>
// The "1" is a string, not int
func (c *ConfigContainer) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	sub, err := c.sub(prefix)
	if err != nil {
		return err
	}
	return mapstructure.Decode(sub, obj)
}

func (c *ConfigContainer) Sub(key string) (config.Configure, error) {
	sub, err := c.sub(key)
	if err != nil {
		return nil, err
	}

	return &ConfigContainer{
		data: sub,
	}, nil

}

func (c *ConfigContainer) sub(key string) (map[string]interface{}, error) {
	if key == "" {
		return c.data, nil
	}
	value, ok := c.data[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("the key is not found: %s", key))
	}
	res, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("the value of this key is not a structure: %s", key))
	}
	return res, nil
}

func (c *ConfigContainer) OnChange(key string, fn func(value string)) {
	logs.Warn("Unsupported operation")
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {
	if v := c.data[key]; v != nil {
		return config.ParseBool(v)
	}
	return false, fmt.Errorf("not exist key: %q", key)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultVal
func (c *ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// Int returns the integer value for a given key.
func (c *ConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key].(string))
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *ConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key].(string), 10, 64)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultVal
	}
	return v

}

// Float returns the float value for a given key.
func (c *ConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key].(string), 64)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// String returns the string value for a given key.
func (c *ConfigContainer) String(key string) (string, error) {
	if v, ok := c.data[key].(string); ok {
		return v, nil
	}
	return "", nil
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	v, err := c.String(key)
	if v == "" || err != nil {
		return defaultVal
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigContainer) Strings(key string) ([]string, error) {
	v, err := c.String(key)
	if v == "" || err != nil {
		return nil, err
	}
	return strings.Split(v, ";"), nil
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	v, err := c.Strings(key)
	if v == nil || err != nil {
		return defaultVal
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section].(map[string]interface{}); ok {
		mapstr := make(map[string]string)
		for k, val := range v {
			mapstr[k] = config.ToString(val)
		}
		return mapstr, nil
	}
	return nil, fmt.Errorf("section '%s' not found", section)
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := xml.MarshalIndent(c.data, "  ", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

// Set writes a new value for key.
func (c *ConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *ConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	config.Register("xml", &Config{})
}
