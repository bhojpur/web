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

// It implements generation and verification of image CAPTCHAs.
// An example for using the captcha library is shown below
//
// ```
// package controllers
//
// import (
// 	"github.com/bhojpur/image"
// 	"github.com/bhojpur/cache"
// 	"github.com/bhojpur/web/pkg/core/utils/captcha"
// )
//
// var cpt *captcha.Captcha
//
// func init() {
// 	// use web cache system store the captcha data
// 	store := cache.NewMemoryCache()
// 	cpt = captcha.NewWithFilter("/captcha/", store)
// }
//
// type MainController struct {
// 	webapp.Controller
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
	context2 "context"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	imageutl "github.com/bhojpur/image/pkg/synthesis"
	logs "github.com/bhojpur/logger/pkg/engine"
	token "github.com/bhojpur/token/pkg/utils"
	"github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/engine"
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
type Captcha struct {
	// bhojpur cache store
	store Storage

	// url prefix for captcha image
	URLPrefix string

	// specify captcha id input field name
	FieldIDName string
	// specify captcha result input field name
	FieldCaptchaName string

	// captcha image width and height
	StdWidth  int
	StdHeight int

	// captcha chars nums
	ChallengeNums int

	// captcha expiration seconds
	Expiration time.Duration

	// cache key prefix
	CachePrefix string
}

// generate key string
func (c *Captcha) key(id string) string {
	return c.CachePrefix + id
}

// generate rand chars with default chars
func (c *Captcha) genRandChars() []byte {
	return token.RandomCreateBytes(c.ChallengeNums, defaultChars...)
}

// Handler bhojpur filter handler for serve captcha image
func (c *Captcha) Handler(ctx *context.Context) {
	var chars []byte

	id := path.Base(ctx.Request.RequestURI)
	if i := strings.Index(id, "."); i != -1 {
		id = id[:i]
	}

	key := c.key(id)

	if len(ctx.Input.Query("reload")) > 0 {
		chars = c.genRandChars()
		if err := c.store.Put(context2.Background(), key, chars, c.Expiration); err != nil {
			ctx.Output.SetStatus(500)
			ctx.WriteString("captcha reload error")
			logs.Error("Reload Create Captcha Error:", err)
			return
		}
	} else {
		val, _ := c.store.Get(context2.Background(), key)
		if v, ok := val.([]byte); ok {
			chars = v
		} else {
			ctx.Output.SetStatus(404)
			ctx.WriteString("captcha not found")
			return
		}
	}

	img := imageutl.NewImage(chars, c.StdWidth, c.StdHeight)
	if _, err := img.WriteTo(ctx.ResponseWriter); err != nil {
		logs.Error("Write Captcha Image Error:", err)
	}
}

// CreateCaptchaHTML template func for output html
func (c *Captcha) CreateCaptchaHTML() template.HTML {
	value, err := c.CreateCaptcha()
	if err != nil {
		logs.Error("Create Captcha Error:", err)
		return ""
	}

	// create html
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`+
		`<a class="captcha" href="javascript:">`+
		`<img onclick="this.src=('%s%s.png?reload='+(new Date()).getTime())" class="captcha-img" src="%s%s.png">`+
		`</a>`, c.FieldIDName, value, c.URLPrefix, value, c.URLPrefix, value))
}

// CreateCaptcha create a new captcha id
func (c *Captcha) CreateCaptcha() (string, error) {
	// generate captcha id
	id := string(token.RandomCreateBytes(15))

	// get the captcha chars
	chars := c.genRandChars()

	// save to store
	if err := c.store.Put(context2.Background(), c.key(id), chars, c.Expiration); err != nil {
		return "", err
	}

	return id, nil
}

// VerifyReq verify from a request
func (c *Captcha) VerifyReq(req *http.Request) bool {
	req.ParseForm()
	return c.Verify(req.Form.Get(c.FieldIDName), req.Form.Get(c.FieldCaptchaName))
}

// Verify direct verify id and challenge string
func (c *Captcha) Verify(id string, challenge string) (success bool) {
	if len(challenge) == 0 || len(id) == 0 {
		return
	}

	var chars []byte

	key := c.key(id)
	val, _ := c.store.Get(context2.Background(), key)
	if v, ok := val.([]byte); ok {
		chars = v
	} else {
		return
	}

	defer func() {
		// finally remove it
		c.store.Delete(context2.Background(), key)
	}()

	if len(chars) != len(challenge) {
		return
	}
	// verify challenge
	for i, c := range chars {
		if c != challenge[i]-48 {
			return
		}
	}

	return true
}

// NewCaptcha create a new captcha.Captcha
func NewCaptcha(urlPrefix string, store Storage) *Captcha {
	cpt := &Captcha{}
	cpt.store = store
	cpt.FieldIDName = fieldIDName
	cpt.FieldCaptchaName = fieldCaptchaName
	cpt.ChallengeNums = challengeNums
	cpt.Expiration = expiration
	cpt.CachePrefix = cachePrefix
	cpt.StdWidth = stdWidth
	cpt.StdHeight = stdHeight

	if len(urlPrefix) == 0 {
		urlPrefix = defaultURLPrefix
	}

	if urlPrefix[len(urlPrefix)-1] != '/' {
		urlPrefix += "/"
	}

	cpt.URLPrefix = urlPrefix

	return cpt
}

// NewWithFilter create a new captcha.Captcha and auto AddFilter for serve captacha image
// and add a template func for output html
func NewWithFilter(urlPrefix string, store Storage) *Captcha {
	cpt := NewCaptcha(urlPrefix, store)

	// create filter for serve captcha image
	engine.InsertFilter(cpt.URLPrefix+"*", engine.BeforeRouter, cpt.Handler)

	// add to template func map
	engine.AddFuncMap("create_captcha", cpt.CreateCaptchaHTML)

	return cpt
}

type Storage interface {
	// Get a cached value by key.
	Get(ctx context2.Context, key string) (interface{}, error)
	// Set a cached value with key and expire time.
	Put(ctx context2.Context, key string, val interface{}, timeout time.Duration) error
	// Delete cached value by key.
	Delete(ctx context2.Context, key string) error
}
