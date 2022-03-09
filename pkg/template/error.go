package template

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
	"bufio"
	"fmt"
	"os"
)

// The Error type is being used to address an error during lexing, parsing or
// execution. If you want to return an error object (for example in your own
// tag or filter) fill this object with as much information as you have.
// Make sure "Sender" is always given (if you're returning an error within
// a filter, make Sender equals 'filter:yourfilter'; same goes for tags: 'tag:mytag').
// It's okay if you only fill in ErrorMsg if you don't have any other details at hand.
type Error struct {
	Template  *Template
	Filename  string
	Line      int
	Column    int
	Token     *Token
	Sender    string
	OrigError error
}

func (e *Error) updateFromTokenIfNeeded(template *Template, t *Token) *Error {
	if e.Template == nil {
		e.Template = template
	}

	if e.Token == nil {
		e.Token = t
		if e.Line <= 0 {
			e.Line = t.Line
			e.Column = t.Col
		}
	}

	return e
}

// Returns a nice formatted error string.
func (e *Error) Error() string {
	s := "[Error"
	if e.Sender != "" {
		s += " (where: " + e.Sender + ")"
	}
	if e.Filename != "" {
		s += " in " + e.Filename
	}
	if e.Line > 0 {
		s += fmt.Sprintf(" | Line %d Col %d", e.Line, e.Column)
		if e.Token != nil {
			s += fmt.Sprintf(" near '%s'", e.Token.Val)
		}
	}
	s += "] "
	s += e.OrigError.Error()
	return s
}

// RawLine returns the affected line from the original template, if available.
func (e *Error) RawLine() (line string, available bool, outErr error) {
	if e.Line <= 0 || e.Filename == "<string>" {
		return "", false, nil
	}

	filename := e.Filename
	if e.Template != nil {
		filename = e.Template.set.resolveFilename(e.Template, e.Filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return "", false, err
	}
	defer func() {
		err := file.Close()
		if err != nil && outErr == nil {
			outErr = err
		}
	}()

	scanner := bufio.NewScanner(file)
	l := 0
	for scanner.Scan() {
		l++
		if l == e.Line {
			return scanner.Text(), true, nil
		}
	}
	return "", false, nil
}
