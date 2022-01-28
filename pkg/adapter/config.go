package adapter

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
	"github.com/bhojpur/web/pkg/adapter/session"
	newCfg "github.com/bhojpur/web/pkg/core/config"
	web "github.com/bhojpur/web/pkg/engine"
)

// Config is the main struct for BConfig
type Config web.Config

// Listen holds for http and https related config
type Listen web.Listen

// WebConfig holds web related config
type WebConfig web.WebConfig

// SessionConfig holds session related config
type SessionConfig web.SessionConfig

// LogConfig holds Log related config
type LogConfig web.LogConfig

var (
	// BConfig is the default config for Application
	BConfig *Config
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *bhojpurAppConfig
	// AppPath is the absolute path to the app
	AppPath string
	// GlobalSessions is the instance for the session manager
	GlobalSessions *session.Manager

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
	// WorkPath is the absolute path to project root directory
	WorkPath string
)

func init() {
	BConfig = (*Config)(web.BConfig)
	AppPath = web.AppPath

	WorkPath = web.WorkPath

	AppConfig = &bhojpurAppConfig{innerConfig: (newCfg.Configure)(web.AppConfig)}
}

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	return web.LoadAppConfig(adapterName, configPath)
}

type bhojpurAppConfig struct {
	innerConfig newCfg.Configure
}

func (b *bhojpurAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val); err != nil {
		return b.innerConfig.Set(key, val)
	}
	return nil
}

func (b *bhojpurAppConfig) String(key string) string {
	if v, err := b.innerConfig.String(BConfig.RunMode + "::" + key); v != "" && err != nil {
		return v
	}
	res, _ := b.innerConfig.String(key)
	return res
}

func (b *bhojpurAppConfig) Strings(key string) []string {
	if v, err := b.innerConfig.Strings(BConfig.RunMode + "::" + key); len(v) > 0 && err != nil {
		return v
	}
	res, _ := b.innerConfig.Strings(key)
	return res
}

func (b *bhojpurAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *bhojpurAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *bhojpurAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *bhojpurAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *bhojpurAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *bhojpurAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *bhojpurAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *bhojpurAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
