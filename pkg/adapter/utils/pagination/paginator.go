package pagination

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
	"net/http"

	"github.com/bhojpur/web/pkg/core/utils/pagination"
)

// Paginator within the state of a http request.
type Paginator pagination.Paginator

// PageNums Returns the total number of pages.
func (p *Paginator) PageNums() int {
	return (*pagination.Paginator)(p).PageNums()
}

// Nums Returns the total number of items (e.g. from doing SQL count).
func (p *Paginator) Nums() int64 {
	return (*pagination.Paginator)(p).Nums()
}

// SetNums Sets the total number of items.
func (p *Paginator) SetNums(nums interface{}) {
	(*pagination.Paginator)(p).SetNums(nums)
}

// Page Returns the current page.
func (p *Paginator) Page() int {
	return (*pagination.Paginator)(p).Page()
}

// Pages Returns a list of all pages.
//
// Usage (in a view template):
//
//  {{range $index, $page := .paginator.Pages}}
//    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
//      <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
//    </li>
//  {{end}}
func (p *Paginator) Pages() []int {
	return (*pagination.Paginator)(p).Pages()
}

// PageLink Returns URL for a given page index.
func (p *Paginator) PageLink(page int) string {
	return (*pagination.Paginator)(p).PageLink(page)
}

// PageLinkPrev Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	return (*pagination.Paginator)(p).PageLinkPrev()
}

// PageLinkNext Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	return (*pagination.Paginator)(p).PageLinkNext()
}

// PageLinkFirst Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return (*pagination.Paginator)(p).PageLinkFirst()
}

// PageLinkLast Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return (*pagination.Paginator)(p).PageLinkLast()
}

// HasPrev Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return (*pagination.Paginator)(p).HasPrev()
}

// HasNext Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return (*pagination.Paginator)(p).HasNext()
}

// IsActive Returns true if the given page index points to the current page.
func (p *Paginator) IsActive(page int) bool {
	return (*pagination.Paginator)(p).IsActive(page)
}

// Offset Returns the current offset.
func (p *Paginator) Offset() int {
	return (*pagination.Paginator)(p).Offset()
}

// HasPages Returns true if there is more than one page.
func (p *Paginator) HasPages() bool {
	return (*pagination.Paginator)(p).HasPages()
}

// NewPaginator Instantiates a paginator struct for the current http request.
func NewPaginator(req *http.Request, per int, nums interface{}) *Paginator {
	return (*Paginator)(pagination.NewPaginator(req, per, nums))
}
