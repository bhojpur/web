package yaml

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

// depend on github.com/bhojpur/web/pkg/core/yaml2
//
// go install github.com/bhojpur/web/pkg/core/yaml2
//
// Usage:
//  import(
//   _ "github.com/bhojpur/web/pkg/core/config/yaml"
//     "github.com/bhojpur/web/pkg/core/config"
//  )
//
//  cnf, err := config.NewConfig("yaml", "config.yaml")

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	goyaml2 "github.com/bhojpur/web/pkg/core/yaml2"
	"gopkg.in/yaml.v2"

	logs "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/core/config"
)

// Config is a yaml config parser and implements Config interface.
type Config struct{}

// Parse returns a ConfigContainer with parsed yaml config map.
func (yaml *Config) Parse(filename string) (y config.Configure, err error) {
	cnf, err := ReadYmlReader(filename)
	if err != nil {
		return
	}
	y = &ConfigContainer{
		data: cnf,
	}
	return
}

// ParseData parse yaml data
func (yaml *Config) ParseData(data []byte) (config.Configure, error) {
	cnf, err := parseYML(data)
	if err != nil {
		return nil, err
	}

	return &ConfigContainer{
		data: cnf,
	}, nil
}

// ReadYmlReader Read yaml file to map.
// if json like, use json package, unless goyaml2 package.
func ReadYmlReader(path string) (cnf map[string]interface{}, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return parseYML(buf)
}

// parseYML parse yaml formatted []byte to map.
func parseYML(buf []byte) (cnf map[string]interface{}, err error) {
	if len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" {
		log.Println("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			log.Println("It is Json Map")
			return
		}
	}

	data, err := goyaml2.Read(bytes.NewReader(buf))
	if err != nil {
		log.Println("Goyaml2 ERR>", string(buf), err)
		return
	}

	if data == nil {
		log.Println("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Not a Map? >> ", string(buf), data)
		cnf = nil
	}
	cnf = config.ExpandValueEnvForMap(cnf)
	return
}

// ConfigContainer is a config which represents the yaml configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

// Unmarshaler is similar to Sub
func (c *ConfigContainer) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	sub, err := c.sub(prefix)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(sub)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, obj)
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
	tmpData := c.data
	keys := strings.Split(key, ".")
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch v.(type) {
			case map[string]interface{}:
				{
					tmpData = v.(map[string]interface{})
					if idx == len(keys)-1 {
						return tmpData, nil
					}
				}
			default:
				return nil, errors.New(fmt.Sprintf("the key is invalid: %s", key))
			}
		}
	}

	return tmpData, nil
}

func (c *ConfigContainer) OnChange(key string, fn func(value string)) {
	// do nothing
	logs.Warn("Unsupported operation: OnChange")
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {
	v, err := c.getData(key)
	if err != nil {
		return false, err
	}
	return config.ParseBool(v)
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
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}
	return 0, errors.New("not int value")
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
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int64); ok {
		return vv, nil
	}
	return 0, errors.New("not bool value")
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
	if v, err := c.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
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
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv, nil
		}
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

	if v, ok := c.data[section]; ok {
		return v.(map[string]string), nil
	}
	return nil, errors.New("not exist section")
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = goyaml2.Write(f, c.data)
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
	return c.getData(key)
}

func (c *ConfigContainer) getData(key string) (interface{}, error) {

	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(c.key(key), ".")
	tmpData := c.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch v.(type) {
			case map[string]interface{}:
				{
					tmpData = v.(map[string]interface{})
					if idx == len(keys)-1 {
						return tmpData, nil
					}
				}
			default:
				{
					return v, nil
				}

			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func (c *ConfigContainer) key(key string) string {
	return key
}

func init() {
	config.Register("yaml", &Config{})
}
