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

	session "github.com/bhojpur/session/pkg/engine"
)

type oldToNewProviderAdapter struct {
	delegate Provider
}

func (o *oldToNewProviderAdapter) SessionInit(ctx context.Context, gclifetime int64, config string) error {
	return o.delegate.SessionInit(gclifetime, config)
}

func (o *oldToNewProviderAdapter) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	store, err := o.delegate.SessionRead(sid)
	return &oldToNewStoreAdapter{
		delegate: store,
	}, err
}

func (o *oldToNewProviderAdapter) SessionExist(ctx context.Context, sid string) (bool, error) {
	return o.delegate.SessionExist(sid), nil
}

func (o *oldToNewProviderAdapter) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	s, err := o.delegate.SessionRegenerate(oldsid, sid)
	return &oldToNewStoreAdapter{
		delegate: s,
	}, err
}

func (o *oldToNewProviderAdapter) SessionDestroy(ctx context.Context, sid string) error {
	return o.delegate.SessionDestroy(sid)
}

func (o *oldToNewProviderAdapter) SessionAll(ctx context.Context) int {
	return o.delegate.SessionAll()
}

func (o *oldToNewProviderAdapter) SessionGC(ctx context.Context) {
	o.delegate.SessionGC()
}

type newToOldProviderAdapter struct {
	delegate session.Provider
}

func (n *newToOldProviderAdapter) SessionInit(gclifetime int64, config string) error {
	return n.delegate.SessionInit(context.Background(), gclifetime, config)
}

func (n *newToOldProviderAdapter) SessionRead(sid string) (Store, error) {
	s, err := n.delegate.SessionRead(context.Background(), sid)
	if adt, ok := s.(*oldToNewStoreAdapter); err == nil && ok {
		return adt.delegate, err
	}
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

func (n *newToOldProviderAdapter) SessionExist(sid string) bool {
	res, _ := n.delegate.SessionExist(context.Background(), sid)
	return res
}

func (n *newToOldProviderAdapter) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := n.delegate.SessionRegenerate(context.Background(), oldsid, sid)
	if adt, ok := s.(*oldToNewStoreAdapter); err == nil && ok {
		return adt.delegate, err
	}
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

func (n *newToOldProviderAdapter) SessionDestroy(sid string) error {
	return n.delegate.SessionDestroy(context.Background(), sid)
}

func (n *newToOldProviderAdapter) SessionAll() int {
	return n.delegate.SessionAll(context.Background())
}

func (n *newToOldProviderAdapter) SessionGC() {
	n.delegate.SessionGC(context.Background())
}
