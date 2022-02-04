package utils

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
	"io"

	utils "github.com/bhojpur/mail/pkg/engine"
)

// Email is the type used for email messages
type Email utils.Email

// Attachment is a struct representing an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment utils.Attachment

// NewEMail create new Email struct with config json.
// config json is followed from Email struct fields.
func NewEMail(config string) *Email {
	return (*Email)(utils.NewEmail(config))
}

// Bytes Make all send information to byte
func (e *Email) Bytes() ([]byte, error) {
	return (*utils.Email)(e).Bytes()
}

// AttachFile Add attach file to the send mail
func (e *Email) AttachFile(args ...string) (*Attachment, error) {
	a, err := (*utils.Email)(e).AttachFile(args...)
	if err != nil {
		return nil, err
	}
	return (*Attachment)(a), err
}

// Attach is used to attach content from an io.Reader to the email.
// Parameters include an io.Reader, the desired filename for the attachment, and the Content-Type.
func (e *Email) Attach(r io.Reader, filename string, args ...string) (*Attachment, error) {
	a, err := (*utils.Email)(e).Attach(r, filename, args...)
	if err != nil {
		return nil, err
	}
	return (*Attachment)(a), err
}

// Send will send out the mail
func (e *Email) Send() error {
	return (*utils.Email)(e).Send()
}
