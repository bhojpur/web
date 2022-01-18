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
	"testing"
)

func TestStringTemplateValidate(t *testing.T) {
	var MissingError = errors.New("is missing")
	var st *StringTemplate

	st = &StringTemplate{MinLength: 4, MaxLength: 10}
	if err := st.Validate(nil); err != nil {
		t.Errorf("Unexpected validate result: %v", err)
	}
	if err := st.Validate(""); err != nil {
		t.Errorf("Unexpected validate result: %v", err)
	}

	st = &StringTemplate{Required: true, MaxLength: 100}
	if err := st.Validate(nil); err != MissingError {
		t.Errorf("Unexpected validate result: %v", err)
	}
	if err := st.Validate(""); err != MissingError {
		t.Errorf("Unexpected validate result: %v", err)
	}
	if err := st.Validate("John"); err != nil {
		t.Errorf("Unexpected validate result: %v", err)
	}

	st = &StringTemplate{MinLength: 4, MaxLength: 6}
	if err := st.Validate("Sam"); err != (TooShortError{Minimum: 4}) {
		t.Errorf("Unexpected validate result: %v", err)
	}
	if err := st.Validate("John"); err != nil {
		t.Errorf("Unexpected validate result: %v", err)
	}
	if err := st.Validate("Alexender"); err != (TooLongError{Maximum: 6}) {
		t.Errorf("Unexpected validate result: %v", err)
	}
}

func TestFormTemplateParse(t *testing.T) {
	formTemplate := NewFormTemplate()
	formTemplate.AddField(&StringTemplate{Name: "firstName"})
	formTemplate.AddField(&StringTemplate{Name: "lastName"})

	values := make(url.Values)
	values["firstName"] = []string{"John"}
	values["lastName"] = []string{"Smith"}

	form := formTemplate.Parse(values)

	if actual := form.Fields["firstName"].Unparsed; actual != "John" {
		t.Errorf(`Expected unparsed "firstName" to be "John" but it was %#v`, actual)
	}

	if actual := form.Fields["firstName"].Parsed; actual != "John" {
		t.Errorf(`Expected parsed "firstName" to be "John" but it was %#v`, actual)
	}

	if actual := form.Fields["lastName"].Unparsed; actual != "Smith" {
		t.Errorf(`Expected unparsed "lastName" to be "Smith" but it was %#v`, actual)
	}

	if actual := form.Fields["lastName"].Parsed; actual != "Smith" {
		t.Errorf(`Expected parsed "lastName" to be "Smith" but it was %#v`, actual)
	}
}

func TestFormTemplateParseWithParseError(t *testing.T) {
	formTemplate := NewFormTemplate()
	formTemplate.AddField(&IntTemplate{Name: "age"})

	values := make(url.Values)
	values["age"] = []string{"foo"}

	form := formTemplate.Parse(values)

	if actual := form.Fields["age"].Unparsed; actual != "foo" {
		t.Errorf(`Expected unparsed "age" to be "foo" but it was %#v`, actual)
	}

	if actual := form.Fields["age"].Parsed; actual != nil {
		t.Errorf(`Expected parsed "actual" to be <nil> but it was %#v`, actual)
	}
}

func TestFormTemplateNew(t *testing.T) {
	formTemplate := NewFormTemplate()
	formTemplate.AddField(&StringTemplate{Name: "name"})
	form := formTemplate.New()

	if actual := form.Fields["name"].Unparsed; actual != "" {
		t.Errorf(`Expected empty "name" to be "" but it was %#v`, actual)
	}

	if actual := form.Fields["name"].Parsed; actual != "" {
		t.Errorf(`Expected empty "name" to be "" but it was %#v`, actual)
	}
}

func TestFormTemplateValidate(t *testing.T) {
	var MissingError = errors.New("is missing")
	formTemplate := NewFormTemplate()
	formTemplate.AddField(&StringTemplate{Name: "name", Required: true})
	form := formTemplate.New()
	formTemplate.Validate(form)

	if form.Fields["name"].Error != MissingError {
		t.Error("Validate didn't")
	}

	formTemplate = NewFormTemplate()
	formTemplate.AddField(&StringTemplate{Name: "name"})
	formTemplate.CustomValidate = func(f *Form) {
		f.Fields["name"].Error = errors.New("Custom Error")
	}
	form = formTemplate.New()
	formTemplate.Validate(form)

	if form.Fields["name"].Error.Error() != "Custom Error" {
		t.Errorf(`Expected "Custom Error" but it ws %#v`, form.Fields["name"].Error)
	}
}
