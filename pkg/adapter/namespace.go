package adapter

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

	adtContext "github.com/bhojpur/web/pkg/adapter/context"
	"github.com/bhojpur/web/pkg/context"

	web "github.com/bhojpur/web/pkg/engine"
)

type namespaceCond func(*adtContext.Context) bool

// LinkNamespace used as link action
type LinkNamespace func(*Namespace)

// Namespace is store all the info
type Namespace web.Namespace

// NewNamespace get new Namespace
func NewNamespace(prefix string, params ...LinkNamespace) *Namespace {
	nps := oldToNewLinkNs(params)
	return (*Namespace)(web.NewNamespace(prefix, nps...))
}

func oldToNewLinkNs(params []LinkNamespace) []web.LinkNamespace {
	nps := make([]web.LinkNamespace, 0, len(params))
	for _, p := range params {
		nps = append(nps, func(namespace *web.Namespace) {
			p((*Namespace)(namespace))
		})
	}
	return nps
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
	(*web.Namespace)(n).Cond(func(context *context.Context) bool {
		return cond((*adtContext.Context)(context))
	})
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
	nfs := oldToNewFilter(filter)
	(*web.Namespace)(n).Filter(action, nfs...)
	return n
}

func oldToNewFilter(filter []FilterFunc) []web.FilterFunc {
	nfs := make([]web.FilterFunc, 0, len(filter))
	for _, f := range filter {
		nfs = append(nfs, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
	return nfs
}

// Router same as bhojpur.Rourer
func (n *Namespace) Router(rootpath string, c ControllerInterface, mappingMethods ...string) *Namespace {
	(*web.Namespace)(n).Router(rootpath, c, mappingMethods...)
	return n
}

// AutoRouter same as bhojpur.AutoRouter
func (n *Namespace) AutoRouter(c ControllerInterface) *Namespace {
	(*web.Namespace)(n).AutoRouter(c)
	return n
}

// AutoPrefix same as bhojpur.AutoPrefix
func (n *Namespace) AutoPrefix(prefix string, c ControllerInterface) *Namespace {
	(*web.Namespace)(n).AutoPrefix(prefix, c)
	return n
}

// Get same as bhojpur.Get
func (n *Namespace) Get(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Get(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Post same as bhojpur.Post
func (n *Namespace) Post(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Post(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Delete same as bhojpur.Delete
func (n *Namespace) Delete(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Delete(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Put same as bhojpur.Put
func (n *Namespace) Put(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Put(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Head same as bhojpur.Head
func (n *Namespace) Head(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Head(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Options same as bhojpur.Options
func (n *Namespace) Options(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Options(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Patch same as bhojpur.Patch
func (n *Namespace) Patch(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Patch(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Any same as bhojpur.Any
func (n *Namespace) Any(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Any(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Handler same as bhojpur.Handler
func (n *Namespace) Handler(rootpath string, h http.Handler) *Namespace {
	(*web.Namespace)(n).Handler(rootpath, h)
	return n
}

// Include add include class
func (n *Namespace) Include(cList ...ControllerInterface) *Namespace {
	nL := oldToNewCtrlIntfs(cList)
	(*web.Namespace)(n).Include(nL...)
	return n
}

// Namespace add nest Namespace
// usage:
// ns := bhojpur.NewNamespace(“/v1”).
// Namespace(
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
// )
func (n *Namespace) Namespace(ns ...*Namespace) *Namespace {
	nns := oldToNewNs(ns)
	(*web.Namespace)(n).Namespace(nns...)
	return n
}

func oldToNewNs(ns []*Namespace) []*web.Namespace {
	nns := make([]*web.Namespace, 0, len(ns))
	for _, n := range ns {
		nns = append(nns, (*web.Namespace)(n))
	}
	return nns
}

// AddNamespace register Namespace into bhojpur.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	nnl := oldToNewNs(nl)
	web.AddNamespace(nnl...)
}

// NSCond is Namespace Condition
func NSCond(cond namespaceCond) LinkNamespace {
	return func(namespace *Namespace) {
		web.NSCond(func(b *context.Context) bool {
			return cond((*adtContext.Context)(b))
		})
	}
}

// NSBefore Namespace BeforeRouter filter
func NSBefore(filterList ...FilterFunc) LinkNamespace {
	return func(namespace *Namespace) {
		nfs := oldToNewFilter(filterList)
		web.NSBefore(nfs...)
	}
}

// NSAfter add Namespace FinishRouter filter
func NSAfter(filterList ...FilterFunc) LinkNamespace {
	return func(namespace *Namespace) {
		nfs := oldToNewFilter(filterList)
		web.NSAfter(nfs...)
	}
}

// NSInclude Namespace Include ControllerInterface
func NSInclude(cList ...ControllerInterface) LinkNamespace {
	return func(namespace *Namespace) {
		nfs := oldToNewCtrlIntfs(cList)
		web.NSInclude(nfs...)
	}
}

// NSRouter call Namespace Router
func NSRouter(rootpath string, c ControllerInterface, mappingMethods ...string) LinkNamespace {
	return func(namespace *Namespace) {
		web.Router(rootpath, c, mappingMethods...)
	}
}

// NSGet call Namespace Get
func NSGet(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSGet(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSPost call Namespace Post
func NSPost(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.Post(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSHead call Namespace Head
func NSHead(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSHead(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSPut call Namespace Put
func NSPut(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSPut(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSDelete call Namespace Delete
func NSDelete(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSDelete(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSAny call Namespace Any
func NSAny(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSAny(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSOptions call Namespace Options
func NSOptions(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSOptions(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSPatch call Namespace Patch
func NSPatch(rootpath string, f FilterFunc) LinkNamespace {
	return func(ns *Namespace) {
		web.NSPatch(rootpath, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
}

// NSAutoRouter call Namespace AutoRouter
func NSAutoRouter(c ControllerInterface) LinkNamespace {
	return func(ns *Namespace) {
		web.NSAutoRouter(c)
	}
}

// NSAutoPrefix call Namespace AutoPrefix
func NSAutoPrefix(prefix string, c ControllerInterface) LinkNamespace {
	return func(ns *Namespace) {
		web.NSAutoPrefix(prefix, c)
	}
}

// NSNamespace add sub Namespace
func NSNamespace(prefix string, params ...LinkNamespace) LinkNamespace {
	return func(ns *Namespace) {
		nps := oldToNewLinkNs(params)
		web.NSNamespace(prefix, nps...)
	}
}

// NSHandler add handler
func NSHandler(rootpath string, h http.Handler) LinkNamespace {
	return func(ns *Namespace) {
		web.NSHandler(rootpath, h)
	}
}
