package captcha

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

// Captcha implements generation and verification of image CAPTCHAs.
// an example for use captcha
//
// ```
// package controllers
//
// import (
// 	websvr "github.com/bhojpur/web/pkg/engine"
// 	"github.com/bhojpur/web/pkg/cache"
// 	"github.com/bhojpur/web/pkg/utils/captcha"
// )
//
// var cpt *captcha.Captcha
//
// func init() {
// 	// use bhojpur cache system store the captcha data
// 	store := cache.NewMemoryCache()
// 	cpt = captcha.NewWithFilter("/captcha/", store)
// }
//
// type MainController struct {
// 	websvr.Controller
// }
//
// func (this *MainController) Get() {
// 	this.TplName = "index.tpl"
// }
//
// func (this *MainController) Post() {
// 	this.TplName = "index.tpl"
//
// 	this.Data["Success"] = cpt.VerifyReq(this.Ctx.Request)
// }
// ```
//
// template usage
//
// ```
// {{.Success}}
// <form action="/" method="post">
// 	{{create_captcha}}
// 	<input name="captcha" type="text">
// </form>
// ```

import (
	"html/template"
	"net/http"
	"time"

	"github.com/bhojpur/web/pkg/captcha"
	ctxsvr "github.com/bhojpur/web/pkg/context"

	"github.com/bhojpur/web/pkg/adapter/cache"
	"github.com/bhojpur/web/pkg/adapter/context"
)

var (
	defaultChars = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

const (
	// default captcha attributes
	challengeNums    = 6
	expiration       = 600 * time.Second
	fieldIDName      = "captcha_id"
	fieldCaptchaName = "captcha"
	cachePrefix      = "captcha_"
	defaultURLPrefix = "/captcha/"
)

// Captcha struct
type Captcha captcha.Captcha

// Handler bhojpur filter handler for serve captcha image
func (c *Captcha) Handler(ctx *context.Context) {
	(*captcha.Captcha)(c).Handler((*ctxsvr.Context)(ctx))
}

// CreateCaptchaHTML template func for output html
func (c *Captcha) CreateCaptchaHTML() template.HTML {
	return (*captcha.Captcha)(c).CreateCaptchaHTML()
}

// CreateCaptcha create a new captcha id
func (c *Captcha) CreateCaptcha() (string, error) {
	return (*captcha.Captcha)(c).CreateCaptcha()
}

// VerifyReq verify from a request
func (c *Captcha) VerifyReq(req *http.Request) bool {
	return (*captcha.Captcha)(c).VerifyReq(req)
}

// Verify direct verify id and challenge string
func (c *Captcha) Verify(id string, challenge string) (success bool) {
	return (*captcha.Captcha)(c).Verify(id, challenge)
}

// NewCaptcha create a new captcha.Captcha
func NewCaptcha(urlPrefix string, store cache.Cache) *Captcha {
	return (*Captcha)(captcha.NewCaptcha(urlPrefix, cache.CreateOldToNewAdapter(store)))
}

// NewWithFilter create a new captcha.Captcha and auto AddFilter for serve captacha image
// and add a template func for output html
func NewWithFilter(urlPrefix string, store cache.Cache) *Captcha {
	return (*Captcha)(captcha.NewWithFilter(urlPrefix, cache.CreateOldToNewAdapter(store)))
}
