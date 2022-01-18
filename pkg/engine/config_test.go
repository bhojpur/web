package engine

import (
	"encoding/json"
	"reflect"
	"testing"

	beeJson "github.com/bhojpur/web/pkg/core/config/json"
)

func TestDefaults(t *testing.T) {
	if BasConfig.WebConfig.FlashName != "BHOJPUR_FLASH" {
		t.Errorf("FlashName was not set to default.")
	}

	if BasConfig.WebConfig.FlashSeparator != "BHOJPURFLASH" {
		t.Errorf("FlashName was not set to default.")
	}
}

func TestAssignConfig_01(t *testing.T) {
	_BasConfig := &Config{}
	_BasConfig.AppName = "bhojpur_test"
	jcf := &beeJson.JSONConfig{}
	ac, _ := jcf.ParseData([]byte(`{"AppName":"bhojpur_json"}`))
	assignSingleConfig(_BasConfig, ac)
	if _BasConfig.AppName != "bhojpur_json" {
		t.Log(_BasConfig)
		t.FailNow()
	}
}

func TestAssignConfig_02(t *testing.T) {
	_BasConfig := &Config{}
	bs, _ := json.Marshal(newBasConfig())

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

	jcf := &beeJson.JSONConfig{}
	bs, _ = json.Marshal(configMap)
	ac, _ := jcf.ParseData(bs)

	for _, i := range []interface{}{_BasConfig, &_BasConfig.Listen, &_BasConfig.WebConfig, &_BasConfig.Log, &_BasConfig.WeBasConfig.Session} {
		assignSingleConfig(i, ac)
	}

	if _BasConfig.MaxMemory != 1024 {
		t.Log(_BasConfig.MaxMemory)
		t.FailNow()
	}

	if !_BasConfig.Listen.Graceful {
		t.Log(_BasConfig.Listen.Graceful)
		t.FailNow()
	}

	if _BasConfig.WebConfig.XSRFExpire != 32 {
		t.Log(_BasConfig.WebConfig.XSRFExpire)
		t.FailNow()
	}

	if _BasConfig.WebConfig.Session.SessionProviderConfig != "file" {
		t.Log(_BasConfig.WebConfig.Session.SessionProviderConfig)
		t.FailNow()
	}

	if !_BasConfig.Log.FileLineNum {
		t.Log(_BasConfig.Log.FileLineNum)
		t.FailNow()
	}

}

func TestAssignConfig_03(t *testing.T) {
	jcf := &beeJson.JSONConfig{}
	ac, _ := jcf.ParseData([]byte(`{"AppName":"bhojpur"}`))
	ac.Set("AppName", "test_app")
	ac.Set("RunMode", "online")
	ac.Set("StaticDir", "download:down download2:down2")
	ac.Set("StaticExtensionsToGzip", ".css,.js,.html,.jpg,.png")
	ac.Set("StaticCacheFileSize", "87456")
	ac.Set("StaticCacheFileNum", "1254")
	assignConfig(ac)

	t.Logf("%#v", BasConfig)

	if BasConfig.AppName != "test_app" {
		t.FailNow()
	}

	if BasConfig.RunMode != "online" {
		t.FailNow()
	}
	if BasConfig.WebConfig.StaticDir["/download"] != "down" {
		t.FailNow()
	}
	if BasConfig.WebConfig.StaticDir["/download2"] != "down2" {
		t.FailNow()
	}
	if BasConfig.WebConfig.StaticCacheFileSize != 87456 {
		t.FailNow()
	}
	if BasConfig.WebConfig.StaticCacheFileNum != 1254 {
		t.FailNow()
	}
	if len(BasConfig.WebConfig.StaticExtensionsToGzip) != 5 {
		t.FailNow()
	}
}
