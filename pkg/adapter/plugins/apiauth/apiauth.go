package apiauth

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

// Package apiauth provides handlers to enable apiauth support.
//
// Simple Usage:
//	import(
//		bhojpur "github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/plugins/apiauth"
//	)
//
//	func main(){
//		// apiauth every request
//		bhojpur.InsertFilter("*", bhojpur.BeforeRouter,apiauth.APIBaiscAuth("appid","appkey"))
//		bhojpur.Run()
//	}
//
// Advanced Usage:
//
//	func getAppSecret(appid string) string {
//		// get appsecret by appid
//		// maybe store in configure, maybe in database
//	}
//
//	bhojpur.InsertFilter("*", bhojpur.BeforeRouter,apiauth.APISecretAuth(getAppSecret, 360))
//
// Information:
//
// In the request user should include these params in the query
//
// 1. appid
//
//		 appid is assigned to the application
//
// 2. signature
//
//	get the signature use apiauth.Signature()
//
//	when you send to server remember use url.QueryEscape()
//
// 3. timestamp:
//
//       send the request time, the format is yyyy-mm-dd HH:ii:ss

import (
	"net/url"

	bhojpur "github.com/bhojpur/web/pkg/adapter"
	"github.com/bhojpur/web/pkg/adapter/context"
	beecontext "github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/filter/apiauth"
)

// AppIDToAppSecret is used to get appsecret throw appid
type AppIDToAppSecret apiauth.AppIDToAppSecret

// APIBasicAuth use the basic appid/appkey as the AppIdToAppSecret
func APIBasicAuth(appid, appkey string) bhojpur.FilterFunc {
	f := apiauth.APIBasicAuth(appid, appkey)
	return func(c *context.Context) {
		f((*beecontext.Context)(c))
	}
}

// APIBaiscAuth calls APIBasicAuth for previous callers
func APIBaiscAuth(appid, appkey string) bhojpur.FilterFunc {
	return APIBasicAuth(appid, appkey)
}

// APISecretAuth use AppIdToAppSecret verify and
func APISecretAuth(f AppIDToAppSecret, timeout int) bhojpur.FilterFunc {
	ft := apiauth.APISecretAuth(apiauth.AppIDToAppSecret(f), timeout)
	return func(ctx *context.Context) {
		ft((*beecontext.Context)(ctx))
	}
}

// Signature used to generate signature with the appsecret/method/params/RequestURI
func Signature(appsecret, method string, params url.Values, requestURL string) string {
	return apiauth.Signature(appsecret, method, params, requestURL)
}
