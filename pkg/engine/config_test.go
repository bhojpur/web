package engine

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
	"reflect"
	"testing"

	webJson "github.com/bhojpur/web/pkg/core/config/json"
)

func TestDefaults(t *testing.T) {
	if BConfig.WebConfig.FlashName != "BHOJPUR_FLASH" {
		t.Errorf("FlashName was not set to default.")
	}

	if BConfig.WebConfig.FlashSeparator != "BHOJPURFLASH" {
		t.Errorf("FlashName was not set to default.")
	}
}

func TestLoadAppConfig(t *testing.T) {
	println(1 << 30)
}

func TestAssignConfig_01(t *testing.T) {
	_BConfig := &Config{}
	_BConfig.AppName = "bhojpur_test"
	jcf := &webJson.JSONConfig{}
	ac, _ := jcf.ParseData([]byte(`{"AppName":"bhojpur_json"}`))
	assignSingleConfig(_BConfig, ac)
	if _BConfig.AppName != "bhojpur_json" {
		t.Log(_BConfig)
		t.FailNow()
	}
}

func TestAssignConfig_02(t *testing.T) {
	_BConfig := &Config{}
	bs, _ := json.Marshal(newBConfig())

	jsonMap := M{}
	json.Unmarshal(bs, &jsonMap)

	configMap := M{}
	for k, v := range jsonMap {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			for k1, v1 := range v.(M) {
				if reflect.TypeOf(v1).Kind() == reflect.Map {
					for k2, v2 := range v1.(M) {
						configMap[k2] = v2
					}
				} else {
					configMap[k1] = v1
				}
			}
		} else {
			configMap[k] = v
		}
	}
	configMap["MaxMemory"] = 1024
	configMap["Graceful"] = true
	configMap["XSRFExpire"] = 32
	configMap["SessionProviderConfig"] = "file"
	configMap["FileLineNum"] = true

	jcf := &webJson.JSONConfig{}
	bs, _ = json.Marshal(configMap)
	ac, _ := jcf.ParseData(bs)

	for _, i := range []interface{}{_BConfig, &_BConfig.Listen, &_BConfig.WebConfig, &_BConfig.Log, &_BConfig.WebConfig.Session} {
		assignSingleConfig(i, ac)
	}

	if _BConfig.MaxMemory != 1024 {
		t.Log(_BConfig.MaxMemory)
		t.FailNow()
	}

	if !_BConfig.Listen.Graceful {
		t.Log(_BConfig.Listen.Graceful)
		t.FailNow()
	}

	if _BConfig.WebConfig.XSRFExpire != 32 {
		t.Log(_BConfig.WebConfig.XSRFExpire)
		t.FailNow()
	}

	if _BConfig.WebConfig.Session.SessionProviderConfig != "file" {
		t.Log(_BConfig.WebConfig.Session.SessionProviderConfig)
		t.FailNow()
	}

	if !_BConfig.Log.FileLineNum {
		t.Log(_BConfig.Log.FileLineNum)
		t.FailNow()
	}
}

func TestAssignConfig_03(t *testing.T) {
	jcf := &webJson.JSONConfig{}
	ac, _ := jcf.ParseData([]byte(`{"AppName":"bhojpur"}`))
	ac.Set("AppName", "test_app")
	ac.Set("RunMode", "online")
	ac.Set("StaticDir", "download:down download2:down2")
	ac.Set("StaticExtensionsToGzip", ".css,.js,.html,.jpg,.png")
	ac.Set("StaticCacheFileSize", "87456")
	ac.Set("StaticCacheFileNum", "1254")
	assignConfig(ac)

	t.Logf("%#v", BConfig)

	if BConfig.AppName != "test_app" {
		t.FailNow()
	}

	if BConfig.RunMode != "online" {
		t.FailNow()
	}
	if BConfig.WebConfig.StaticDir["/download"] != "down" {
		t.FailNow()
	}
	if BConfig.WebConfig.StaticDir["/download2"] != "down2" {
		t.FailNow()
	}
	if BConfig.WebConfig.StaticCacheFileSize != 87456 {
		t.FailNow()
	}
	if BConfig.WebConfig.StaticCacheFileNum != 1254 {
		t.FailNow()
	}
	if len(BConfig.WebConfig.StaticExtensionsToGzip) != 5 {
		t.FailNow()
	}
}
