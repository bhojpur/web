package engine

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
	"strings"

	ctxsvr "github.com/bhojpur/web/pkg/context"
)

type namespaceCond func(*ctxsvr.Context) bool

// LinkNamespace used as link action
type LinkNamespace func(*Namespace)

// Namespace is store all the info
type Namespace struct {
	prefix   string
	handlers *ControllerRegister
}

// NewNamespace get new Namespace
func NewNamespace(prefix string, params ...LinkNamespace) *Namespace {
	ns := &Namespace{
		prefix:   prefix,
		handlers: NewControllerRegister(),
	}
	for _, p := range params {
		p(ns)
	}
	return ns
}

// Cond set condition function
// if cond return true can run this namespace, else can't
// usage:
// ns.Cond(func (ctx *context.Context) bool{
//       if ctx.Input.Domain() == "api.bhojpur.net" {
//         return true
//       }
//       return false
//   })
// Cond as the first filter
func (n *Namespace) Cond(cond namespaceCond) *Namespace {
	fn := func(ctx *ctxsvr.Context) {
		if !cond(ctx) {
			exception("405", ctx)
		}
	}
	if v := n.handlers.filters[BeforeRouter]; len(v) > 0 {
		mr := new(FilterRouter)
		mr.tree = NewTree()
		mr.pattern = "*"
		mr.filterFunc = fn
		mr.tree.AddRouter("*", true)
		n.handlers.filters[BeforeRouter] = append([]*FilterRouter{mr}, v...)
	} else {
		n.handlers.InsertFilter("*", BeforeRouter, fn)
	}
	return n
}

// Filter add filter in the Namespace
// action has before & after
// FilterFunc
// usage:
// Filter("before", func (ctx *context.Context){
//       _, ok := ctx.Input.Session("uid").(int)
//       if !ok && ctx.Request.RequestURI != "/login" {
//          ctx.Redirect(302, "/login")
//        }
//   })
func (n *Namespace) Filter(action string, filter ...FilterFunc) *Namespace {
	var a int
	if action == "before" {
		a = BeforeRouter
	} else if action == "after" {
		a = FinishRouter
	}
	for _, f := range filter {
		n.handlers.InsertFilter("*", a, f, WithReturnOnOutput(true))
	}
	return n
}

// Router same as bhojpur.Router
// refer: https://godoc.org/github.com/bhojpur/web#Router
func (n *Namespace) Router(rootpath string, c ControllerInterface, mappingMethods ...string) *Namespace {
	n.handlers.Add(rootpath, c, mappingMethods...)
	return n
}

// AutoRouter same as bhojpur.AutoRouter
// refer: https://godoc.org/github.com/bhojpur/web#AutoRouter
func (n *Namespace) AutoRouter(c ControllerInterface) *Namespace {
	n.handlers.AddAuto(c)
	return n
}

// AutoPrefix same as bhojpur.AutoPrefix
// refer: https://godoc.org/github.com/bhojpur/web#AutoPrefix
func (n *Namespace) AutoPrefix(prefix string, c ControllerInterface) *Namespace {
	n.handlers.AddAutoPrefix(prefix, c)
	return n
}

// Get same as bhojpur.Get
// refer: https://godoc.org/github.com/bhojpur/web#Get
func (n *Namespace) Get(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Get(rootpath, f)
	return n
}

// Post same as bhojpur.Post
// refer: https://godoc.org/github.com/bhojpur/web#Post
func (n *Namespace) Post(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Post(rootpath, f)
	return n
}

// Delete same as bhojpur.Delete
// refer: https://godoc.org/github.com/bhojpur/web#Delete
func (n *Namespace) Delete(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Delete(rootpath, f)
	return n
}

// Put same as bhojpur.Put
// refer: https://godoc.org/github.com/bhojpur/web#Put
func (n *Namespace) Put(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Put(rootpath, f)
	return n
}

// Head same as bhojpur.Head
// refer: https://godoc.org/github.com/bhojpur/web#Head
func (n *Namespace) Head(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Head(rootpath, f)
	return n
}

// Options same as bhojpur.Options
// refer: https://godoc.org/github.com/bhojpur/web#Options
func (n *Namespace) Options(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Options(rootpath, f)
	return n
}

// Patch same as bhojpur.Patch
// refer: https://godoc.org/github.com/bhojpur/web#Patch
func (n *Namespace) Patch(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Patch(rootpath, f)
	return n
}

// Any same as bhojpur.Any
// refer: https://godoc.org/github.com/bhojpur/web#Any
func (n *Namespace) Any(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Any(rootpath, f)
	return n
}

// Handler same as bhojpur.Handler
// refer: https://godoc.org/github.com/bhojpur/web#Handler
func (n *Namespace) Handler(rootpath string, h http.Handler) *Namespace {
	n.handlers.Handler(rootpath, h)
	return n
}

// Include add include class
// refer: https://godoc.org/github.com/bhojpur/web#Include
func (n *Namespace) Include(cList ...ControllerInterface) *Namespace {
	n.handlers.Include(cList...)
	return n
}

// Namespace add nest Namespace
// usage:
//ns := bhojpur.NewNamespace(“/v1”).
//Namespace(
//    bhojpur.NewNamespace("/shop").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("shopinfo"))
//    }),
//    bhojpur.NewNamespace("/order").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("orderinfo"))
//    }),
//    bhojpur.NewNamespace("/crm").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("crminfo"))
//    }),
//)
func (n *Namespace) Namespace(ns ...*Namespace) *Namespace {
	for _, ni := range ns {
		for k, v := range ni.handlers.routers {
			if _, ok := n.handlers.routers[k]; ok {
				addPrefix(v, ni.prefix)
				n.handlers.routers[k].AddTree(ni.prefix, v)
			} else {
				t := NewTree()
				t.AddTree(ni.prefix, v)
				addPrefix(t, ni.prefix)
				n.handlers.routers[k] = t
			}
		}
		if ni.handlers.enableFilter {
			for pos, filterList := range ni.handlers.filters {
				for _, mr := range filterList {
					t := NewTree()
					t.AddTree(ni.prefix, mr.tree)
					mr.tree = t
					n.handlers.insertFilterRouter(pos, mr)
				}
			}
		}
	}
	return n
}

// AddNamespace register Namespace into bhojpur.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		for k, v := range n.handlers.routers {
			if _, ok := BhojpurApp.Handlers.routers[k]; ok {
				addPrefix(v, n.prefix)
				BhojpurApp.Handlers.routers[k].AddTree(n.prefix, v)
			} else {
				t := NewTree()
				t.AddTree(n.prefix, v)
				addPrefix(t, n.prefix)
				BhojpurApp.Handlers.routers[k] = t
			}
		}
		if n.handlers.enableFilter {
			for pos, filterList := range n.handlers.filters {
				for _, mr := range filterList {
					t := NewTree()
					t.AddTree(n.prefix, mr.tree)
					mr.tree = t
					BhojpurApp.Handlers.insertFilterRouter(pos, mr)
				}
			}
		}
	}
}

func addPrefix(t *Tree, prefix string) {
	for _, v := range t.fixrouters {
		addPrefix(v, prefix)
	}
	if t.wildcard != nil {
		addPrefix(t.wildcard, prefix)
	}
	for _, l := range t.leaves {
		if c, ok := l.runObject.(*ControllerInfo); ok {
			if !strings.HasPrefix(c.pattern, prefix) {
				c.pattern = prefix + c.pattern
			}
		}
	}
}

// NSCond is Namespace Condition
func NSCond(cond namespaceCond) LinkNamespace {
	return func(ns *Namespace) {
		ns.Cond(cond)
	}
}

// NSBefore Namespace BeforeRouter filter
func NSBefore(filterList ...FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Filter("before", filterList...)
	}
}

// NSAfter add Namespace FinishRouter filter
func NSAfter(filterList ...FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Filter("after", filterList...)
	}
}

// NSInclude Namespace Include ControllerInterface
func NSInclude(cList ...ControllerInterface) LinkNamespace {
	return func(ns *Namespace) {
		ns.Include(cList...)
	}
}

// NSRouter call Namespace Router
func NSRouter(rootpath string, c ControllerInterface, mappingMethods ...string) LinkNamespace {
	return func(ns *Namespace) {
		ns.Router(rootpath, c, mappingMethods...)
	}
}

// NSGet call Namespace Get
func NSGet(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Get(rootpath, f)
	}
}

// NSPost call Namespace Post
func NSPost(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Post(rootpath, f)
	}
}

// NSHead call Namespace Head
func NSHead(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Head(rootpath, f)
	}
}

// NSPut call Namespace Put
func NSPut(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Put(rootpath, f)
	}
}

// NSDelete call Namespace Delete
func NSDelete(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Delete(rootpath, f)
	}
}

// NSAny call Namespace Any
func NSAny(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Any(rootpath, f)
	}
}

// NSOptions call Namespace Options
func NSOptions(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Options(rootpath, f)
	}
}

// NSPatch call Namespace Patch
func NSPatch(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		ns.Patch(rootpath, f)
	}
}

// NSAutoRouter call Namespace AutoRouter
func NSAutoRouter(c ControllerInterface) LinkNamespace {
	return func(ns *Namespace) {
		ns.AutoRouter(c)
	}
}

// NSAutoPrefix call Namespace AutoPrefix
func NSAutoPrefix(prefix string, c ControllerInterface) LinkNamespace {
	return func(ns *Namespace) {
		ns.AutoPrefix(prefix, c)
	}
}

// NSNamespace add sub Namespace
func NSNamespace(prefix string, params ...LinkNamespace) LinkNamespace {
	return func(ns *Namespace) {
		n := NewNamespace(prefix, params...)
		ns.Namespace(n)
	}
}

// NSHandler add handler
func NSHandler(rootpath string, h http.Handler) LinkNamespace {
	return func(ns *Namespace) {
		ns.Handler(rootpath, h)
	}
}
