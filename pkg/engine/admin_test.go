package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/core/admin"
)

type SampleDatabaseCheck struct {
}

type SampleCacheCheck struct {
}

func (dc *SampleDatabaseCheck) Check() error {
	return nil
}

func (cc *SampleCacheCheck) Check() error {
	return errors.New("no cache detected")
}

func TestList_01(t *testing.T) {
	m := make(M)
	list("BasConfig", BasConfig, m)
	t.Log(m)
	om := oldMap()
	for k, v := range om {
		if fmt.Sprint(m[k]) != fmt.Sprint(v) {
			t.Log(k, "old-key", v, "new-key", m[k])
			t.FailNow()
		}
	}
}

func oldMap() M {
	m := make(M)
	m["BasConfig.AppName"] = BasConfig.AppName
	m["BasConfig.RunMode"] = BasConfig.RunMode
	m["BasConfig.RouterCaseSensitive"] = BasConfig.RouterCaseSensitive
	m["BasConfig.ServerName"] = BasConfig.ServerName
	m["BasConfig.RecoverPanic"] = BasConfig.RecoverPanic
	m["BasConfig.CopyRequestBody"] = BasConfig.CopyRequestBody
	m["BasConfig.EnableGzip"] = BasConfig.EnableGzip
	m["BasConfig.MaxMemory"] = BasConfig.MaxMemory
	m["BasConfig.EnableErrorsShow"] = BasConfig.EnableErrorsShow
	m["BasConfig.Listen.Graceful"] = BasConfig.Listen.Graceful
	m["BasConfig.Listen.ServerTimeOut"] = BasConfig.Listen.ServerTimeOut
	m["BasConfig.Listen.ListenTCP4"] = BasConfig.Listen.ListenTCP4
	m["BasConfig.Listen.EnableHTTP"] = BasConfig.Listen.EnableHTTP
	m["BasConfig.Listen.HTTPAddr"] = BasConfig.Listen.HTTPAddr
	m["BasConfig.Listen.HTTPPort"] = BasConfig.Listen.HTTPPort
	m["BasConfig.Listen.EnableHTTPS"] = BasConfig.Listen.EnableHTTPS
	m["BasConfig.Listen.HTTPSAddr"] = BasConfig.Listen.HTTPSAddr
	m["BasConfig.Listen.HTTPSPort"] = BasConfig.Listen.HTTPSPort
	m["BasConfig.Listen.HTTPSCertFile"] = BasConfig.Listen.HTTPSCertFile
	m["BasConfig.Listen.HTTPSKeyFile"] = BasConfig.Listen.HTTPSKeyFile
	m["BasConfig.Listen.EnableAdmin"] = BasConfig.Listen.EnableAdmin
	m["BasConfig.Listen.AdminAddr"] = BasConfig.Listen.AdminAddr
	m["BasConfig.Listen.AdminPort"] = BasConfig.Listen.AdminPort
	m["BasConfig.Listen.EnableFcgi"] = BasConfig.Listen.EnableFcgi
	m["BasConfig.Listen.EnableStdIo"] = BasConfig.Listen.EnableStdIo
	m["BasConfig.WebConfig.AutoRender"] = BasConfig.WebConfig.AutoRender
	m["BasConfig.WebConfig.EnableDocs"] = BasConfig.WebConfig.EnableDocs
	m["BasConfig.WebConfig.FlashName"] = BasConfig.WebConfig.FlashName
	m["BasConfig.WebConfig.FlashSeparator"] = BasConfig.WebConfig.FlashSeparator
	m["BasConfig.WebConfig.DirectoryIndex"] = BasConfig.WebConfig.DirectoryIndex
	m["BasConfig.WebConfig.StaticDir"] = BasConfig.WebConfig.StaticDir
	m["BasConfig.WebConfig.StaticExtensionsToGzip"] = BasConfig.WebConfig.StaticExtensionsToGzip
	m["BasConfig.WebConfig.StaticCacheFileSize"] = BasConfig.WebConfig.StaticCacheFileSize
	m["BasConfig.WebConfig.StaticCacheFileNum"] = BasConfig.WebConfig.StaticCacheFileNum
	m["BasConfig.WebConfig.TemplateLeft"] = BasConfig.WebConfig.TemplateLeft
	m["BasConfig.WebConfig.TemplateRight"] = BasConfig.WebConfig.TemplateRight
	m["BasConfig.WebConfig.ViewsPath"] = BasConfig.WebConfig.ViewsPath
	m["BasConfig.WebConfig.EnableXSRF"] = BasConfig.WebConfig.EnableXSRF
	m["BasConfig.WebConfig.XSRFExpire"] = BasConfig.WebConfig.XSRFExpire
	m["BasConfig.WebConfig.Session.SessionOn"] = BasConfig.WebConfig.Session.SessionOn
	m["BasConfig.WebConfig.Session.SessionProvider"] = BasConfig.WebConfig.Session.SessionProvider
	m["BasConfig.WebConfig.Session.SessionName"] = BasConfig.WebConfig.Session.SessionName
	m["BasConfig.WebConfig.Session.SessionGCMaxLifetime"] = BasConfig.WebConfig.Session.SessionGCMaxLifetime
	m["BasConfig.WebConfig.Session.SessionProviderConfig"] = BasConfig.WebConfig.Session.SessionProviderConfig
	m["BasConfig.WebConfig.Session.SessionCookieLifeTime"] = BasConfig.WebConfig.Session.SessionCookieLifeTime
	m["BasConfig.WebConfig.Session.SessionAutoSetCookie"] = BasConfig.WebConfig.Session.SessionAutoSetCookie
	m["BasConfig.WebConfig.Session.SessionDomain"] = BasConfig.WebConfig.Session.SessionDomain
	m["BasConfig.WebConfig.Session.SessionDisableHTTPOnly"] = BasConfig.WebConfig.Session.SessionDisableHTTPOnly
	m["BasConfig.Log.AccessLogs"] = BasConfig.Log.AccessLogs
	m["BasConfig.Log.EnableStaticLogs"] = BasConfig.Log.EnableStaticLogs
	m["BasConfig.Log.AccessLogsFormat"] = BasConfig.Log.AccessLogsFormat
	m["BasConfig.Log.FileLineNum"] = BasConfig.Log.FileLineNum
	m["BasConfig.Log.Outputs"] = BasConfig.Log.Outputs
	return m
}

func TestWriteJSON(t *testing.T) {
	t.Log("Testing the adding of JSON to the response")

	w := httptest.NewRecorder()
	originalBody := []int{1, 2, 3}

	res, _ := json.Marshal(originalBody)

	writeJSON(w, res)

	decodedBody := []int{}
	err := json.NewDecoder(w.Body).Decode(&decodedBody)

	if err != nil {
		t.Fatal("Could not decode response body into slice.")
	}

	for i := range decodedBody {
		if decodedBody[i] != originalBody[i] {
			t.Fatalf("Expected %d but got %d in decoded body slice", originalBody[i], decodedBody[i])
		}
	}
}

func TestHealthCheckHandlerDefault(t *testing.T) {
	endpointPath := "/healthcheck"

	admin.AddHealthCheck("database", &SampleDatabaseCheck{})
	admin.AddHealthCheck("cache", &SampleCacheCheck{})

	req, err := http.NewRequest("GET", endpointPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(heathCheck)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "database") {
		t.Errorf("Expected 'database' in generated template.")
	}

}

func TestBuildHealthCheckResponseList(t *testing.T) {
	healthCheckResults := [][]string{
		[]string{
			"error",
			"Database",
			"Error occured whie starting the db",
		},
		[]string{
			"success",
			"Cache",
			"Cache started successfully",
		},
	}

	responseList := buildHealthCheckResponseList(&healthCheckResults)

	if len(responseList) != len(healthCheckResults) {
		t.Errorf("invalid response map length: got %d want %d",
			len(responseList), len(healthCheckResults))
	}

	responseFields := []string{"name", "message", "status"}

	for _, response := range responseList {
		for _, field := range responseFields {
			_, ok := response[field]
			if !ok {
				t.Errorf("expected %s to be in the response %v", field, response)
			}
		}

	}

}

func TestHealthCheckHandlerReturnsJSON(t *testing.T) {

	admin.AddHealthCheck("database", &SampleDatabaseCheck{})
	admin.AddHealthCheck("cache", &SampleCacheCheck{})

	req, err := http.NewRequest("GET", "/healthcheck?json=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(heathCheck)

	handler.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	decodedResponseBody := []map[string]interface{}{}
	expectedResponseBody := []map[string]interface{}{}

	expectedJSONString := []byte(`
		[
			{
				"message":"database",
				"name":"success",
				"status":"OK"
			},
			{
				"message":"cache",
				"name":"error",
				"status":"no cache detected"
			}
		]
	`)

	json.Unmarshal(expectedJSONString, &expectedResponseBody)

	json.Unmarshal(w.Body.Bytes(), &decodedResponseBody)

	if len(expectedResponseBody) != len(decodedResponseBody) {
		t.Errorf("invalid response map length: got %d want %d",
			len(decodedResponseBody), len(expectedResponseBody))
	}
	assert.Equal(t, len(expectedResponseBody), len(decodedResponseBody))
	assert.Equal(t, 2, len(decodedResponseBody))

	var database, cache map[string]interface{}
	if decodedResponseBody[0]["message"] == "database" {
		database = decodedResponseBody[0]
		cache = decodedResponseBody[1]
	} else {
		database = decodedResponseBody[1]
		cache = decodedResponseBody[0]
	}

	assert.Equal(t, expectedResponseBody[0], database)
	assert.Equal(t, expectedResponseBody[1], cache)

}
