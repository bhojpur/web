package json

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
	"encoding/json"
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
)

// JSONConfig is a json config parser and implements Config interface.
type JSONConfig struct {
}

// Parse returns a ConfigContainer with parsed json config map.
func (js *JSONConfig) Parse(filename string) (config.Configure, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return js.ParseData(content)
}

// ParseData returns a ConfigContainer with json string
func (js *JSONConfig) ParseData(data []byte) (config.Configure, error) {
	x := &JSONConfigContainer{
		data: make(map[string]interface{}),
	}
	err := json.Unmarshal(data, &x.data)
	if err != nil {
		var wrappingArray []interface{}
		err2 := json.Unmarshal(data, &wrappingArray)
		if err2 != nil {
			return nil, err
		}
		x.data["rootArray"] = wrappingArray
	}

	x.data = config.ExpandValueEnvForMap(x.data)

	return x, nil
}

// JSONConfigContainer is a config which represents the json configuration.
// Only when get value, support key as section:name type.
type JSONConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

func (c *JSONConfigContainer) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	sub, err := c.sub(prefix)
	if err != nil {
		return err
	}
	return mapstructure.Decode(sub, obj)
}

func (c *JSONConfigContainer) Sub(key string) (config.Configure, error) {
	sub, err := c.sub(key)
	if err != nil {
		return nil, err
	}
	return &JSONConfigContainer{
		data: sub,
	}, nil
}

func (c *JSONConfigContainer) sub(key string) (map[string]interface{}, error) {
	if key == "" {
		return c.data, nil
	}
	value, ok := c.data[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("key is not found: %s", key))
	}

	res, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("the type of value is invalid, key: %s", key))
	}
	return res, nil
}

func (c *JSONConfigContainer) OnChange(key string, fn func(value string)) {
	logs.Warn("unsupported operation")
}

// Bool returns the boolean value for a given key.
func (c *JSONConfigContainer) Bool(key string) (bool, error) {
	val := c.getData(key)
	if val != nil {
		return config.ParseBool(val)
	}
	return false, fmt.Errorf("not exist key: %q", key)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultval
func (c *JSONConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	if v, err := c.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

// Int returns the integer value for a given key.
func (c *JSONConfigContainer) Int(key string) (int, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		} else if v, ok := val.(string); ok {
			return strconv.Atoi(v)
		}
		return 0, errors.New("not valid value")
	}
	return 0, errors.New("not exist key:" + key)
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaultval
func (c *JSONConfigContainer) DefaultInt(key string, defaultVal int) int {
	if v, err := c.Int(key); err == nil {
		return v
	}
	return defaultVal
}

// Int64 returns the int64 value for a given key.
func (c *JSONConfigContainer) Int64(key string) (int64, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		}
		return 0, errors.New("not int64 value")
	}
	return 0, errors.New("not exist key:" + key)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaultval
func (c *JSONConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := c.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

// Float returns the float value for a given key.
func (c *JSONConfigContainer) Float(key string) (float64, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		}
		return 0.0, errors.New("not float64 value")
	}
	return 0.0, errors.New("not exist key:" + key)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaultval
func (c *JSONConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := c.Float(key); err == nil {
		return v
	}
	return defaultVal
}

// String returns the string value for a given key.
func (c *JSONConfigContainer) String(key string) (string, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v, nil
		}
	}
	return "", nil
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultval
func (c *JSONConfigContainer) DefaultString(key string, defaultVal string) string {
	// TODO FIXME should not use "" to replace non existence
	if v, err := c.String(key); v != "" && err == nil {
		return v
	}
	return defaultVal
}

// Strings returns the []string value for a given key.
func (c *JSONConfigContainer) Strings(key string) ([]string, error) {
	stringVal, err := c.String(key)
	if stringVal == "" || err != nil {
		return nil, err
	}
	return strings.Split(stringVal, ";"), nil
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultval
func (c *JSONConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	if v, err := c.Strings(key); v != nil && err == nil {
		return v
	}
	return defaultVal
}

// GetSection returns map for the given section
func (c *JSONConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		return v.(map[string]string), nil
	}
	return nil, errors.New("nonexist section " + section)
}

// SaveConfigFile save the config into file
func (c *JSONConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

// Set writes a new value for key.
func (c *JSONConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *JSONConfigContainer) DIY(key string) (v interface{}, err error) {
	val := c.getData(key)
	if val != nil {
		return val, nil
	}
	return nil, errors.New("not exist key")
}

// section.key or key
func (c *JSONConfigContainer) getData(key string) interface{} {
	if len(key) == 0 {
		return nil
	}

	c.RLock()
	defer c.RUnlock()

	sectionKeys := strings.Split(key, "::")
	if len(sectionKeys) >= 2 {
		curValue, ok := c.data[sectionKeys[0]]
		if !ok {
			return nil
		}
		for _, key := range sectionKeys[1:] {
			if v, ok := curValue.(map[string]interface{}); ok {
				if curValue, ok = v[key]; !ok {
					return nil
				}
			}
		}
		return curValue
	}
	if v, ok := c.data[key]; ok {
		return v
	}
	return nil
}

func init() {
	config.Register("json", &JSONConfig{})
}
