package order_clause

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
	"strings"

	"github.com/bhojpur/web/pkg/client/orm/clauses"
)

type Sort int8

const (
	None       Sort = 0
	Ascending  Sort = 1
	Descending Sort = 2
)

type Option func(order *Order)

type Order struct {
	column string
	sort   Sort
	isRaw  bool
}

func Clause(options ...Option) *Order {
	o := &Order{}
	for _, option := range options {
		option(o)
	}

	return o
}

func (o *Order) GetColumn() string {
	return o.column
}

func (o *Order) GetSort() Sort {
	return o.sort
}

func (o *Order) SortString() string {
	switch o.GetSort() {
	case Ascending:
		return "ASC"
	case Descending:
		return "DESC"
	}

	return ``
}

func (o *Order) IsRaw() bool {
	return o.isRaw
}

func ParseOrder(expressions ...string) []*Order {
	var orders []*Order
	for _, expression := range expressions {
		sort := Ascending
		column := strings.ReplaceAll(expression, clauses.ExprSep, clauses.ExprDot)
		if column[0] == '-' {
			sort = Descending
			column = column[1:]
		}

		orders = append(orders, &Order{
			column: column,
			sort:   sort,
		})
	}

	return orders
}

func Column(column string) Option {
	return func(order *Order) {
		order.column = strings.ReplaceAll(column, clauses.ExprSep, clauses.ExprDot)
	}
}

func sort(sort Sort) Option {
	return func(order *Order) {
		order.sort = sort
	}
}

func SortAscending() Option {
	return sort(Ascending)
}

func SortDescending() Option {
	return sort(Descending)
}

func SortNone() Option {
	return sort(None)
}

func Raw() Option {
	return func(order *Order) {
		order.isRaw = true
	}
}
