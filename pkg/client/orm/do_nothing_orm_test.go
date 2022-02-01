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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoNothingOrm(t *testing.T) {
	o := &DoNothingOrm{}
	err := o.DoTxWithCtxAndOpts(nil, nil, nil)
	assert.Nil(t, err)

	err = o.DoTxWithCtx(nil, nil)
	assert.Nil(t, err)

	err = o.DoTx(nil)
	assert.Nil(t, err)

	err = o.DoTxWithOpts(nil, nil)
	assert.Nil(t, err)

	assert.Nil(t, o.Driver())

	assert.Nil(t, o.QueryM2MWithCtx(nil, nil, ""))
	assert.Nil(t, o.QueryM2M(nil, ""))
	assert.Nil(t, o.ReadWithCtx(nil, nil))
	assert.Nil(t, o.Read(nil))

	txOrm, err := o.BeginWithCtxAndOpts(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.BeginWithCtx(nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.BeginWithOpts(nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.Begin()
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	assert.Nil(t, o.RawWithCtx(nil, ""))
	assert.Nil(t, o.Raw(""))

	i, err := o.InsertMulti(0, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.Insert(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertOrUpdateWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertOrUpdate(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertMultiWithCtx(nil, 0, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.LoadRelatedWithCtx(nil, nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.LoadRelated(nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	assert.Nil(t, o.QueryTableWithCtx(nil, nil))
	assert.Nil(t, o.QueryTable(nil))

	assert.Nil(t, o.Read(nil))
	assert.Nil(t, o.ReadWithCtx(nil, nil))
	assert.Nil(t, o.ReadForUpdateWithCtx(nil, nil))
	assert.Nil(t, o.ReadForUpdate(nil))

	ok, i, err := o.ReadOrCreate(nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)
	assert.False(t, ok)

	ok, i, err = o.ReadOrCreateWithCtx(nil, nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)
	assert.False(t, ok)

	i, err = o.Delete(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.DeleteWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.Update(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.UpdateWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	assert.Nil(t, o.DBStats())

	to := &DoNothingTxOrm{}
	assert.Nil(t, to.Commit())
	assert.Nil(t, to.Rollback())
}
