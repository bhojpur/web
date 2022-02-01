package session

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
	"net/http"

	session "github.com/bhojpur/session/pkg/engine"
)

type NewToOldStoreAdapter struct {
	delegate session.Store
}

func CreateNewToOldStoreAdapter(s session.Store) Store {
	return &NewToOldStoreAdapter{
		delegate: s,
	}
}

func (n *NewToOldStoreAdapter) Set(key, value interface{}) error {
	return n.delegate.Set(context.Background(), key, value)
}

func (n *NewToOldStoreAdapter) Get(key interface{}) interface{} {
	return n.delegate.Get(context.Background(), key)
}

func (n *NewToOldStoreAdapter) Delete(key interface{}) error {
	return n.delegate.Delete(context.Background(), key)
}

func (n *NewToOldStoreAdapter) SessionID() string {
	return n.delegate.SessionID(context.Background())
}

func (n *NewToOldStoreAdapter) SessionRelease(w http.ResponseWriter) {
	n.delegate.SessionRelease(context.Background(), w)
}

func (n *NewToOldStoreAdapter) Flush() error {
	return n.delegate.Flush(context.Background())
}

type oldToNewStoreAdapter struct {
	delegate Store
}

func (o *oldToNewStoreAdapter) Set(ctx context.Context, key, value interface{}) error {
	return o.delegate.Set(key, value)
}

func (o *oldToNewStoreAdapter) Get(ctx context.Context, key interface{}) interface{} {
	return o.delegate.Get(key)
}

func (o *oldToNewStoreAdapter) Delete(ctx context.Context, key interface{}) error {
	return o.delegate.Delete(key)
}

func (o *oldToNewStoreAdapter) SessionID(ctx context.Context) string {
	return o.delegate.SessionID()
}

func (o *oldToNewStoreAdapter) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	o.delegate.SessionRelease(w)
}

func (o *oldToNewStoreAdapter) Flush(ctx context.Context) error {
	return o.delegate.Flush()
}
