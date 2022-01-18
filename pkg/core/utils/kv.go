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

type KV interface {
	GetKey() interface{}
	GetValue() interface{}
}

// SimpleKV is common structure to store key-value pairs.
// When you need something like Pair, you can use this
type SimpleKV struct {
	Key   interface{}
	Value interface{}
}

var _ KV = new(SimpleKV)

func (s *SimpleKV) GetKey() interface{} {
	return s.Key
}

func (s *SimpleKV) GetValue() interface{} {
	return s.Value
}

// KVs interface
type KVs interface {
	GetValueOr(key interface{}, defValue interface{}) interface{}
	Contains(key interface{}) bool
	IfContains(key interface{}, action func(value interface{})) KVs
}

// SimpleKVs will store SimpleKV collection as map
type SimpleKVs struct {
	kvs map[interface{}]interface{}
}

var _ KVs = new(SimpleKVs)

// GetValueOr returns the value for a given key, if non-existant
// it returns defValue
func (kvs *SimpleKVs) GetValueOr(key interface{}, defValue interface{}) interface{} {
	v, ok := kvs.kvs[key]
	if ok {
		return v
	}
	return defValue
}

// Contains checks if a key exists
func (kvs *SimpleKVs) Contains(key interface{}) bool {
	_, ok := kvs.kvs[key]
	return ok
}

// IfContains invokes the action on a key if it exists
func (kvs *SimpleKVs) IfContains(key interface{}, action func(value interface{})) KVs {
	v, ok := kvs.kvs[key]
	if ok {
		action(v)
	}
	return kvs
}

// NewKVs creates the *KVs instance
func NewKVs(kvs ...KV) KVs {
	res := &SimpleKVs{
		kvs: make(map[interface{}]interface{}, len(kvs)),
	}
	for _, kv := range kvs {
		res.kvs[kv.GetKey()] = kv.GetValue()
	}
	return res
}
