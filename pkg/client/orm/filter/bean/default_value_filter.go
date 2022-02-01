package bean

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
	"reflect"
	"strings"

	logs "github.com/bhojpur/logger/pkg/engine"

	"github.com/bhojpur/web/pkg/client/orm"
	"github.com/bhojpur/web/pkg/core/bean"
)

// DefaultValueFilterChainBuilder only works for InsertXXX method,
// But InsertOrUpdate and InsertOrUpdateWithCtx is more dangerous than other methods.
// so we won't handle those two methods unless you set includeInsertOrUpdate to true
// And if the element is not pointer, this filter doesn't work
type DefaultValueFilterChainBuilder struct {
	factory                bean.AutoWireBeanFactory
	compatibleWithOldStyle bool

	// only the includeInsertOrUpdate is true, this filter will handle those two methods
	includeInsertOrUpdate bool
}

// NewDefaultValueFilterChainBuilder will create an instance of DefaultValueFilterChainBuilder
// In Bhojpur Web v1.x, the default value config looks like orm:default(xxxx)
// But the default value in 2.x is default:xxx
// so if you want to be compatible with v1.x, please pass true as compatibleWithOldStyle
func NewDefaultValueFilterChainBuilder(typeAdapters map[string]bean.TypeAdapter,
	includeInsertOrUpdate bool,
	compatibleWithOldStyle bool) *DefaultValueFilterChainBuilder {
	factory := bean.NewTagAutoWireBeanFactory()

	if compatibleWithOldStyle {
		newParser := factory.FieldTagParser
		factory.FieldTagParser = func(field reflect.StructField) *bean.FieldMetadata {
			if newParser != nil && field.Tag.Get(bean.DefaultValueTagKey) != "" {
				return newParser(field)
			} else {
				res := &bean.FieldMetadata{}
				ormMeta := field.Tag.Get("orm")
				ormMetaParts := strings.Split(ormMeta, ";")
				for _, p := range ormMetaParts {
					if strings.HasPrefix(p, "default(") && strings.HasSuffix(p, ")") {
						res.DftValue = p[8 : len(p)-1]
					}
				}
				return res
			}
		}
	}

	for k, v := range typeAdapters {
		factory.Adapters[k] = v
	}

	return &DefaultValueFilterChainBuilder{
		factory:                factory,
		compatibleWithOldStyle: compatibleWithOldStyle,
		includeInsertOrUpdate:  includeInsertOrUpdate,
	}
}

func (d *DefaultValueFilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {
	return func(ctx context.Context, inv *orm.Invocation) []interface{} {
		switch inv.Method {
		case "Insert", "InsertWithCtx":
			d.handleInsert(ctx, inv)
			break
		case "InsertOrUpdate", "InsertOrUpdateWithCtx":
			d.handleInsertOrUpdate(ctx, inv)
			break
		case "InsertMulti", "InsertMultiWithCtx":
			d.handleInsertMulti(ctx, inv)
			break
		}
		return next(ctx, inv)
	}
}

func (d *DefaultValueFilterChainBuilder) handleInsert(ctx context.Context, inv *orm.Invocation) {
	d.setDefaultValue(ctx, inv.Args[0])
}

func (d *DefaultValueFilterChainBuilder) handleInsertOrUpdate(ctx context.Context, inv *orm.Invocation) {
	if d.includeInsertOrUpdate {
		ins := inv.Args[0]
		if ins == nil {
			return
		}

		pkName := inv.GetPkFieldName()
		pkField := reflect.Indirect(reflect.ValueOf(ins)).FieldByName(pkName)

		if pkField.IsZero() {
			d.setDefaultValue(ctx, ins)
		}
	}
}

func (d *DefaultValueFilterChainBuilder) handleInsertMulti(ctx context.Context, inv *orm.Invocation) {
	mds := inv.Args[1]

	if t := reflect.TypeOf(mds).Kind(); t != reflect.Array && t != reflect.Slice {
		// do nothing
		return
	}

	mdsArr := reflect.Indirect(reflect.ValueOf(mds))
	for i := 0; i < mdsArr.Len(); i++ {
		d.setDefaultValue(ctx, mdsArr.Index(i).Interface())
	}
	logs.Warn("%v", mdsArr.Index(0).Interface())
}

func (d *DefaultValueFilterChainBuilder) setDefaultValue(ctx context.Context, ins interface{}) {
	err := d.factory.AutoWire(ctx, nil, ins)
	if err != nil {
		logs.Error("try to wire the bean for orm.Insert failed. "+
			"the default value is not set: %v, ", err)
	}
}
