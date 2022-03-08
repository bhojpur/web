package app

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

// Tagger is the interface that describes a collection of tags that gives
// context to something.
type Tagger interface {
	// Returns a collection of tags.
	Tags() Tags
}

// Tags represent key-value pairs that give context to what they are used with.
type Tags map[string]string

func (t Tags) Tags() Tags {
	return t
}

// Set sets a tag with the given name and value. The value is converted to a
// string.
func (t Tags) Set(name string, v interface{}) {
	t[name] = toString(v)
}

// Get returns a tag value with the given name.
func (t Tags) Get(name string) string {
	return t[name]
}

// Tag is a key-value pair that adds context to an action.
type Tag struct {
	Name  string
	Value string
}

func (t Tag) Tags() Tags {
	return Tags{t.Name: t.Value}
}

// T creates a tag with the given name and value. The value is converted to a
// string.
func T(name string, value interface{}) Tag {
	return Tag{
		Name:  name,
		Value: toString(value),
	}
}
