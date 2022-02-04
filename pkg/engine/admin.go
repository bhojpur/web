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
	"fmt"
	"net/http"
	"reflect"
	"time"

	logsvr "github.com/bhojpur/logger/pkg/engine"
)

// WebAdminApp is the default adminApp used by admin module.
var webAdminApp *adminApp

// FilterMonitorFunc is default monitor filter when admin module is enable.
// if this func returns, admin module records qps for this request by condition of this function logic.
// usage:
// 	func MyFilterMonitor(method, requestPath string, t time.Duration, pattern string, statusCode int) bool {
//	 	if method == "POST" {
//			return false
//	 	}
//	 	if t.Nanoseconds() < 100 {
//			return false
//	 	}
//	 	if strings.HasPrefix(requestPath, "/bhojpur") {
//			return false
//	 	}
//	 	return true
// 	}
// 	websvr.FilterMonitorFunc = MyFilterMonitor.
var FilterMonitorFunc func(string, string, time.Duration, string, int) bool

func init() {
	FilterMonitorFunc = func(string, string, time.Duration, string, int) bool { return true }
}

func list(root string, p interface{}, m M) {
	pt := reflect.TypeOf(p)
	pv := reflect.ValueOf(p)
	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
		pv = pv.Elem()
	}
	for i := 0; i < pv.NumField(); i++ {
		var key string
		if root == "" {
			key = pt.Field(i).Name
		} else {
			key = root + "." + pt.Field(i).Name
		}
		if pv.Field(i).Kind() == reflect.Struct {
			list(key, pv.Field(i).Interface(), m)
		} else {
			m[key] = pv.Field(i).Interface()
		}
	}
}

func writeJSON(rw http.ResponseWriter, jsonData []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonData)
}

// adminApp is an http.HandlerFunc map used as webAdminApp.
type adminApp struct {
	*HttpServer
}

// Run start Bhojpur Web admin
func (admin *adminApp) Run() {
	logsvr.Debug("now we don't start tasks here, if you use task module," +
		" please invoke task.StartTask, or task will not be executed")
	addr := BConfig.Listen.AdminAddr
	if BConfig.Listen.AdminPort != 0 {
		addr = fmt.Sprintf("%s:%d", BConfig.Listen.AdminAddr, BConfig.Listen.AdminPort)
	}
	logsvr.Info("Bhojpur Web - Administration server running on %s", addr)
	admin.HttpServer.Run(addr)
}

func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {

		c := &adminController{
			servers: make([]*HttpServer, 0, 2),
		}

		// copy config to avoid conflict
		adminCfg := *BConfig
		webAdminApp = &adminApp{
			HttpServer: NewHttpServerWithCfg(&adminCfg),
		}
		// keep in mind that all data should be html escaped to avoid XSS attack
		webAdminApp.Router("/", c, "get:AdminIndex")
		webAdminApp.Router("/qps", c, "get:QpsIndex")
		webAdminApp.Router("/prof", c, "get:ProfIndex")
		webAdminApp.Router("/healthcheck", c, "get:Healthcheck")
		webAdminApp.Router("/task", c, "get:TaskStatus")
		webAdminApp.Router("/listconf", c, "get:ListConf")
		webAdminApp.Router("/metrics", c, "get:PrometheusMetrics")

		go webAdminApp.Run()
	}
	return nil
}
