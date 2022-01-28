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
	"net/url"
	"strings"
)

// FlashData is a tools to maintain data when using across request.
type FlashData struct {
	Data map[string]string
}

// NewFlash return a new empty FlashData struct.
func NewFlash() *FlashData {
	return &FlashData{
		Data: make(map[string]string),
	}
}

// Set message to flash
func (fd *FlashData) Set(key string, msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data[key] = msg
	} else {
		fd.Data[key] = fmt.Sprintf(msg, args...)
	}
}

// Success writes success message to flash.
func (fd *FlashData) Success(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["success"] = msg
	} else {
		fd.Data["success"] = fmt.Sprintf(msg, args...)
	}
}

// Notice writes notice message to flash.
func (fd *FlashData) Notice(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["notice"] = msg
	} else {
		fd.Data["notice"] = fmt.Sprintf(msg, args...)
	}
}

// Warning writes warning message to flash.
func (fd *FlashData) Warning(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["warning"] = msg
	} else {
		fd.Data["warning"] = fmt.Sprintf(msg, args...)
	}
}

// Error writes error message to flash.
func (fd *FlashData) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["error"] = msg
	} else {
		fd.Data["error"] = fmt.Sprintf(msg, args...)
	}
}

// Store does the saving operation of flash data.
// the data are encoded and saved in cookie.
func (fd *FlashData) Store(c *Controller) {
	c.Data["flash"] = fd.Data
	var flashValue string
	for key, value := range fd.Data {
		flashValue += "\x00" + key + "\x23" + BasConfig.WebConfig.FlashSeparator + "\x23" + value + "\x00"
	}
	c.Ctx.SetCookie(BasConfig.WebConfig.FlashName, url.QueryEscape(flashValue), 0, "/")
}

// ReadFromRequest parsed flash data from encoded values in cookie.
func ReadFromRequest(c *Controller) *FlashData {
	flash := NewFlash()
	if cookie, err := c.Ctx.Request.Cookie(BasConfig.WebConfig.FlashName); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, "\x23"+BasConfig.WebConfig.FlashSeparator+"\x23")
				if len(kv) == 2 {
					flash.Data[kv[0]] = kv[1]
				}
			}
		}
		//read one time then delete it
		c.Ctx.SetCookie(BasConfig.WebConfig.FlashName, "", -1, "/")
	}
	c.Data["flash"] = flash.Data
	return flash
}
