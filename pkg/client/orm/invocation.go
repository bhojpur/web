package orm

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
	"context"
	"time"
)

// Invocation represents an "ORM" invocation
type Invocation struct {
	Method string
	// Md may be nil in some cases. It depends on method
	Md interface{}
	// the args are all arguments except context.Context
	Args []interface{}

	mi *modelInfo
	// f is the Orm operation
	f func(ctx context.Context) []interface{}

	// insideTx indicates whether this is inside a transaction
	InsideTx    bool
	TxStartTime time.Time
	TxName      string
}

func (inv *Invocation) GetTableName() string {
	if inv.mi != nil {
		return inv.mi.table
	}
	return ""
}

func (inv *Invocation) execute(ctx context.Context) []interface{} {
	return inv.f(ctx)
}

// GetPkFieldName return the primary key of this table
// if not found, "" is returned
func (inv *Invocation) GetPkFieldName() string {
	if inv.mi.fields.pk != nil {
		return inv.mi.fields.pk.name
	}
	return ""
}
