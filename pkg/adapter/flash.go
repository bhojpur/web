package adapter

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
	web "github.com/bhojpur/web/pkg/engine"
)

// FlashData is a tools to maintain data when using across request.
type FlashData web.FlashData

// NewFlash return a new empty FlashData struct.
func NewFlash() *FlashData {
	return (*FlashData)(web.NewFlash())
}

// Set message to flash
func (fd *FlashData) Set(key string, msg string, args ...interface{}) {
	(*web.FlashData)(fd).Set(key, msg, args...)
}

// Success writes success message to flash.
func (fd *FlashData) Success(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Success(msg, args...)
}

// Notice writes notice message to flash.
func (fd *FlashData) Notice(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Notice(msg, args...)
}

// Warning writes warning message to flash.
func (fd *FlashData) Warning(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Warning(msg, args...)
}

// Error writes error message to flash.
func (fd *FlashData) Error(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Error(msg, args...)
}

// Store does the saving operation of flash data.
// the data are encoded and saved in cookie.
func (fd *FlashData) Store(c *Controller) {
	(*web.FlashData)(fd).Store((*web.Controller)(c))
}

// ReadFromRequest parsed flash data from encoded values in cookie.
func ReadFromRequest(c *Controller) *FlashData {
	return (*FlashData)(web.ReadFromRequest((*web.Controller)(c)))
}
