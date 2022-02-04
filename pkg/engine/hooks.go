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
	"mime"
	"net/http"
	"path/filepath"

	logsvr "github.com/bhojpur/logger/pkg/engine"
	session "github.com/bhojpur/session/pkg/engine"
	"github.com/bhojpur/web/pkg/context"
)

// register MIME type with content type
func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

// register default error http handlers, 404,401,403,500 and 503.
func registerDefaultErrorHandler() error {
	m := map[string]func(http.ResponseWriter, *http.Request){
		"401": unauthorized,
		"402": paymentRequired,
		"403": forbidden,
		"404": notFound,
		"405": methodNotAllowed,
		"500": internalServerError,
		"501": notImplemented,
		"502": badGateway,
		"503": serviceUnavailable,
		"504": gatewayTimeout,
		"417": invalidxsrf,
		"422": missingxsrf,
		"413": payloadTooLarge,
	}
	for e, h := range m {
		if _, ok := ErrorMaps[e]; !ok {
			ErrorHandler(e, h)
		}
	}
	return nil
}

func registerSession() error {
	if BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig, err := AppConfig.String("sessionConfig")
		conf := new(session.ManagerConfig)
		if sessionConfig == "" || err != nil {
			conf.CookieName = BConfig.WebConfig.Session.SessionName
			conf.EnableSetCookie = BConfig.WebConfig.Session.SessionAutoSetCookie
			conf.Gclifetime = BConfig.WebConfig.Session.SessionGCMaxLifetime
			conf.Secure = BConfig.Listen.EnableHTTPS
			conf.CookieLifeTime = BConfig.WebConfig.Session.SessionCookieLifeTime
			conf.ProviderConfig = filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig)
			conf.DisableHTTPOnly = BConfig.WebConfig.Session.SessionDisableHTTPOnly
			conf.Domain = BConfig.WebConfig.Session.SessionDomain
			conf.EnableSidInHTTPHeader = BConfig.WebConfig.Session.SessionEnableSidInHTTPHeader
			conf.SessionNameInHTTPHeader = BConfig.WebConfig.Session.SessionNameInHTTPHeader
			conf.EnableSidInURLQuery = BConfig.WebConfig.Session.SessionEnableSidInURLQuery
			conf.CookieSameSite = BConfig.WebConfig.Session.SessionCookieSameSite
			conf.SessionIDPrefix = BConfig.WebConfig.Session.SessionIDPrefix
		} else {
			if err = json.Unmarshal([]byte(sessionConfig), conf); err != nil {
				return err
			}
		}
		if GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, conf); err != nil {
			return err
		}
		go GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	defer lockViewPaths()
	if err := AddViewPath(BConfig.WebConfig.ViewsPath); err != nil {
		if BConfig.RunMode == DEV {
			logsvr.Warn(err)
		}
		return err
	}
	return nil
}

func registerGzip() error {
	if BConfig.EnableGzip {
		context.InitGzip(
			AppConfig.DefaultInt("gzipMinLength", -1),
			AppConfig.DefaultInt("gzipCompressLevel", -1),
			AppConfig.DefaultStrings("includedMethods", []string{"GET"}),
		)
	}
	return nil
}
