package etcd

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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdConfigureProvider_Parse(t *testing.T) {
	provider := &EtcdConfigureProvider{}
	cfger, err := provider.Parse(readEtcdConfig())
	assert.Nil(t, err)
	assert.NotNil(t, cfger)
}

func TestEtcdConfigure(t *testing.T) {

	provider := &EtcdConfigureProvider{}
	cfger, _ := provider.Parse(readEtcdConfig())

	subCfger, err := cfger.Sub("sub.")
	assert.Nil(t, err)
	assert.NotNil(t, subCfger)

	subSubCfger, err := subCfger.Sub("sub.")
	assert.NotNil(t, subSubCfger)
	assert.Nil(t, err)

	str, err := subSubCfger.String("key1")
	assert.Nil(t, err)
	assert.Equal(t, "sub.sub.key", str)

	// we cannot test it
	subSubCfger.OnChange("watch", func(value string) {
		// do nothing
	})

	defStr := cfger.DefaultString("not_exit", "default value")
	assert.Equal(t, "default value", defStr)

	defInt64 := cfger.DefaultInt64("not_exit", -1)
	assert.Equal(t, int64(-1), defInt64)

	defInt := cfger.DefaultInt("not_exit", -2)
	assert.Equal(t, -2, defInt)

	defFlt := cfger.DefaultFloat("not_exit", 12.3)
	assert.Equal(t, 12.3, defFlt)

	defBl := cfger.DefaultBool("not_exit", true)
	assert.True(t, defBl)

	defStrs := cfger.DefaultStrings("not_exit", []string{"hello"})
	assert.Equal(t, []string{"hello"}, defStrs)

	fl, err := cfger.Float("current.float")
	assert.Nil(t, err)
	assert.Equal(t, 1.23, fl)

	bl, err := cfger.Bool("current.bool")
	assert.Nil(t, err)
	assert.True(t, bl)

	it, err := cfger.Int("current.int")
	assert.Nil(t, err)
	assert.Equal(t, 11, it)

	str, err = cfger.String("current.string")
	assert.Nil(t, err)
	assert.Equal(t, "hello", str)

	tn := &TestEntity{}
	err = cfger.Unmarshaler("current.serialize.", tn)
	assert.Nil(t, err)
	assert.Equal(t, "test", tn.Name)
}

type TestEntity struct {
	Name string    `yaml:"name"`
	Sub  SubEntity `yaml:"sub"`
}

type SubEntity struct {
	SubName string `yaml:"subName"`
}

func readEtcdConfig() string {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		addr = "localhost:2379"
	}

	obj := clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 3 * time.Second,
	}
	cfg, _ := json.Marshal(obj)
	return string(cfg)
}
