package form

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
	"errors"
	"net/url"
	"strconv"
)

type FieldTemplate interface {
	GetName() string
	Parse(unparsed string) (interface{}, error)
	Validate(value interface{}) error
}

type TooShortError struct {
	Minimum int
}

func (e TooShortError) Error() string {
	return "Too short"
}

type TooLongError struct {
	Maximum int
}

func (e TooLongError) Error() string {
	return "Too long"
}

type StringTemplate struct {
	Name      string
	Required  bool
	MinLength int
	MaxLength int
}

func (f *StringTemplate) GetName() string {
	return f.Name
}

func (f *StringTemplate) Parse(unparsed string) (interface{}, error) {
	return unparsed, nil
}

func (f *StringTemplate) Validate(value interface{}) (err error) {
	if value == nil || value == "" {
		if f.Required {
			return errors.New("is missing")
		} else {
			return nil
		}
	}

	v := value.(string)

	if len(v) < f.MinLength {
		return TooShortError{Minimum: f.MinLength}
	}

	if f.MaxLength < len(v) {
		return TooLongError{Maximum: f.MaxLength}
	}

	return
}

type IntTemplate struct {
	Name     string
	Required bool
	Minimum  int64
	Maximum  int64
}

func (f *IntTemplate) GetName() string {
	return f.Name
}

func (f *IntTemplate) Parse(unparsed string) (interface{}, error) {
	if unparsed == "" {
		return nil, nil
	}

	if parsed, err := strconv.ParseInt(unparsed, 10, 64); err == nil {
		return parsed, err
	} else {
		return nil, err
	}
}

func (f *IntTemplate) Validate(value interface{}) error {
	if f.Required && value == nil {
		return errors.New("is missing")
	}

	v := value.(int64)

	if v < f.Minimum {
		return errors.New("too small")
	}

	if f.Maximum < v {
		return errors.New("too big")
	}

	return nil
}

type FormTemplate struct {
	fieldTemplates map[string]FieldTemplate
	CustomValidate func(*Form)
}

func NewFormTemplate() (f *FormTemplate) {
	f = &FormTemplate{}
	f.fieldTemplates = make(map[string]FieldTemplate)
	return
}

func (f *FormTemplate) AddField(fieldTemplate FieldTemplate) {
	f.fieldTemplates[fieldTemplate.GetName()] = fieldTemplate
}

func (f *FormTemplate) Parse(values url.Values) (s *Form) {
	s = new(Form)
	s.Fields = make(map[string]*Field, len(f.fieldTemplates))

	for name, field := range f.fieldTemplates {
		var sf Field

		if fieldValues, ok := values[name]; ok {
			unparsed := fieldValues[len(fieldValues)-1]
			sf.Unparsed = unparsed

			parsed, err := field.Parse(unparsed)
			sf.Parsed = parsed
			sf.Error = err
		}

		s.Fields[name] = &sf
	}

	return
}

func (f *FormTemplate) New() (s *Form) {
	s = new(Form)
	s.Fields = make(map[string]*Field, len(f.fieldTemplates))

	for name, field := range f.fieldTemplates {
		var sf Field
		sf.Parsed, _ = field.Parse("")
		s.Fields[name] = &sf
	}

	return
}

func (f *FormTemplate) Validate(s *Form) {
	for name, fieldTemplate := range f.fieldTemplates {
		s.Fields[name].Error = fieldTemplate.Validate(s.Fields[name].Parsed)
	}

	if f.CustomValidate != nil {
		f.CustomValidate(s)
	}
}

type Field struct {
	Unparsed string
	Parsed   interface{}
	Error    error
}

type Form struct {
	Fields map[string]*Field
}

func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if field.Error != nil {
			return false
		}
	}
	return true
}
