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

// It provides handlers to enable apiauth support.
//
// Simple Usage:
//	import(
//		"github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/filter/apiauth"
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
//		 appid is assigned to the Web Application
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
//

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"time"

	ctxsvr "github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// AppIDToAppSecret gets appsecret through appid
type AppIDToAppSecret func(string) string

// APIBasicAuth uses the basic appid/appkey as the AppIdToAppSecret
func APIBasicAuth(appid, appkey string) websvr.FilterFunc {
	ft := func(aid string) string {
		if aid == appid {
			return appkey
		}
		return ""
	}
	return APISecretAuth(ft, 300)
}

// APISecretAuth uses AppIdToAppSecret verify and
func APISecretAuth(f AppIDToAppSecret, timeout int) websvr.FilterFunc {
	return func(ctx *ctxsvr.Context) {
		if ctx.Input.Query("appid") == "" {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("missing query parameter: appid")
			return
		}
		appsecret := f(ctx.Input.Query("appid"))
		if appsecret == "" {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("appid query parameter missing")
			return
		}
		if ctx.Input.Query("signature") == "" {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("missing query parameter: signature")

			return
		}
		if ctx.Input.Query("timestamp") == "" {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("missing query parameter: timestamp")
			return
		}
		u, err := time.Parse("2018-03-26 15:04:05", ctx.Input.Query("timestamp"))
		if err != nil {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("incorrect timestamp format. Should be in the form 2018-03-26 15:04:05")

			return
		}
		t := time.Now()
		if t.Sub(u).Seconds() > float64(timeout) {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("request timer timeout exceeded. Please try again")
			return
		}
		if ctx.Input.Query("signature") !=
			Signature(appsecret, ctx.Input.Method(), ctx.Request.Form, ctx.Input.URL()) {
			ctx.ResponseWriter.WriteHeader(403)
			ctx.WriteString("authentication failed")
		}
	}
}

// Signature generates signature with appsecret/method/params/RequestURI
func Signature(appsecret, method string, params url.Values, RequestURL string) (result string) {
	var b bytes.Buffer
	keys := make([]string, len(params))
	pa := make(map[string]string)
	for k, v := range params {
		pa[k] = v[0]
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		if key == "signature" {
			continue
		}

		val := pa[key]
		if key != "" && val != "" {
			b.WriteString(key)
			b.WriteString(val)
		}
	}

	stringToSign := fmt.Sprintf("%v\n%v\n%v\n", method, b.String(), RequestURL)

	sha256 := sha256.New
	hash := hmac.New(sha256, []byte(appsecret))
	hash.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
