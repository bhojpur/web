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
	"database/sql"

	"github.com/bhojpur/web/pkg/core/utils"
)

// DoNothingOrm won't do anything, usually you use this to custom your mock Ormer implementation
// I think golang mocking interface is hard to use
// this may help you to integrate with Ormer

var _ Ormer = new(DoNothingOrm)

type DoNothingOrm struct {
}

func (d *DoNothingOrm) Read(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadForUpdate(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return false, 0, nil
}

func (d *DoNothingOrm) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return false, 0, nil
}

func (d *DoNothingOrm) LoadRelated(md interface{}, name string, args ...utils.KV) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils.KV) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) QueryM2M(md interface{}, name string) QueryM2Mer {
	return nil
}

func (d *DoNothingOrm) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) QueryM2Mer {
	return nil
}

func (d *DoNothingOrm) QueryTable(ptrStructOrTableName interface{}) QuerySetter {
	return nil
}

func (d *DoNothingOrm) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) QuerySetter {
	return nil
}

func (d *DoNothingOrm) DBStats() *sql.DBStats {
	return nil
}

func (d *DoNothingOrm) Insert(md interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Update(md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Delete(md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Raw(query string, args ...interface{}) RawSetter {
	return nil
}

func (d *DoNothingOrm) RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSetter {
	return nil
}

func (d *DoNothingOrm) Driver() Driver {
	return nil
}

func (d *DoNothingOrm) Begin() (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithCtx(ctx context.Context) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

// DoNothingTxOrm is similar with DoNothingOrm, usually you use it to test
type DoNothingTxOrm struct {
	DoNothingOrm
}

func (d *DoNothingTxOrm) Commit() error {
	return nil
}

func (d *DoNothingTxOrm) Rollback() error {
	return nil
}
